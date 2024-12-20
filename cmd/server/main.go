package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/RIBorisov/GophKeeper/internal/app/s3"
	"github.com/RIBorisov/GophKeeper/internal/app/server"
	"github.com/RIBorisov/GophKeeper/internal/config"
	"github.com/RIBorisov/GophKeeper/internal/log"
	"github.com/RIBorisov/GophKeeper/internal/service"
	"github.com/RIBorisov/GophKeeper/internal/storage"
)

func main() {
	log.InitLogger(zerolog.Level(0))
	log.Info("Logger has been initialized..")

	if err := initApp(); err != nil {
		log.Fatal("failed to initialize application", "err", err)
	}
}

func initApp() error {
	ctx := context.Background()

	g, ctx := errgroup.WithContext(ctx)
	cfg := config.Load()

	store, err := storage.Load(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to load storage: %w", err)
	}
	log.Info("Storage has been initialized..")

	s3client, err := s3.NewS3Client(ctx, cfg)
	svc := &service.Service{Cfg: cfg, Storage: store, S3Client: s3client}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	g.Go(func() error {
		if err = server.GRPCServe(svc, ch); err != nil {
			log.Fatal("failed to run grpc server")
		}
		return err
	})

	if err = g.Wait(); err != nil {
		return fmt.Errorf("unexpected error occurred: %w", err)
	}

	return nil
}
