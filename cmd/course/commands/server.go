package commands

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/imrenagicom/demo-app/course/catalog"
	"github.com/imrenagicom/demo-app/course/server/apiserver"
	"github.com/imrenagicom/demo-app/internal/config"
	"github.com/imrenagicom/demo-app/internal/instrumentation"
	"github.com/imrenagicom/demo-app/internal/postgres"
	"github.com/imrenagicom/demo-app/internal/redis"
	"github.com/imrenagicom/demo-app/internal/util"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type serverOpts struct {
	envPrefix string
}

func newServer(opts *opts) *cobra.Command {
	serverOpts := &serverOpts{}
	command := &cobra.Command{
		Use:   "server",
		Short: "server subcommands",
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
		},
	}
	command.AddCommand(
		newServerStart(opts, serverOpts),
		newServerSeed(opts, serverOpts),
	)

	command.PersistentFlags().StringVar(&serverOpts.envPrefix, "env-prefix", "COURSE_SERVER", "config prefix")
	return command
}

func newServerStart(opts *opts, serverOpts *serverOpts) *cobra.Command {
	command := &cobra.Command{
		Use: "start",
		RunE: func(c *cobra.Command, args []string) error {
			conf, err := config.NewServer(opts.configPath, serverOpts.envPrefix)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to load config file")
			}
			logFn := instrumentation.InitializeLogger(conf.Log)
			defer logFn()

			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			signal.Notify(ch, syscall.SIGTERM)
			go func() {
				oscall := <-ch
				log.Warn().Msgf("system call:%+v", oscall)
				cancel()
			}()

			log.Debug().Msgf("running migration on %s", opts.migrationDir)
			if err := postgres.Migrate(opts.migrationDir, conf.DB.DatabaseUrl(), true); err != nil {
				log.Fatal().Err(err).Msg("unable to run migration")
			}

			server := apiserver.NewServer(apiserver.ServerOpts{
				Config: conf,
				Clients: &util.Clients{
					DB:    postgres.NewSQLx(conf.DB),
					Redis: redis.New(conf.Redis),
				},
			})
			return server.Run(ctx)
		},
	}
	return command
}

func newServerSeed(opts *opts, serverOpts *serverOpts) *cobra.Command {
	command := &cobra.Command{
		Use:   "seed",
		Short: "seed db",
		RunE: func(c *cobra.Command, args []string) error {
			conf, err := config.NewServer(opts.configPath, serverOpts.envPrefix)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to load config file")
			}
			logFn := instrumentation.InitializeLogger(conf.Log)
			defer logFn()

			ctx := log.With().Logger().WithContext(context.Background())
			ctx, cancel := context.WithCancel(ctx)
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			signal.Notify(ch, syscall.SIGTERM)
			go func() {
				oscall := <-ch
				log.Warn().Msgf("system call:%+v", oscall)
				cancel()
			}()

			clients := &util.Clients{
				DB: postgres.NewSQLx(conf.DB),
			}
			concertStore := catalog.NewStore(clients.DB, clients.Redis)
			catalogSvc := catalog.NewService(concertStore, clients.DB)
			return catalogSvc.Seed(ctx)
		},
	}
	return command
}
