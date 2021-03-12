package types

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"github.com/JKhawaja/errors"
	"github.com/humilityai/temporal"
)

// OriginStats --
type OriginStats struct {
	ID       *UUID          `json:"id"`
	SpanType string         `json:"spantype"`
	Start    time.Time      `json:"start"`
	End      time.Time      `json:"end"`
	Profile  *OriginProfile `json:"profile"`
}

// OriginStatsList handles tracking origin stats by spantype via a hash key
// hash key = origin-id-string + spantype
type OriginStatsList struct {
	index   map[uint32]*OriginStats
	updated map[uint32]*OriginStats
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

// NewOriginStats --
func NewOriginStats(event *Event, spanType string) (*OriginStats, error) {
	i := event.Interaction
	profile, err := NewOriginProfile(event)
	if err != nil {
		return nil, errors.New(err, nil)
	}

	var start, end time.Time
	switch spanType {
	case Hourly:
		start = temporal.HourStart(*i.CreatedAt)
		end = temporal.HourFinish(*i.CreatedAt)
	case Daily:
		start = temporal.DayStart(*i.CreatedAt)
		end = temporal.DayFinish(*i.CreatedAt)
	case Weekly:
		start = temporal.WeekStart(*i.CreatedAt)
		end = temporal.WeekFinish(*i.CreatedAt)
	case Monthly:
		start = temporal.MonthStart(*i.CreatedAt)
		end = temporal.MonthFinish(*i.CreatedAt)
	case AllTime:
		start = time.Time{}
		end = time.Unix(1<<63-1, 0)
	}

	return &OriginStats{
		ID:       &UUID{UUID: event.Origin},
		SpanType: spanType,
		Start:    start,
		End:      end,
		Profile:  profile,
	}, nil
}

// NewOriginStatsList --
func NewOriginStatsList() *OriginStatsList {
	return &OriginStatsList{
		index:   make(map[uint32]*OriginStats),
		updated: make(map[uint32]*OriginStats),
		Mutex:   &sync.Mutex{},
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

// Apply --
func (o *OriginStats) Apply(event *Event) error {
	cat := *event.Interaction.CreatedAt
	if cat.After(o.End) && o.SpanType != AllTime {
		newStats, err := NewOriginStats(event, o.SpanType)
		if err != nil {
			return errors.New(err, nil)
		}
		o = newStats
		return nil
	}

	return o.Profile.Update(event)
}

// Get --
func (o *OriginStatsList) Get(id, spanType string) (*OriginStats, error) {
	hasher := fnv.New32a()
	_, err := hasher.Write([]byte(fmt.Sprintf("%s-%s", id, spanType)))
	if err != nil {
		return nil, errors.New(err, map[string]interface{}{
			"spantype": spanType,
			"id":       id,
		})
	}
	hashedKey := hasher.Sum32()

	stats, ok := o.index[hashedKey]
	if !ok {
		return nil, ErrDNE
	}

	return stats, nil
}

// Apply --
func (o *OriginStatsList) Apply(event *Event) error {
	o.Lock()
	defer o.Unlock()

	for _, spantype := range Spans {
		hasher := fnv.New32a()
		_, err := hasher.Write([]byte(fmt.Sprintf("%s-%s", event.Origin.String(), spantype)))
		if err != nil {
			return errors.New(err, map[string]interface{}{
				"spantype": spantype,
				"id":       event.Origin,
			})
		}
		hashedKey := hasher.Sum32()

		stats, ok := o.index[hashedKey]
		if !ok {
			newStats, err := NewOriginStats(event, spantype)
			if err != nil {
				return errors.New(err, map[string]interface{}{
					"spantype": spantype,
					"id":       event.Origin,
				})
			}
			o.index[hashedKey] = newStats
			o.updated[hashedKey] = newStats
			continue
		}

		err = stats.Apply(event)
		if err != nil {
			return errors.New(err, map[string]interface{}{
				"spantype": spantype,
				"id":       event.Origin,
			})
		}
		o.updated[hashedKey] = stats
	}
	return nil
}

// Load --
func (o *OriginStatsList) Load(list []*OriginStats) error {
	o.Lock()
	defer o.Unlock()

	for _, stats := range list {
		hasher := fnv.New32a()
		_, err := hasher.Write([]byte(fmt.Sprintf("%s-%s", stats.ID.UUID.String(), stats.SpanType)))
		if err != nil {
			return errors.New(err, map[string]interface{}{
				"spantype": stats.SpanType,
				"id":       stats.ID,
			})
		}
		hashedKey := hasher.Sum32()

		o.index[hashedKey] = stats
	}

	return nil
}

// Update --
func (o *OriginStatsList) Update(updateFunc func(object interface{}) error) error {
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

// VisitUpdate --
func (o *OriginProfile) VisitUpdate(s *UserSession) {
	o.VisitStats.Update(s)
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

// MarshalJSON --
func (o *OriginStats) MarshalJSON() ([]byte, error) {
	oCopy := struct {
		ID       *UUID          `json:"id"`
		SpanType string         `json:"spantype"`
		Start    time.Time      `json:"start"`
		End      time.Time      `json:"end"`
		Profile  *OriginProfile `json:"profile"`
	}{
		ID:       o.ID,
		SpanType: o.SpanType,
		Start:    o.Start,
		End:      o.End,
		Profile:  o.Profile,
	}

	if oCopy.SpanType == AllTime {
		oCopy.End = time.Time{}
	}

	return json.Marshal(oCopy)
}
