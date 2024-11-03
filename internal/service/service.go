package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/RIBorisov/GophKeeper/internal/app/s3"
	"github.com/RIBorisov/GophKeeper/internal/config"
	"github.com/RIBorisov/GophKeeper/internal/log"
	"github.com/RIBorisov/GophKeeper/internal/model"
	"github.com/RIBorisov/GophKeeper/internal/storage"
	pb "github.com/RIBorisov/GophKeeper/pkg/server"
)

type Store interface {
	Register(ctx context.Context, in model.UserCredentials) (model.UserID, error)
	GetUser(ctx context.Context, login string) (*storage.UserEntity, error)
	Save(ctx context.Context, data model.Save) (string, error)
	Get(ctx context.Context, id string) (*storage.MetadataEntity, error)
	GetMany(ctx context.Context) ([]*storage.MetadataEntity, error)
}

type Service struct {
	Storage  Store
	S3Client *s3.S3Client
	Cfg      *config.Config
}

// RegisterUser encrypts password, saves user login and password into database.
func (s *Service) RegisterUser(ctx context.Context, in model.UserCredentials) (string, error) {
	encrypted, err := hashPassword(s.Cfg.Service.SecretKey, in.Password)
	in.Password = encrypted

	userID, err := s.Storage.Register(ctx, in)
	if err != nil {
		return "", fmt.Errorf("failed to register user: %w", err)
	}

	token, err := BuildJWTString(s.Cfg.Service.SecretKey, string(userID))
	if err != nil {
		return "", fmt.Errorf("failed to build token: %w", err)
	}
	return token, nil
}

func (s *Service) AuthUser(ctx context.Context, in model.UserCredentials) (string, error) {
	u, err := s.Storage.GetUser(ctx, in.Login)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	if err = comparePasswords(s.Cfg.Service.SecretKey, u.Password, in.Password); err != nil {
		return "", ErrInvalidPassword
	}
	authToken, err := BuildJWTString(s.Cfg.Service.SecretKey, u.ID)
	if err != nil {
		return "", fmt.Errorf("failed to build token: %w", err)
	}

	return authToken, nil
}

// SaveData encrypts and saves user data.
func (s *Service) SaveData(ctx context.Context, in *pb.SaveRequest) (string, error) {
	// TODO: продумать ограничение на байты
	fileName := uuid.New().String()

	switch in.GetData().(type) {
	case *pb.SaveRequest_Binary:
		log.Debug("going to save binary data..")
		if err := s.saveBytes(ctx, fileName, in.GetBinary()); err != nil {
			return "", fmt.Errorf("failed to save binary: %w", err)
		}

	case *pb.SaveRequest_Text:
		log.Debug("going to save text data..")
		if err := s.saveText(ctx, fileName, in.GetText()); err != nil {
			return "", fmt.Errorf("failed to save text: %w", err)
		}

	case *pb.SaveRequest_Card:
		log.Debug("going to save card data..")
		if err := s.saveCard(ctx, fileName, in.GetCard()); err != nil {
			return "", fmt.Errorf("failed to save card: %w", err)
		}

	case *pb.SaveRequest_Credentials:
		log.Debug("going to save credentials..")

		creds := pb.Credentials{Login: in.GetCredentials().Login, Password: in.GetCredentials().Password}
		b, err := proto.Marshal(&creds)
		if err != nil {
			return "", fmt.Errorf("failed to marshal text into pb: %w", err)
		}
		if err = s.saveBytes(ctx, fileName, b); err != nil {
			return "", fmt.Errorf("failed to save credentials")
		}

	default:
		return "", ErrUnsupportedType
	}

	log.Debug("successfully saved data", "kind", in.GetKind().String())

	metaDTO := model.Save{ID: fileName, Kind: model.ToKind(in.GetKind()), Meta: in.GetMeta()}
	id, err := s.SaveMeta(ctx, metaDTO)
	if err != nil {
		return "", fmt.Errorf("failed to save metadata into postgres: %w", err)
	}

	log.Debug("successfully saved metadata", "kind", in.GetKind().String())

	return id, nil
}

func (s *Service) GetData(ctx context.Context, id string) (*pb.GetResponse, error) {
	log.Debug("going to get metadata", "id", id)
	raw, err := s.Storage.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata from storage: %w", err)
	}

	log.Debug("going to get object from s3 storage", "object name", id)
	obj, err := s.S3Client.GetObject(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3: %w", err)
	}

	// decrypts object data
	decrypted, err := Decrypt(s.Cfg.Service.SecretKey, obj)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt object: %w", err)
	}

	res := &pb.GetResponse{
		ID:   id,
		Kind: raw.Kind.ToPB(),
	}

	switch raw.Kind {
	case model.CardCredentials:
		parts := strings.Split(string(decrypted), " ")
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid data format: %w", err)
		}

		res.Data = &pb.GetResponse_Card{Card: &pb.Card{
			Number: parts[0],
			MmYy:   parts[1],
			Cvc:    parts[2],
		}}

	case model.Text:
		res.Data = &pb.GetResponse_Text{Text: string(decrypted)}

	case model.Credentials:
		var credentials pb.Credentials
		if err = proto.Unmarshal(decrypted, &credentials); err != nil {
			return nil, fmt.Errorf("failed to unmarshal credentials: %w", err)
		}
		res.Data = &pb.GetResponse_Credentials{Credentials: &credentials}

	case model.Binary:
		res.Data = &pb.GetResponse_Binary{Binary: decrypted}

	default:
		return nil, fmt.Errorf("unsupported kind: %v", raw.Kind)
	}

	return res, nil
}

// GetUserData retrieves all user data. The method could be used for synchronize purposes.
func (s *Service) GetUserData(ctx context.Context) (*pb.GetManyResponse, error) {
	raw, err := s.Storage.GetMany(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from storage: %w", err)
	}

	var userData []*pb.GetResponse
	for _, r := range raw {
		data, err := s.GetData(ctx, r.ID)
		if err != nil {
			return nil, err
		}
		userData = append(userData, &pb.GetResponse{
			ID:   data.GetID(),
			Kind: data.GetKind(),
			Data: data.GetData(),
		})
	}

	return &pb.GetManyResponse{UserData: userData}, nil
}

var (
	ErrInvalidPassword = errors.New("invalid login and(or) password")
	ErrUnsupportedType = errors.New("unsupported type")
)
