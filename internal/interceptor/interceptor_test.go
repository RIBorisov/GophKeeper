package interceptor

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	s3mock "github.com/RIBorisov/GophKeeper/internal/app/s3/mocks"
	"github.com/RIBorisov/GophKeeper/internal/config"
	"github.com/RIBorisov/GophKeeper/internal/log"
	"github.com/RIBorisov/GophKeeper/internal/service"
	svcmock "github.com/RIBorisov/GophKeeper/internal/service/mocks"
)

func TestUserIDUnaryInterceptor(t *testing.T) {
	const validToken = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
eyJleHAiOjE3MzM1OTQxNzUsIlVzZXJJRCI6IjllZDlmYWZjLTFhYjEtNGMzOC04ODdmLTRkNjViZmFhNDFiNCJ9.
74N-oSqFZOSnd7Me2bYSvagO64_PF0lBhieSrM9jGvw`
	excluded := []string{"/GophKeeperService/Register", "/GophKeeperService/Auth"}
	tests := []struct {
		name           string
		method         string
		excludeMethods []string
		token          string
		wantErr        bool
	}{
		{
			name:           "Positive #1",
			excludeMethods: excluded,
			method:         "/GophKeeperService/Register",
			token:          validToken,
			wantErr:        false,
		},
		{
			name:           "Positive #2",
			excludeMethods: excluded,
			method:         "/GophKeeperService/Auth",
			token:          validToken,
			wantErr:        false,
		},
		{
			name:           "Negative #1",
			method:         "/GophKeeperService/Save",
			excludeMethods: excluded,
			token:          "invalid token",
			wantErr:        true,
		},
		{
			name:           "Negative #2",
			method:         "/GophKeeperService/Get",
			excludeMethods: excluded,
			token:          "invalid token",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.InitLogger(zerolog.Level(0))
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := svcmock.NewMockStoreI(ctrl)

			svc := service.Service{
				Cfg:      &config.Config{},
				S3Client: s3mock.NewMockS3ClientI(ctrl),
				Storage:  store,
			}

			interceptor := UserIDUnaryInterceptor(&svc, tt.excludeMethods)

			ctx := context.Background()
			md := metadata.New(map[string]string{"token": tt.token})
			ctx = metadata.NewIncomingContext(ctx, md)

			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return req, nil
			}

			_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: tt.method}, handler)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
