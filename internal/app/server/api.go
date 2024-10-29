package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

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
func GRPCServe(svc *service.Service) error {
	listen, err := net.Listen("tcp", svc.Cfg.App.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen a port: %w", err)
	}

	exclude := []string{
		"/GophKeeperService/Register",
		"/GophKeeperService/Auth",
	}
	_ = exclude
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor.UserIDUnaryInterceptor(svc, exclude)))
	pb.RegisterGophKeeperServiceServer(s, &GRPCServer{svc: svc})
	reflection.Register(s)
	log.Info("Starting gRPC server...")

	return s.Serve(listen)
}

// Register method uses for customer registration.
func (g *GRPCServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	t, err := g.svc.RegisterUser(ctx, model.UserCredentials{Login: in.GetLogin(), Password: in.GetPassword()})
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, "")
	}

	md := metadata.New(map[string]string{"token": t})
	if err = grpc.SetHeader(ctx, md); err != nil {
		log.Error("failed to set token header")
		return nil, status.Error(codes.Internal, "")
	}

	return &pb.RegisterResponse{Result: "Success"}, nil
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

	md := metadata.New(map[string]string{"token": t})
	if err = grpc.SetHeader(ctx, md); err != nil {
		log.Error("failed to set token header")
		return nil, status.Error(codes.Internal, "")
	}

	return &pb.AuthResponse{Result: "Success"}, nil
}

// Save method uses to save user's data.
func (g *GRPCServer) Save(ctx context.Context, in *pb.SaveRequest) (*pb.SaveResponse, error) {
	id, err := g.svc.SaveData(ctx, in)
	if err != nil {
		log.Error("failed save data", "err", err, "input values", in)
		return nil, status.Error(codes.Internal, "")
	}

	// нужно сохранять разные типы данных
	// через свитч получаем тип данных и формируем модель,
	// с моделью идем в сервис, в конкретный отдельный метод и сохраняем в базе и с3
	// TODO: потом навесить шифрование

	return &pb.SaveResponse{ID: id}, nil
}

// Get method uses to retrieve user's data.
func (g *GRPCServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	data, err := g.svc.GetData(ctx, in.GetID())
	if err != nil {
		log.Error("failed to get data", "metadata id", in.GetID(), "err", err)
		if errors.Is(err, storage.ErrMetadataNotFound) {
			return nil, status.Error(codes.NotFound, "Not found")
		}
		return nil, status.Error(codes.Internal, "")
	}
	res := &pb.GetResponse{ID: data.GetID(), Kind: data.GetKind(), Data: data.GetData()}

	return res, nil
}
