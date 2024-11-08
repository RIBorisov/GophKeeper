package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "Positive #1",
			want: &Config{
				App: AppConfig{
					Addr:        ":50051",
					PgDSN:       "postgresql://admin:password@localhost:5432/gophkeeper",
					CertPath:    "tls/server.crt",
					CertKeyPath: "tls/server.key",
					TLSEnabled:  true,
				},
				Service: ServiceConfig{
					SecretKey: "super",
				},
				S3: S3Config{
					BucketName:      "tests",
					Endpoint:        "localhost:9999",
					AccessKeyID:     "adm",
					SecretAccessKey: "pwd",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SERVER_ADDRESS", ":50051")
			t.Setenv("POSTGRES_DSN", "postgresql://admin:password@localhost:5432/gophkeeper")
			t.Setenv("CERT_PATH", "tls/server.crt")
			t.Setenv("CERT_KEY_PATH", "tls/server.key")
			t.Setenv("TLS_ENABLED", "1")
			t.Setenv("SECRET_KEY", "super")
			t.Setenv("S3_BUCKET_NAME", "tests")
			t.Setenv("S3_ENDPOINT", "localhost:9999")
			t.Setenv("S3_AK_ID", "adm")
			t.Setenv("S3_SECRET_AK", "pwd")
			assert.Equal(t, tt.want, Load())
		})
	}
}
