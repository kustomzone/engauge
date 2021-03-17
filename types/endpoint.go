package types

import (
	"strings"
	"sync"

	"github.com/JKhawaja/errors"
	"github.com/gofrs/uuid"
)

// Endpoint represents the type of an interaction. It is defined by the fields for which a value exists
// and the specific value of that existing field.
// The set of all interaction types is equivalent to the set of all possibly valid selection queries.
// Selection queries will primarily vary by: timespans, users, sessions, and properties
type Endpoint struct {
	ID *UUID `json:"id"`

	// how
	Action *string `json:"action,omitempty"`

	// what
	EntityType *string `json:"entityType,omitempty"`
	EntityID   *string `json:"entityID,omitempty"`

	// where
	OriginType *string `json:"originType,omitempty"`
	OriginID   *string `json:"originID,omitempty"`
}

// EndpointListView --
type EndpointListView struct {
	ID *UUID `json:"id"`

	// how
	Action *string `json:"action,omitempty"`

	// what
	EntityType *string `json:"entityType,omitempty"`
	EntityID   *string `json:"entityID,omitempty"`

	// where
	OriginType *string `json:"originType,omitempty"`
	OriginID   *string `json:"originID,omitempty"`
}

// EndpointListViews --
type EndpointListViews []*EndpointListView

// EndpointResponse --
type EndpointResponse struct {
	ID *UUID `json:"id"`

	// how
	Action *string `json:"action,omitempty"`

	// what
	EntityType *string `json:"entityType,omitempty"`
	EntityID   *string `json:"entityID,omitempty"`

	// where
	OriginType *string `json:"originType,omitempty"`
	OriginID   *string `json:"originID,omitempty"`

	Stats *AllIntervalStats `json:"stats,omitempty"`
}

// Endpoints is a List of endpoint objects
type Endpoints struct {
	List    map[uuid.UUID]*Endpoint `json:"list"`
	index   map[string]uuid.UUID
	updated map[uuid.UUID]*Endpoint
	*sync.Mutex
}

// EndpointsList does not utilize a mutex and can be used for gob serialization
type EndpointsList struct {
	List []*Endpoint
}

// EndpointReward is the format for an Itype when sent over the API for viewing
type EndpointReward struct {
	ID     string   `json:"id"`
	Reward *float64 `json:"reward,omitempty"`
}

// EndpointsRewards is a list of EndpointReward objects
type EndpointsRewards []EndpointReward

// EndpointProfile --
type EndpointProfile struct {
	Total            int64                 `json:"total"`
	UserTypeStats    *SimpleStats          `json:"userTypeStats"`
	DeviceTypeStats  *SimpleStats          `json:"deviceTypeStats"`
	SessionTypeStats *SimpleStats          `json:"sessionTypeStats"`
	SessionStats     *SessionStatsList     `json:"sessionStats"` // tracks duration into session, and prior total interaction stats
	PropertyStats    *NamedSimpleStatsList `json:"propertyStats"`
}

// ListView --
func (e *Endpoint) ListView() *EndpointListView {
	return &EndpointListView{
		ID:         e.ID,
		Action:     e.Action,
		EntityType: e.EntityType,
		EntityID:   e.EntityID,
		OriginType: e.OriginType,
		OriginID:   e.OriginID,
	}
}

// NewEndpoints will create a new Endpoints List management object
func NewEndpoints() *Endpoints {
	return &Endpoints{
		List:    make(map[uuid.UUID]*Endpoint),
		index:   make(map[string]uuid.UUID),
		updated: make(map[uuid.UUID]*Endpoint),
		Mutex:   &sync.Mutex{},
	}
}

// NewEndpointsList --
func NewEndpointsList() *EndpointsList {
	return &EndpointsList{
		List: make([]*Endpoint, 0),
	}
}

// NewEndpoint will generate a new endpoint object (with a unique id) from an interaction
func NewEndpoint(i *Interaction) (*Endpoint, error) {
	endpoint := i.Endpoint()
	endpoint.ID = NewUUID()

	return endpoint, nil
}

// NewEndpointProfile --
func NewEndpointProfile(i *Interaction) (*EndpointProfile, error) {
	device := i.Device()
	user := i.User()
	session := i.Session()

	userTypeStats, err := NewSimpleStats(user.Type)
	if err != nil {
		return nil, errors.New(err, nil)
	}

	deviceTypeStats, err := NewSimpleStats(device.Type)
	if err != nil {
		return nil, errors.New(err, nil)
	}

	sessionTypeStats, err := NewSimpleStats(session.Type)
	if err != nil {
		return nil, errors.New(err, nil)
	}

	propertyStats := NewNamedSimpleStatsList()
	if i.Properties != nil {
		for name, value := range i.Properties {
			err := propertyStats.Update(name, value)
			if err != nil {
				return nil, errors.New(err, nil)
			}
		}
	}

	return &EndpointProfile{
		Total:            1,
		UserTypeStats:    userTypeStats,
		DeviceTypeStats:  deviceTypeStats,
		SessionTypeStats: sessionTypeStats,
		SessionStats:     NewSessionStatsList(),
		PropertyStats:    propertyStats,
	}, nil
}

// Response --
func (e *Endpoints) Response() []*Endpoint {
	var list []*Endpoint
	for _, ep := range e.List {
		list = append(list, ep)
	}
	return list
}

// ID --
func (i *Endpoints) ID(object interface{}) uuid.UUID {
	i.Lock()
	defer i.Unlock()

	endpoint, ok := object.(*Endpoint)
	if !ok {
		return UUIDNil.UUID
	}

	id, ok := i.index[endpoint.String()]
	if !ok {
		return UUIDNil.UUID
	}

	return id
}

// Range --
func (i *Endpoints) Range(rangeFunc func(key, value interface{}) bool) {
	i.Lock()
	defer i.Unlock()

	for key, it := range i.List {
		if !rangeFunc(key, it) {
			break
		}
	}
}

// Apply --
func (i *Endpoints) Apply(event *Event) error {
	i.Lock()
	defer i.Unlock()

	interaction := event.Interaction

	ep, ok := i.List[event.Endpoint]
	if ok {
		err := ep.Apply(event)
		if err != nil {
			return errors.New(err, nil)
		}

		i.updated[ep.ID.UUID] = ep
		return nil
	}

	ne, err := NewEndpoint(interaction)
	if err != nil {
		return errors.New(err, nil)
	}

	i.List[ne.ID.UUID] = ne
	i.index[ne.String()] = ne.ID.UUID
	i.updated[ne.ID.UUID] = ne

	event.Endpoint = ne.ID.UUID

	return nil
}

// Apply --
func (e *EndpointsList) Apply(event *Event) error {
	interaction := event.Interaction

	for _, ep := range e.List {
		if ep.ID.UUID == event.Endpoint {
			err := ep.Apply(event)
			if err != nil {
				return errors.New(err, nil)
			}
			return nil
		}
	}

	ne, err := NewEndpoint(interaction)
	if err != nil {
		return errors.New(err, nil)
	}

	e.List = append(e.List, ne)

	return nil
}

// Update --
func (i *Endpoints) Update(updateFunc func(object interface{}) error) error {
	i.Lock()
	defer i.Unlock()

	for id, endpoint := range i.updated {
		err := updateFunc(endpoint)
		if err != nil {
			return errors.New(err, nil)
		}
		delete(i.updated, id)
	}

	return nil
}

// Remove --
func (i *Endpoints) Remove(key interface{}) error {
	i.Lock()
	defer i.Unlock()

	id, ok := key.(uuid.UUID)
	if !ok {
		return ErrKeyType
	}

	delete(i.List, id)
	return nil
}

// Set --
func (i *Endpoints) Set(key, value interface{}) error {
	i.Lock()
	defer i.Unlock()

	id, ok := key.(uuid.UUID)
	if !ok {
		return errors.New(ErrKeyType, map[string]interface{}{
			"key":   key,
			"value": value,
		})
	}

	endpoint, ok := value.(*Endpoint)
	if !ok {
		return errors.New(ErrValueType, map[string]interface{}{
			"key":   key,
			"value": value,
		})
	}

	i.List[id] = endpoint
	i.index[endpoint.String()] = id
	return nil
}

// Get --
func (i *Endpoints) Get(key interface{}) interface{} {
	i.Lock()
	defer i.Unlock()

	id, ok := key.(uuid.UUID)
	if !ok {
		return nil
	}

	ep, ok := i.List[id]
	if !ok {
		return nil
	}

	return ep
}

// Apply --
func (i *Endpoint) Apply(event *Event) error {
	return nil
}

// SuperTypes will return whether or not the endpoint argument is a sub-type.
// A subtype must match on the Action field, but it only has to match on other
// fields if those fields have a value specified by the supertype.
func (i *Endpoint) SuperTypes(sub *Endpoint) bool {
	super := i

	if *super.Action != *sub.Action {
		return false
	}

	if super.OriginType != nil {
		if sub.OriginType == nil {
			return false
		}

		if *super.OriginType != *sub.OriginType {
			return false
		}
	}

	if super.OriginID != nil {
		if sub.OriginID == nil {
			return false
		}

		if *super.OriginID != *sub.OriginID {
			return false
		}
	}

	if super.EntityType != nil {
		if sub.EntityType == nil {
			return false
		}

		if *super.EntityType != *sub.EntityType {
			return false
		}
	}

	if super.EntityID != nil {
		if sub.EntityID == nil {
			return false
		}

		if *super.EntityID != *sub.EntityID {
			return false
		}
	}

	return true
}

// Len will return the length of the list of Endpoints
func (i *Endpoints) Len() int {
	return len(i.List)
}

// Eq checks if the IType is equivalent to the provided object argument.
// This function works for any arbitrary argument without error or panic.
func (i *Endpoint) Eq(input interface{}) bool {
	i2, ok := input.(*Endpoint)
	if !ok {
		return false
	}

	if i.Action != nil && i2.Action != nil {
		if *i.Action != *i2.Action {
			return false
		}
	} else {
		return false
	}

	if i.EntityType != nil && i2.EntityType != nil {
		if *i.EntityType != *i2.EntityType {
			return false
		}
	} else {
		return false
	}

	if i.EntityID != nil && i2.EntityID != nil {
		if *i.EntityID != *i2.EntityID {
			return false
		}
	} else {
		return false
	}

	if i.OriginType != nil && i2.OriginType != nil {
		if *i.OriginType != *i2.OriginType {
			return false
		}
	} else {
		return false
	}

	if i.OriginID != nil && i2.OriginID != nil {
		if *i.OriginID != *i2.OriginID {
			return false
		}
	} else {
		return false
	}

	return true
}

// Apply --
func (i *EndpointProfile) Update(event *Event) error {
	interaction := event.Interaction
	sess := event.Session

	i.Total++

	if interaction.UserType != nil {
		if i.UserTypeStats != nil {
			err := i.UserTypeStats.Update(*interaction.UserType)
			if err != nil {
				return errors.New(err, nil)
			}
		} else {
			ets, err := NewSimpleStats(*interaction.UserType)
			if err != nil {
				return errors.New(err, nil)
			}
			i.UserTypeStats = ets
		}
	}

	if interaction.DeviceType != nil {
		if i.DeviceTypeStats != nil {
			err := i.DeviceTypeStats.Update(*interaction.DeviceType)
			if err != nil {
				return errors.New(err, nil)
			}
		} else {
			ets, err := NewSimpleStats(*interaction.DeviceType)
			if err != nil {
				return errors.New(err, nil)
			}
			i.DeviceTypeStats = ets
		}
	}

	if interaction.SessionType != nil {
		if i.SessionTypeStats != nil {
			err := i.SessionTypeStats.Update(*interaction.SessionType)
			if err != nil {
				return errors.New(err, nil)
			}
		} else {
			ets, err := NewSimpleStats(*interaction.SessionType)
			if err != nil {
				return errors.New(err, nil)
			}
			i.SessionTypeStats = ets
		}
	}

	if interaction.Properties != nil {
		for name, value := range interaction.Properties {
			err := i.PropertyStats.Update(name, value)
			if err != nil {
				return errors.New(err, nil)
			}
		}
	}

	err := i.SessionStats.Update(sess)
	if err != nil {
		return errors.New(err, nil)
	}
	return nil
}

// String --
func (i *Endpoint) String() string {
	s := make([]string, 0, 5)

	s = append(s, pstr(i.Action))
	s = append(s, pstr(i.EntityType))
	s = append(s, pstr(i.EntityID))
	s = append(s, pstr(i.OriginType))
	s = append(s, pstr(i.OriginID))

	return strings.Join(s, "-")
}

func pstr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
