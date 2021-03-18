package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/JKhawaja/errors"
)

var (
	/* special keys */
	// Conversion is a special key (used in automated analysis)
	Conversion = "conversion"

	/* timezone */

	// DefaultTimeZone is the default timezone for the system (defaults to local time)
	DefaultTimeZone = "Local"
)

// Interaction represents the full structure
// of an interaction object
// must be unique on: {Action, UserID, Timestamp}
type Interaction struct {
	// how
	Action *string `json:"action,omitempty"`

	// what
	EntityType *string `json:"entityType,omitempty"`
	EntityID   *string `json:"entityID,omitempty"`

	// where
	OriginType *string `json:"originType,omitempty"`
	OriginID   *string `json:"originID,omitempty"`

	// who
	UserType *string `json:"userType,omitempty"`
	UserID   *string `json:"userID,omitempty"`

	DeviceType *string `json:"deviceType,omitempty"`
	DeviceID   *string `json:"deviceID,omitempty"`

	// why (context)
	SessionType *string `json:"sessionType,omitempty"`
	SessionID   *string `json:"sessionID,omitempty"`

	// when
	Timestamp  *string    `json:"timestamp,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	ReceivedAt *time.Time `json:"receivedAt,omitempty"`

	// metadata: entity-properties, origin-properties, user-properties, session-properties, etc.
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// Validate will return an error if interaction object is not valid.
func (i *Interaction) Validate() error {
	// action type is always required
	if i.Action == nil {
		return errors.New(ErrActionType, nil)
	}

	// user-id is always required
	if i.UserID == nil {
		return errors.New(ErrUser, nil)
	}

	// all properties must have numerical, numerical-array, text, or text-array values
	if i.Properties != nil {
		for key, value := range i.Properties {
			switch value.(type) {
			case float64, []float64:
				continue
			case string, []string:
				continue
			default:
				delete(i.Properties, key)
			}
		}
	}

	return nil
}

// Date --
func (i *Interaction) Date() string {
	t := *i.CreatedAt
	return fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day())
}

// Entity --
func (i *Interaction) Entity() *Entity {
	var entityType, entityID string
	if i.EntityType != nil {
		entityType = *i.EntityType
	}

	if i.EntityID != nil {
		entityID = *i.EntityID
	}

	return &Entity{
		EntityType: entityType,
		EntityID:   entityID,
	}
}

// String --
func (i *Interaction) String() string {
	s := make([]string, 0, 10)
	t := *i.CreatedAt

	s = append(s, pstr(i.Action))
	s = append(s, pstr(i.EntityType))
	s = append(s, pstr(i.EntityID))
	s = append(s, pstr(i.OriginType))
	s = append(s, pstr(i.OriginID))
	s = append(s, pstr(i.UserType))
	s = append(s, pstr(i.UserID))
	s = append(s, pstr(i.DeviceType))
	s = append(s, pstr(i.DeviceID))
	s = append(s, t.String())

	return strings.Join(s, "-")
}

// Endpoint will return the endpoint of the interaction
func (i *Interaction) Endpoint() *Endpoint {
	return &Endpoint{
		Action:     i.Action,
		EntityType: i.EntityType,
		EntityID:   i.EntityID,
		OriginType: i.OriginType,
		OriginID:   i.OriginID,
	}
}

// Origin will return the full origin object of the interaction
func (i *Interaction) Origin() *Origin {
	return &Origin{
		OriginType: i.OriginType,
		OriginID:   i.OriginID,
	}
}

// CSV --
func (i *Interaction) CSV() []string {
	s := make([]string, 0, 15)

	s = append(s, pstr(i.Action))
	s = append(s, pstr(i.EntityType))
	s = append(s, pstr(i.EntityID))
	s = append(s, pstr(i.OriginType))
	s = append(s, pstr(i.OriginID))
	s = append(s, pstr(i.UserType))
	s = append(s, pstr(i.UserID))
	s = append(s, pstr(i.DeviceType))
	s = append(s, pstr(i.DeviceID))
	s = append(s, pstr(i.SessionType))
	s = append(s, pstr(i.SessionID))
	s = append(s, pstr(i.Timestamp))

	if i.CreatedAt != nil {
		s = append(s, i.CreatedAt.String())
	} else {
		s = append(s, "")
	}

	if i.ReceivedAt != nil {
		s = append(s, i.ReceivedAt.String())
	} else {
		s = append(s, "")
	}

	if i.Properties != nil {
		data, _ := json.Marshal(i.Properties)
		s = append(s, string(data))
	} else {
		s = append(s, "")
	}

	return s
}
