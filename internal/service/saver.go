package service

import (
	"context"
	"fmt"
	"os"

	"github.com/RIBorisov/GophKeeper/internal/log"
	"github.com/RIBorisov/GophKeeper/internal/model"
	pb "github.com/RIBorisov/GophKeeper/pkg/server"
)

func (s *Service) saveText(ctx context.Context, fileName, text string) error {
	file, err := saveFile(fileName, s.Cfg.Service.SecretKey, []byte(text))
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Error("failed to close file", "fileName", fileName)
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	size := fileInfo.Size()

	if err = s.S3Client.PutObject(ctx, fileName, file, size); err != nil {
		return fmt.Errorf("failed to put file: %w", err)
	}
	return nil
}

func (s *Service) saveBytes(ctx context.Context, fileName string, text []byte) error {
	file, err := saveFile(fileName, s.Cfg.Service.SecretKey, text)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Error("failed to close file", "fileName", fileName)
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	size := fileInfo.Size()

	if err = s.S3Client.PutObject(ctx, fileName, file, size); err != nil {
		return fmt.Errorf("failed to put file: %w", err)
	}

	return nil
}

func (s *Service) saveCard(ctx context.Context, fileName string, in *pb.Card) error {
	fileData := in.GetNumber() + " " + in.GetMmYy() + " " + in.GetCvc()

	file, err := saveFile(fileName, s.Cfg.Service.SecretKey, []byte(fileData))
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Error("failed to close file", "fileName", fileName)
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	size := fileInfo.Size()

	if err = s.S3Client.PutObject(ctx, fileName, file, size); err != nil {
		return fmt.Errorf("failed to put file: %w", err)
	}

	return nil
}

// SaveMeta saves metadata into db.
func (s *Service) SaveMeta(ctx context.Context, metadata model.Save) (string, error) {
	log.Debug("going to save metadata into postgres..")

	id, err := s.Storage.Save(ctx, metadata)
	if err != nil {
		return "", fmt.Errorf("failed to save: %w", err)
	}

	return id, nil
}

// saveFile encrypts data, saves and returns it as *os.File.
func saveFile(fname, secret string, data []byte) (*os.File, error) {
	encrypted, err := Encrypt(secret, data)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt: %w", err)
	}

	if err := os.WriteFile("tmp/"+fname, encrypted, 0600); err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	file, err := os.Open("tmp/" + fname)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}
