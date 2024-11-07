package main

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/RIBorisov/GophKeeper/internal/app/client"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		tls     bool
		wantErr bool
		err     error
	}{
		{
			name:    "Positive #1",
			tls:     true,
			wantErr: true,
			err:     fs.ErrNotExist,
		},
		{
			name:    "Positive #2",
			tls:     false,
			wantErr: true,
			err:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.NewClient(tt.tls)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}
