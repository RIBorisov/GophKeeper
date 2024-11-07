package main

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"

	"github.com/RIBorisov/GophKeeper/internal/app/client"
)

func TestRunClient(t *testing.T) {
	tests := []struct {
		name    string
		tls     bool
		wantErr bool
		err     error
	}{
		{
			name:    "Positive #1",
			tls:     false,
			wantErr: true,
			err:     client.ErrToManyIncorrectValues,
		},
		{
			name:    "Positive #2",
			tls:     true,
			wantErr: true,
			err:     fs.ErrNotExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eg := &errgroup.Group{}
			eg.Go(func() error {
				return runClient(tt.tls)
			})

			assert.ErrorIs(t, eg.Wait(), tt.err)
		})
	}
}
