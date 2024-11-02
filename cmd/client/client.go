package main

import (
	"context"

	"github.com/RIBorisov/GophKeeper/internal/app/client"
	"github.com/RIBorisov/GophKeeper/internal/log"
)

func main() {
	ctx := context.Background()
	c, err := client.NewClient(ctx)
	if err != nil {
		log.Fatal("failed to start client", "err", err)
	}

	c.ListenAction(ctx)
}
