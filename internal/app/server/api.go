package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/RIBorisov/GophKeeper/internal/app/cert"
	"github.com/RIBorisov/GophKeeper/internal/interceptor"
	"github.com/RIBorisov/GophKeeper/internal/log"
	"github.com/RIBorisov/GophKeeper/internal/model"
	"github.com/RIBorisov/GophKeeper/internal/service"
	"github.com/RIBorisov/GophKeeper/internal/storage"
	pb "github.com/RIBorisov/GophKeeper/pkg/server"
)

// GRPCServer describes server.
type GRPCServer struct {
	pb.UnimplementedGophKeeperServiceServer
	svc *service.Service
}

// GRPCServe runs the gRPC server.
func GRPCServe(svc *service.Service, stopCh chan os.Signal) error {
	listen, err := net.Listen("tcp", svc.Cfg.App.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen a port: %w", err)
	}

	exclude := []string{
		"/GophKeeperService/Register",
		"/GophKeeperService/Auth",
	}

	if err = cert.PrepareTLS(svc.Cfg.App.CertPath, svc.Cfg.App.CertKeyPath); err != nil {
		return fmt.Errorf("failed to prepare TLS: %w", err)
	}

	creds, err := credentials.NewServerTLSFromFile(svc.Cfg.App.CertPath, svc.Cfg.App.CertKeyPath)
	if err != nil {
		return fmt.Errorf("failed to construct TLS credentials: %w", err)
	}
	s := grpc.NewServer(
		grpc.Creds(creds),
		grpc.ChainUnaryInterceptor(
			interceptor.UnaryServerInterceptor(),
			interceptor.UserIDUnaryInterceptor(svc, exclude),
		),
	)
	pb.RegisterGophKeeperServiceServer(s, &GRPCServer{svc: svc})
	reflection.Register(s)
	log.Info("Starting gRPC server..", "Addr", svc.Cfg.App.Addr)

	go func() {
		<-stopCh
		log.Info("got signal to stop server..")
		s.GracefulStop()
		log.Info("server gracefully stopped..")
	}()

	return s.Serve(listen)
}

// Register method uses for customer registration.
func (g *GRPCServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Debug("Going to register", "user", in.GetLogin())
	t, err := g.svc.RegisterUser(ctx, model.UserCredentials{Login: in.GetLogin(), Password: in.GetPassword()})
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, "")
	}

	return &pb.RegisterResponse{Token: t}, nil
}

// Auth method uses for customer authentication.
func (g *GRPCServer) Auth(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {
	log.Debug("going to authenticate", "user", in.GetLogin())
	t, err := g.svc.AuthUser(ctx, model.UserCredentials{Login: in.GetLogin(), Password: in.GetPassword()})
	if err != nil {
		if errors.Is(err, service.ErrInvalidPassword) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		log.Error("failed to authenticate", "Login", in.GetLogin(), "err", err.Error())

		return nil, status.Error(codes.Internal, "")
	}

	return &pb.AuthResponse{Token: t}, nil
}

// Save method uses to save user's data.
func (g *GRPCServer) Save(ctx context.Context, in *pb.SaveRequest) (*pb.SaveResponse, error) {
	id, err := g.svc.SaveData(ctx, in)
	if err != nil {
		log.Error("failed save data", "err", err, "input values", in)
		return nil, status.Error(codes.Internal, "")
	}

	return &pb.SaveResponse{ID: id}, nil
}

// Get method uses to retrieve user's data.
func (g *GRPCServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	data, err := g.svc.GetData(ctx, in.GetID())
	if err != nil {
		log.Error("failed to get data", "id", in.GetID(), "err", err)
		if errors.Is(err, storage.ErrMetadataNotFound) {
			return nil, status.Error(codes.NotFound, "Not found")
		}
		return nil, status.Error(codes.Internal, "")
	}

	return &pb.GetResponse{
		ID:   data.GetID(),
		Kind: data.GetKind(),
		Data: data.GetData(),
	}, nil
}

// GetMany method uses to retrieve all user's data.
func (g *GRPCServer) GetMany(ctx context.Context, _ *pb.GetManyRequest) (*pb.GetManyResponse, error) {
	data, err := g.svc.GetUserData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user data: %w", err)
	}

	return data, nil
}
