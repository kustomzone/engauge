package types

import (
	"strings"
	"sync"

	"github.com/JKhawaja/errors"
	"github.com/gofrs/uuid"
)

// Origin --
type Origin struct {
	ID         *UUID   `json:"id"`
	OriginType *string `json:"originType,omitempty"`
	OriginID   *string `json:"originID,omitempty"`
}

// Origins is a list of Origin objects.
type Origins struct {
	List    map[uuid.UUID]*Origin
	index   map[string]uuid.UUID
	updated map[uuid.UUID]*Origin
	*sync.Mutex
}

// OriginProfile --
type OriginProfile struct {
	Total            int64                 `json:"total"`
	ActionStats      *SimpleStats          `json:"actionStats,omitempty"`
	EntityTypeStats  *SimpleStats          `json:"entityTypeStats,omitempty"`
	UserTypeStats    *SimpleStats          `json:"userTypeStats,omitempty"`
	DeviceTypeStats  *SimpleStats          `json:"deviceTypeStats,omitempty"`
	SessionTypeStats *SimpleStats          `json:"sessionTypeStats,omitempty"`
	PropertyStats    *NamedSimpleStatsList `json:"propertyStats,omitempty"`
	VisitStats       *SessionStatsList     `json:"visitStats,omitempty"`
}

// OriginsList does not utilize a mutex and can be used for gob serialization
type OriginsList struct {
	List []*Origin
}

// OriginCount is primarily used inside of sessions
// to track how many interactions have occured at an
// origin, and how many visits to the origin have occurred
// within the session.
type OriginCount struct {
	Origin *Origin `json:"origin"`
	Count  int64   `json:"count"`
	Visits int64   `json:"visits"`
}

// OriginCounts --
type OriginCounts struct {
	List []*OriginCount
}

// OriginResponse --
type OriginResponse struct {
	ID         *UUID             `json:"id"`
	OriginType *string           `json:"originType,omitempty"`
	OriginID   *string           `json:"originID,omitempty"`
	Stats      *AllIntervalStats `json:"stats,omitempty"`
}

// NewOrigin will create and return a new origin object with a unique id
func NewOrigin(event *Event) *Origin {
	i := event.Interaction
	return &Origin{
		NewUUID(),
		i.OriginType,
		i.OriginID,
	}
}

// NewOrigins will return a pointer to an Origins object
// which manages a list of Origin object pointers.
func NewOrigins() *Origins {
	return &Origins{
		List:    make(map[uuid.UUID]*Origin),
		updated: make(map[uuid.UUID]*Origin),
		index:   make(map[string]uuid.UUID),
		Mutex:   &sync.Mutex{},
	}
}

// NewOriginsList --
func NewOriginsList() *OriginsList {
	return &OriginsList{
		List: make([]*Origin, 0),
	}
}

// NewOriginCounts --
func NewOriginCounts() *OriginCounts {
	return &OriginCounts{
		List: make([]*OriginCount, 0),
	}
}

// NewOriginProfile --
func NewOriginProfile(event *Event) (*OriginProfile, error) {
	i := event.Interaction
	sess := event.Session

	actionStats, err := NewSimpleStats(*i.Action)
	if err != nil {
		return nil, errors.New(err, nil)
	}

	var entityTypeStats, userTypeStats, deviceTypeStats, sessionTypeStats *SimpleStats
	if i.EntityType != nil {
		ets, err := NewSimpleStats(*i.EntityType)
		if err != nil {
			return nil, errors.New(err, nil)
		}
		entityTypeStats = ets
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

	propertyStats := NewNamedSimpleStatsList()
	if i.Properties != nil {
		for key, value := range i.Properties {
			err := propertyStats.Update(key, value)
			if err != nil {
				return nil, errors.New(err, nil)
			}
		}
	}

	endpointProfiles := NewEndpointsList()
	err = endpointProfiles.Apply(event)
	if err != nil {
		return nil, errors.New(err, nil)
	}

	visitStats := NewSessionStatsList()
	err = visitStats.Update(sess)
	if err != nil {
		return nil, errors.New(err, nil)
	}

	return &OriginProfile{
		Total:            1,
		ActionStats:      actionStats,
		EntityTypeStats:  entityTypeStats,
		UserTypeStats:    userTypeStats,
		DeviceTypeStats:  deviceTypeStats,
		SessionTypeStats: sessionTypeStats,
		PropertyStats:    propertyStats,
		VisitStats:       visitStats,
	}, nil
}

// Len --
func (o *Origins) Len() int {
	return len(o.List)
}

// Apply --
func (o *Origins) Apply(event *Event) {
	o.Lock()
	defer o.Unlock()

	_, ok := o.List[event.Origin]
	if !ok {
		// add if new origin
		newOrigin := NewOrigin(event)
		o.List[newOrigin.ID.UUID] = newOrigin
		o.index[newOrigin.String()] = newOrigin.ID.UUID
		o.updated[newOrigin.ID.UUID] = newOrigin

		// add new origin-id to event (for further event processing)
		event.Origin = newOrigin.ID.UUID
	}
}

// Contains --
func (o *Origins) Contains(rep string) bool {
	_, ok := o.index[rep]
	return ok
}

// ID --
func (o *Origins) ID(object interface{}) uuid.UUID {
	o.Lock()
	defer o.Unlock()

	origin, ok := object.(*Origin)
	if !ok {
		return UUIDNil.UUID
	}

	id, ok := o.index[origin.String()]
	if !ok {
		return UUIDNil.UUID
	}

	return id
}

// Update --
func (o *Origins) Update(updateFunc func(object interface{}) error) error {
	o.Lock()
	defer o.Unlock()

	for id, origin := range o.updated {
		err := updateFunc(origin)
		if err != nil {
			return errors.New(err, nil)
		}
		delete(o.updated, id)
	}

	return nil
}

// Remove --
func (o *Origins) Remove(key interface{}) error {
	o.Lock()
	defer o.Unlock()

	id, ok := key.(uuid.UUID)
	if ok {
		return ErrKeyType
	}

	delete(o.List, id)
	return nil
}

// Set --
func (o *Origins) Set(key, value interface{}) error {
	o.Lock()
	defer o.Unlock()

	id, ok := key.(uuid.UUID)
	if !ok {
		return errors.New(ErrKeyType, map[string]interface{}{
			"key":   key,
			"value": value,
		})
	}

	origin, ok := value.(*Origin)
	if !ok {
		return errors.New(ErrKeyType, map[string]interface{}{
			"key":   key,
			"value": value,
		})
	}

	o.List[id] = origin
	o.index[origin.String()] = id

	return nil
}

// Get --
func (o *Origins) Get(key interface{}) interface{} {
	o.Lock()
	defer o.Unlock()

	id, ok := key.(uuid.UUID)
	if !ok {
		return nil
	}

	origin, ok := o.List[id]
	if !ok {
		return nil
	}

	return origin
}

// Eq will check if an object is equal to this Origin object.
func (o *Origin) Eq(input interface{}) bool {
	o2, ok := input.(*Origin)
	if !ok {
		return false
	}

	if *o.OriginType != *o2.OriginType {
		return false
	}

	if *o.OriginID != *o2.OriginID {
		return false
	}

	return true
}

// Get will return a pointer to the OriginCount object and a boolean
// that specifies whether a matching OriginCount was found or not.
func (o *OriginCounts) Get(origin *Origin) (*OriginCount, bool) {
	for _, oc := range o.List {
		if oc.Origin.Eq(origin) {
			return oc, true
		}
	}

	return nil, false
}

// AddUnique will return whether or not the origin id
// was added to the list of origin counts.
func (o *OriginCounts) AddUnique(origin *Origin) bool {
	if !o.Contains(origin) {
		o.List = append(o.List, &OriginCount{
			Origin: origin,
			Count:  1,
		})
		return true
	}

	return false
}

// Increment --
func (o *OriginCounts) Increment(origin *Origin) {
	for _, oc := range o.List {
		if oc.Origin.Eq(origin) {
			oc.Count++
		}
	}
}

// IncrementVisit --
func (o *OriginCounts) IncrementVisit(origin *Origin) {
	for _, oc := range o.List {
		if oc.Origin.Eq(origin) {
			oc.Visits++
			return
		}
	}
}

// Contains is self-explanatory
func (o *OriginCounts) Contains(origin *Origin) bool {
	for _, oc := range o.List {
		if oc.Origin.Eq(origin) {
			return true
		}
	}

	return false
}

// String --
func (o *Origin) String() string {
	s := make([]string, 0, 2)

	s = append(s, pstr(o.OriginType))
	s = append(s, pstr(o.OriginID))

	return strings.Join(s, "-")
}

func (o *OriginsList) qualifies(id uuid.UUID) bool {
	for _, origin := range o.List {
		if origin.ID.UUID == id {
			return true
		}
	}

	return false
}

// Update --
func (o *OriginProfile) Update(event *Event) error {
	o.Total++

	i := event.Interaction
	s := event.Session

	err := o.ActionStats.Update(*i.Action)
	if err != nil {
		return errors.New(err, nil)
	}

	if i.EntityType != nil {
		if o.EntityTypeStats != nil {
			err := o.EntityTypeStats.Update(*i.EntityType)
			if err != nil {
				return errors.New(err, nil)
			}
		} else {
			ets, err := NewSimpleStats(*i.EntityType)
			if err != nil {
				return errors.New(err, nil)
			}
			o.EntityTypeStats = ets
		}
	}

	if i.UserType != nil {
		if o.UserTypeStats != nil {
			err := o.UserTypeStats.Update(*i.UserType)
			if err != nil {
				return errors.New(err, nil)
			}
		} else {
			ets, err := NewSimpleStats(*i.UserType)
			if err != nil {
				return errors.New(err, nil)
			}
			o.UserTypeStats = ets
		}
	}

	if i.DeviceType != nil {
		if o.DeviceTypeStats != nil {
			err := o.DeviceTypeStats.Update(*i.DeviceType)
			if err != nil {
				return errors.New(err, nil)
			}
		} else {
			ets, err := NewSimpleStats(*i.DeviceType)
			if err != nil {
				return errors.New(err, nil)
			}
			o.DeviceTypeStats = ets
		}
	}

	if i.SessionType != nil {
		if o.SessionTypeStats != nil {
			err := o.SessionTypeStats.Update(*i.SessionType)
			if err != nil {
				return errors.New(err, nil)
			}
		} else {
			ets, err := NewSimpleStats(*i.SessionType)
			if err != nil {
				return errors.New(err, nil)
			}
			o.SessionTypeStats = ets
		}
	}

	if i.Properties != nil {
		for name, value := range i.Properties {
			err := o.PropertyStats.Update(name, value)
			if err != nil {
				return errors.New(err, nil)
			}
		}
	}

	err = o.VisitStats.Update(s)
	if err != nil {
		return errors.New(err, nil)
	}

	return nil
}
