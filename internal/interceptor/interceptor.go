package interceptor

import (
	"context"
	"fmt"
	"slices"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/RIBorisov/GophKeeper/internal/log"
	"github.com/RIBorisov/GophKeeper/internal/model"
	"github.com/RIBorisov/GophKeeper/internal/service"
)

// UserIDUnaryInterceptor checks JWT token and injects userID into context.
func UserIDUnaryInterceptor(svc *service.Service, excludeMethods []string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// do not intercept if call method one of the excluded
		if slices.Contains(excludeMethods, info.FullMethod) {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Error("failed to get metadata")
			return nil, status.Error(codes.Unauthenticated, "Access denied")
		}

		token := md.Get("token")
		if len(token) == 0 {
			log.Error("got empty token")
			return nil, status.Error(codes.Unauthenticated, "Access denied")
		}

		userID := svc.GetUserID(token[0], svc.Cfg.Service.SecretKey)

		if userID == "" {
			return nil, status.Error(codes.Unauthenticated, "Access denied")
		}
		ctx = context.WithValue(ctx, model.CtxUserIDKey, userID)

		return handler(ctx, req)
	}
}

func interceptorLogger() logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l := log.GetLogger().With().Fields(fields).Logger()

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
			l.Debug().Msg(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	opts := []logging.Option{logging.WithLogOnEvents(logging.StartCall, logging.FinishCall)}

	return logging.UnaryServerInterceptor(interceptorLogger(), opts...)
}
