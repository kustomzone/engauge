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
	ID             *UUID        `json:"id"`
	EntityType     *string      `json:"entityType,omitempty"`
	EntityID       *string      `json:"entityID,omitempty"`
	AlltimeStats   *EntityStats `json:"allTimeStats,omitempty"`
	HourlyStats    *EntityStats `json:"hourlyStats,omitempty"`
	DailyStats     *EntityStats `json:"dailyStats,omitempty"`
	WeeklyStats    *EntityStats `json:"weeklyStats,omitempty"`
	MonthlyStats   *EntityStats `json:"monthlyStats,omitempty"`
	QuarterlyStats *EntityStats `json:"quarterlyStats,omitempty"`
	YearlyStats    *EntityStats `json:"yearlyStats,omitempty"`
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
