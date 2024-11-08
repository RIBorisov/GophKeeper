package main

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/RIBorisov/GophKeeper/internal/log"
)

func Test_initApp(t *testing.T) {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer shutdownRelease()
	log.InitLogger(zerolog.Level(0))
	go func() {
		assert.Nil(t, initApp())
	}()
	<-shutdownCtx.Done()
}
