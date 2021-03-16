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

// EntityStats --
type EntityStats struct {
	ID       *UUID          `json:"id"`
	SpanType string         `json:"spantype"`
	Start    time.Time      `json:"start"`
	End      time.Time      `json:"end"`
	Profile  *EntityProfile `json:"profile"`
}

// EntityStatsList handles tracking entity stats by spantype via a hash key
// hash key = entity-id-string + spantype
type EntityStatsList struct {
	index   map[uint32]*EntityStats
	updated map[uint32]*EntityStats
	*sync.Mutex
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

// NewEntityStatsList --
func NewEntityStatsList() *EntityStatsList {
	return &EntityStatsList{
		index:   make(map[uint32]*EntityStats),
		updated: make(map[uint32]*EntityStats),
		Mutex:   &sync.Mutex{},
	}
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

// NewEntityStats --
func NewEntityStats(event *Event, interval string) (*EntityStats, error) {
	i := event.Interaction
	profile, err := NewEntityProfile(i)
	if err != nil {
		return nil, errors.New(err, nil)
	}

	var start, end time.Time
	switch interval {
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
	case Quarterly:
		start = temporal.QuarterStart(*i.CreatedAt)
		end = temporal.QuarterFinish(*i.CreatedAt)
	case Yearly:
		start = temporal.YearStart(*i.CreatedAt)
		end = temporal.YearFinish(*i.CreatedAt)
	case AllTime:
		start = time.Time{}
		end = time.Unix(1<<63-1, 0)
	}

	return &EntityStats{
		ID:       &UUID{UUID: event.Entity},
		SpanType: interval,
		Start:    start,
		End:      end,
		Profile:  profile,
	}, nil
}

// Update --
func (e *EntityProfile) Update(i *Interaction) error {
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

// Apply --
func (e *EntityStats) Apply(event *Event) error {
	cat := *event.Interaction.CreatedAt
	if cat.After(e.End) && e.SpanType != AllTime {
		newStats, err := NewEntityStats(event, e.SpanType)
		if err != nil {
			return errors.New(err, nil)
		}
		e = newStats
		return nil
	}

	return e.Profile.Update(event.Interaction)
}

// Get --
func (e *EntityStatsList) Get(id, spanType string) (*EntityStats, error) {
	hasher := fnv.New32a()
	_, err := hasher.Write([]byte(fmt.Sprintf("%s-%s", id, spanType)))
	if err != nil {
		return nil, errors.New(err, map[string]interface{}{
			"spantype": spanType,
			"id":       id,
		})
	}
	hashedKey := hasher.Sum32()

	stats, ok := e.index[hashedKey]
	if !ok {
		return nil, ErrDNE
	}

	return stats, nil
}

// Apply --
func (e *EntityStatsList) Apply(event *Event) error {
	e.Lock()
	defer e.Unlock()

	for _, spantype := range Intervals {
		hasher := fnv.New32a()
		_, err := hasher.Write([]byte(fmt.Sprintf("%s-%s", event.Entity.String(), spantype)))
		if err != nil {
			return errors.New(err, map[string]interface{}{
				"spantype": spantype,
				"id":       event.Entity,
			})
		}
		hashedKey := hasher.Sum32()

		stats, ok := e.index[hashedKey]
		if !ok {
			newStats, err := NewEntityStats(event, spantype)
			if err != nil {
				return errors.New(err, map[string]interface{}{
					"spantype": spantype,
					"id":       event.Entity,
				})
			}
			e.index[hashedKey] = newStats
			e.updated[hashedKey] = newStats
			continue
		}

		err = stats.Apply(event)
		if err != nil {
			return errors.New(err, map[string]interface{}{
				"spantype": spantype,
				"id":       event.Entity,
			})
		}
		e.updated[hashedKey] = stats
	}
	return nil
}

// Load --
func (e *EntityStatsList) Load(list []*EntityStats) error {
	e.Lock()
	defer e.Unlock()

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

		e.index[hashedKey] = stats
	}

	return nil
}

// Update --
func (e *EntityStatsList) Update(updateFunc func(object interface{}) error) error {
	e.Lock()
	defer e.Unlock()

	for id, es := range e.updated {
		err := updateFunc(es)
		if err != nil {
			return errors.New(err, nil)
		}
		delete(e.updated, id)
	}

	return nil
}

// MarshalJSON --
func (e *EntityStats) MarshalJSON() ([]byte, error) {
	eCopy := struct {
		ID       *UUID          `json:"id"`
		SpanType string         `json:"spantype"`
		Start    time.Time      `json:"start"`
		End      time.Time      `json:"end"`
		Profile  *EntityProfile `json:"profile"`
	}{
		ID:       e.ID,
		SpanType: e.SpanType,
		Start:    e.Start,
		End:      e.End,
		Profile:  e.Profile,
	}

	if eCopy.SpanType == AllTime {
		eCopy.End = time.Time{}
	}

	return json.Marshal(eCopy)
}
