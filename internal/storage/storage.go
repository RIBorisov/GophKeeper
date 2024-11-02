package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/RIBorisov/GophKeeper/internal/config"
	"github.com/RIBorisov/GophKeeper/internal/log"
	"github.com/RIBorisov/GophKeeper/internal/model"
)

// DB provides methods for interacting with the database.
// It encapsulates the database connection pool and configuration settings.
type DB struct {
	*DBPool
	cfg *config.Config
}

// Load method initializes database pool.
func Load(ctx context.Context, cfg *config.Config) (*DB, error) {
	pool, err := NewPool(ctx, cfg.App.PgDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain pool: %w", err)
	}
	return &DB{DBPool: pool, cfg: cfg}, nil
}

// Register method saves user into database.
func (d *DB) Register(ctx context.Context, in model.UserCredentials) (model.UserID, error) {
	const (
		stmt          = `INSERT INTO users (login, password) VALUES (@login, @password) RETURNING "id"`
		constrainName = "idx__login_is_unique"
	)
	args := pgx.NamedArgs{"login": in.Login, "password": in.Password}

	var userID string
	if err := d.pool.QueryRow(ctx, stmt, args).Scan(&userID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == constrainName {
			return "", ErrUserExists
		}

		return "", fmt.Errorf("failed to scan row into userID: %w", err)
	}

	return model.UserID(userID), nil
}

// UserEntity describes user database entity.
type UserEntity struct {
	ID       string `db:"id"`
	Login    string `db:"login"`
	Password string `db:"id"`
}

// GetUser method gets user entity from database.
func (d *DB) GetUser(ctx context.Context, login string) (*UserEntity, error) {
	const stmt = `SELECT id, login, password FROM users WHERE login = @login`

	log.Debug("going to get from db", "user", login)
	u := &UserEntity{}
	var pass []byte
	arg := pgx.NamedArgs{"login": login}
	err := d.pool.QueryRow(ctx, stmt, arg).Scan(&u.ID, &u.Login, &pass)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to scan row: %w", err)
	}
	u.Password = string(pass)
	log.Debug("got row from database", "UserEntity", *u)

	return u, nil
}

// Save method saves user data (e.g. kind, meta) into database.
func (d *DB) Save(ctx context.Context, data model.Save) (string, error) {
	const stmt = `INSERT INTO metadata(id, user_id, kind, metadata) VALUES (@id, @user_id, @kind, @metadata) RETURNING id`
	userID, err := getUserFromCtx(ctx)
	if err != nil {
		return "", err
	}
	var id string
	args := pgx.NamedArgs{"id": data.ID, "user_id": userID, "kind": data.Kind.String(), "metadata": data.Meta}
	if err := d.pool.QueryRow(ctx, stmt, args).Scan(&id); err != nil {
		log.Error("failed to execute", "stmt", stmt, "userID", userID, "err", err)
		return "", fmt.Errorf("failed to exec stmt: %w", err)
	}

	return id, nil
}

type MetadataEntity struct {
	ID       string     `db:"id"`
	Metadata string     `db:"metadata"`
	Kind     model.Kind `db:"kind"`
}

// Get method gets metadata from database.
func (d *DB) Get(ctx context.Context, id string) (*MetadataEntity, error) {
	const stmt = `SELECT metadata, kind FROM metadata WHERE id=@id AND user_id=@userID`
	userID, err := getUserFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var e MetadataEntity

	err = d.pool.QueryRow(ctx, stmt, pgx.NamedArgs{"id": id, "userID": userID}).Scan(&e.Metadata, &e.Kind)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMetadataNotFound
		}
		return nil, fmt.Errorf("failed to scan row into metadata: %w", err)
	}

	return &e, nil
}

func getUserFromCtx(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(model.CtxUserIDKey).(string)
	if !ok {
		return "", ErrUserNotFound
	}
	log.Debug("Got user from ctx", "userID", userID)

	return userID, nil
}

var (
	// ErrUserExists is returned when the user already exists.
	ErrUserExists = errors.New("user already exists")
	// ErrUserNotFound is returned when the requested user is not found.
	ErrUserNotFound = errors.New("user not found")
	// ErrMetadataNotFound is returned when the requested ID is not found.
	ErrMetadataNotFound = errors.New("metadata not found")
)
