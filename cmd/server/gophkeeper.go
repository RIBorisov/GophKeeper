package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/RIBorisov/GophKeeper/internal/app/server"
	"github.com/RIBorisov/GophKeeper/internal/config"
	"github.com/RIBorisov/GophKeeper/internal/log"
	"github.com/RIBorisov/GophKeeper/internal/service"
	"github.com/RIBorisov/GophKeeper/internal/storage"
)

func main() {
	log.InitLogger(zerolog.Level(0))
	log.Info("Logger has been initialized..")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	cfg := config.Load()

	store, err := storage.Load(ctx, cfg)
	if err != nil {
		log.Fatal("failed to load storage", err)
	}
	log.Info("Storage has been initialized..")

	svc := &service.Service{Cfg: cfg, Storage: store}

	g.Go(func() error { return server.GRPCServe(svc) })
	if err = g.Wait(); err != nil {
		log.Fatal("unexpected error occurred", err)
	}
}
