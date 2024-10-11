package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/RIBorisov/GophKeeper/internal/config"
	"github.com/RIBorisov/GophKeeper/internal/log"
)

func main() {
	log.InitLogger(zerolog.Level(0))
	log.Info("Logger has been initialized...")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	cfg := config.Load()
	log.Debug("Loaded config", "config", &cfg)
	log.Warning("WARN Message")
	log.Info("INFO Message")
	log.Error("ERROR Message")
	_ = cfg
	err := g.Wait()
	if err != nil {
		log.Fatal("FATAL ERROR")
	}

}
