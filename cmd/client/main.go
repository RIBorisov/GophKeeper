package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/RIBorisov/GophKeeper/internal/app/client"
	"github.com/RIBorisov/GophKeeper/internal/log"
)

var (
	buildDate    = "N/A"
	buildVersion = "N/A"
)

func main() {
	eg := &errgroup.Group{}
	eg.Go(func() error {
		return runClient(true)
	})

	if err := eg.Wait(); err != nil {
		if !errors.Is(err, client.ErrToManyIncorrectValues) {
			log.Fatal("failed to run client", "err", err)
		}
		fmt.Println("Exit application..")
		os.Exit(1)
	}
}

func runClient(tlsEnabled bool) error {
	c, err := client.NewClient(tlsEnabled)
	if err != nil {
		return fmt.Errorf("failed to start client: %w", err)
	}
	return c.ListenAction(context.Background(), buildDate, buildVersion)
}
