package interceptor

import (
	"context"
	"slices"

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
