package types

import (
	"errors"
)

var (
	// ErrActionType --
	ErrActionType = errors.New("missing action type")
	// ErrEntityType --
	ErrEntityType = errors.New("missing entity type")
	// ErrDuplicate --
	ErrDuplicate = errors.New("duplicate value")
	// ErrCollision --
	ErrCollision = errors.New("key collision")
	// ErrDNE --
	ErrDNE = errors.New("does not exist")
	// ErrUser --
	ErrUser = errors.New("missing user id")
	// ErrDataType --
	ErrDataType = errors.New("incorrect data type")
	// ErrValueType --
	ErrValueType = errors.New("incorrect value type")
	// ErrKeyType --
	ErrKeyType = errors.New("incorrect key type")
	// ErrID --
	ErrID = errors.New("missing id")
	// ErrStatus --
	ErrStatus = errors.New("missing status")
	// ErrSplit --
	ErrSplit = errors.New("missing split")
	// ErrOriginType --
	ErrOriginType = errors.New("missing origin type")
	// ErrSplitInvalid --
	ErrSplitInvalid = errors.New("invalid split value")
	// ErrVariants --
	ErrVariants = errors.New("missing variants")
	// ErrExists --
	ErrExists = errors.New("already exists")
	// ErrITypes --
	ErrITypes = errors.New("missing itypes")
	// ErrRewards --
	ErrRewards = errors.New("missing rewards")
	// ErrLength --
	ErrLength = errors.New("incorrect length")
	// ErrTimestamp --
	ErrTimestamp = errors.New("invalid timestamp format")
	// ErrBounds --
	ErrBounds = errors.New("out of bounds")
	// ErrAssertion --
	ErrAssertion = errors.New("type assertion error")
	// ErrResourceType --
	ErrResourceType = errors.New("invalid resource type")
)
