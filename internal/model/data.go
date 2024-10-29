package model

import (
	"github.com/google/uuid"

	"github.com/RIBorisov/GophKeeper/pkg/server"
)

type (
	// Kind is a type of saving data.
	Kind string

	// UserID is a user identification value.
	UserID string

	// ID represent unique value of stored data entity.
	ID uuid.UUID
)

const (
	// Unspecified kind of saving data.
	Unspecified Kind = "unknown"

	// CardCredentials kind of saving data.
	CardCredentials Kind = "card"

	// Text kind of saving data.
	Text Kind = "text"

	// Credentials kind of saving data.
	Credentials Kind = "credentials"

	// Binary kind of saving data.
	Binary Kind = "binary"
)

// Save model struct.
type Save struct {
	ID   string
	Kind Kind
	Data string
	Meta string
}

// ToPB method converts custom type Kind to protobuf value.
func (k Kind) ToPB() server.Kind {
	switch k {
	case CardCredentials:
		return server.Kind_CARD
	case Text:
		return server.Kind_TEXT
	case Credentials:
		return server.Kind_CREDENTIALS
	case Binary:
		return server.Kind_BINARY
	default:
		return server.Kind_UNSPECIFIED
	}
}

// ToKind method converts protobuf to custom Kind.
func ToKind(k server.Kind) Kind {
	switch k {
	case server.Kind_CARD:
		return CardCredentials
	case server.Kind_TEXT:
		return Text
	case server.Kind_CREDENTIALS:
		return Credentials
	case server.Kind_BINARY:
		return Binary
	default:
		return Unspecified
	}
}

// String represents k as string.
func (k Kind) String() string {
	return string(k)
}
