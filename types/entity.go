package types

import (
	"strings"
	"sync"

	"github.com/JKhawaja/errors"
	"github.com/gofrs/uuid"
)

// Entity --
type Entity struct {
	ID         *UUID  `json:"id"`
	EntityType string `json:"entityType"`
	EntityID   string `json:"entityID"`
}

// Entities --
type Entities struct {
	List    map[uuid.UUID]*Entity
	index   map[string]uuid.UUID
	updated map[uuid.UUID]*Entity
	*sync.Mutex
}

// EntityResponse --
type EntityResponse struct {
	ID         *UUID             `json:"id"`
	EntityType *string           `json:"entityType,omitempty"`
	EntityID   *string           `json:"entityID,omitempty"`
	Stats      *AllIntervalStats `json:"stats,omitempty"`
}

// EntityProfile --
type EntityProfile struct {
	Total            int64                 `json:"total"`
	ActionStats      *SimpleStats          `json:"actionStats,omitempty"`
	UserTypeStats    *SimpleStats          `json:"userTypeStats,omitempty"`
	DeviceTypeStats  *SimpleStats          `json:"deviceTypeStats,omitempty"`
	SessionTypeStats *SimpleStats          `json:"sessionTypeStats,omitempty"`
	PropertyStats    *NamedSimpleStatsList `json:"propertyStats,omitempty"`
}

// NewEntities --
func NewEntities() *Entities {
	return &Entities{
		List:    make(map[uuid.UUID]*Entity),
		index:   make(map[string]uuid.UUID),
		updated: make(map[uuid.UUID]*Entity),
		Mutex:   &sync.Mutex{},
	}
}

// NewEntity --
func NewEntity(i *Interaction) *Entity {
	entity := i.Entity()
	entity.ID = NewUUID()
	return entity
}

// NewEntityProfile --
func NewEntityProfile(i *Interaction) (*EntityProfile, error) {
	var userTypeStats, deviceTypeStats, sessionTypeStats *SimpleStats

	actionStats, err := NewSimpleStats(*i.Action)
	if err != nil {
		return nil, errors.New(err, nil)
	}

	if i.UserType != nil {
		uts, err := NewSimpleStats(*i.UserType)
		if err != nil {
			return nil, errors.New(err, nil)
		}
		userTypeStats = uts
	}

	if i.DeviceType != nil {
		dts, err := NewSimpleStats(*i.DeviceType)
		if err != nil {
			return nil, errors.New(err, nil)
		}
		deviceTypeStats = dts
	}

	if i.SessionType != nil {
		sts, err := NewSimpleStats(*i.SessionType)
		if err != nil {
			return nil, errors.New(err, nil)
		}
		sessionTypeStats = sts
	}

	// properties
	propertyStats := NewNamedSimpleStatsList()
	if i.Properties != nil {
		for key, value := range i.Properties {
			err := propertyStats.Update(key, value)
			if err != nil {
				return nil, errors.New(err, nil)
			}
		}
	}

	return &EntityProfile{
		Total:            1,
		UserTypeStats:    userTypeStats,
		ActionStats:      actionStats,
		DeviceTypeStats:  deviceTypeStats,
		SessionTypeStats: sessionTypeStats,
		PropertyStats:    propertyStats,
	}, nil
}

// Apply --
func (e *Entities) Apply(event *Event) {
	e.Lock()
	defer e.Unlock()

	_, ok := e.List[event.Entity]
	if !ok {
		// add if new origin
		newEntity := NewEntity(event.Interaction)
		e.List[newEntity.ID.UUID] = newEntity
		e.index[newEntity.String()] = newEntity.ID.UUID
		e.updated[newEntity.ID.UUID] = newEntity

		// add new entity-id to event (for further event processing)
		event.Entity = newEntity.ID.UUID
	}
}

// Update --
func (e *Entities) Update(updateFunc func(object interface{}) error) error {
	e.Lock()
	defer e.Unlock()

	for id, entity := range e.updated {
		err := updateFunc(entity)
		if err != nil {
			return errors.New(err, nil)
		}
		delete(e.updated, id)
	}

	return nil
}

// ID --
func (e *Entities) ID(object interface{}) uuid.UUID {
	e.Lock()
	defer e.Unlock()

	entity, ok := object.(*Entity)
	if !ok {
		return UUIDNil.UUID
	}

	id, ok := e.index[entity.String()]
	if !ok {
		return UUIDNil.UUID
	}

	return id
}

// Set --
func (e *Entities) Set(entity *Entity) {
	e.Lock()
	defer e.Unlock()

	e.List[entity.ID.UUID] = entity
	e.index[entity.String()] = entity.ID.UUID
}

// Len --
func (e *Entities) Len() int {
	return len(e.List)
}

// String --
func (e *Entity) String() string {
	s := make([]string, 0, 2)

	s = append(s, e.EntityType)
	s = append(s, e.EntityID)

	return strings.Join(s, "-")
}

// Update --
func (e *EntityProfile) Update(event *Event) error {
	i := event.Interaction

	e.Total++

	err := e.ActionStats.Update(*i.Action)
	if err != nil {
		return errors.New(err, nil)
	}

	if i.UserType != nil {
		if e.UserTypeStats != nil {
			err := e.UserTypeStats.Update(*i.UserType)
			if err != nil {
				return errors.New(err, nil)
			}
		} else {
			ets, err := NewSimpleStats(*i.UserType)
			if err != nil {
				return errors.New(err, nil)
			}
			e.UserTypeStats = ets
		}
	}

	if i.DeviceType != nil {
		if e.DeviceTypeStats != nil {
			err := e.DeviceTypeStats.Update(*i.DeviceType)
			if err != nil {
				return errors.New(err, nil)
			}
		} else {
			ets, err := NewSimpleStats(*i.DeviceType)
			if err != nil {
				return errors.New(err, nil)
			}
			e.DeviceTypeStats = ets
		}
	}

	if i.SessionType != nil {
		if e.SessionTypeStats != nil {
			err := e.SessionTypeStats.Update(*i.SessionType)
			if err != nil {
				return errors.New(err, nil)
			}
		} else {
			ets, err := NewSimpleStats(*i.SessionType)
			if err != nil {
				return errors.New(err, nil)
			}
			e.SessionTypeStats = ets
		}
	}

	if i.Properties != nil {
		for name, value := range i.Properties {
			err := e.PropertyStats.Update(name, value)
			if err != nil {
				return errors.New(err, nil)
			}
		}
	}

	return nil
}
