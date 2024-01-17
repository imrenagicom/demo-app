package commands

import (
	"github.com/spf13/cobra"
)

type opts struct {
	configPath   string
	migrationDir string
}

func NewCommand() *cobra.Command {
	opts := &opts{}
	command := &cobra.Command{
		Use:   "course",
		Short: "all cli commands for course service",
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
		},
	}
	command.AddCommand(
		newServer(opts),
	)
	command.PersistentFlags().StringVar(&opts.configPath, "config", "/etc/course/conf/server.yaml", "path to config file")
	command.PersistentFlags().StringVar(&opts.migrationDir, "migration", "/etc/course/migrations", "migration directory")
	return command
}
