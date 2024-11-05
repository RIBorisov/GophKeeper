package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_initApp(t *testing.T) {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer shutdownRelease()

	go func() {
		assert.Nil(t, initApp())
	}()
	<-shutdownCtx.Done()

}
