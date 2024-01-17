package grpc

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func Logger() logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l := log.Ctx(ctx).With().Fields(fields).Logger()
		switch lvl {
		case logging.LevelDebug:
			l.Debug().Msg(msg)
		case logging.LevelInfo:
			l.Info().Msg(msg)
		case logging.LevelWarn:
			l.Warn().Msg(msg)
		case logging.LevelError:
			l.Error().Msg(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

var loggingOpts = []logging.Option{
	logging.WithLogOnEvents(
		logging.StartCall,
		logging.PayloadReceived,
		logging.PayloadSent,
		logging.FinishCall,
	),
}

func StreamServerGRPCLoggerInterceptor(opts ...logging.Option) grpc.StreamServerInterceptor {
	options := loggingOpts
	if len(opts) > 0 {
		options = opts
	}
	return logging.StreamServerInterceptor(Logger(), options...)
}

func UnaryServerGRPCLoggerInterceptor(opts ...logging.Option) grpc.UnaryServerInterceptor {
	options := loggingOpts
	if len(opts) > 0 {
		options = opts
	}
	return logging.UnaryServerInterceptor(Logger(), options...)
}

func UnaryClientGRPCLoggerInterceptor(opts ...logging.Option) grpc.UnaryClientInterceptor {
	options := loggingOpts
	if len(opts) > 0 {
		options = opts
	}
	return logging.UnaryClientInterceptor(Logger(), options...)
}

func StreamClientGRPCLoggerInterceptor(opts ...logging.Option) grpc.StreamClientInterceptor {
	options := loggingOpts
	if len(opts) > 0 {
		options = opts
	}
	return logging.StreamClientInterceptor(Logger(), options...)
}

func UnaryServerAppLoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := log.With().Str("request_id", uuid.New().String()).Logger()
		return handler(log.WithContext(ctx), req)
	}
}

func StreamServerAppLoggerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, newWrappedStream(ss))
		if err != nil {
			log.Error().Err(err).Msgf("Error: %v", err)
			return err
		}
		return nil
	}
}

type wrappedStream struct {
	grpc.ServerStream
}

func (w *wrappedStream) Context() context.Context {
	log := log.With().Str("request_id", uuid.New().String()).
		Logger()
	return log.WithContext(context.Background())
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}
