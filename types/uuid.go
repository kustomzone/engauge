package types

import (
	"github.com/gofrs/uuid"
)

var (
	UUIDNil = &UUID{uuid.Nil}
)

// UUID --
type UUID struct{ uuid.UUID }

// NewUUID --
func NewUUID() *UUID {
	id, _ := uuid.NewV4()
	return &UUID{UUID: id}
}

// UUIDFromString --
func UUIDFromString(input string) *UUID {
	id := uuid.FromStringOrNil(input)
	if id == uuid.Nil {
		return nil
	}
	return &UUID{id}
}
