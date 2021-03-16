package types

import (
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"github.com/JKhawaja/errors"
	"github.com/humilityai/temporal"
)

// EndpointStats --
type EndpointStats struct {
	ID       *UUID            `json:"id"`
	SpanType string           `json:"spanType"`
	Start    time.Time        `json:"start"`
	End      time.Time        `json:"end"`
	Profile  *EndpointProfile `json:"profile"`
}

// EndpointProfile --
type EndpointProfile struct {
	Total            int64                 `json:"total"`
	UserTypeStats    *SimpleStats          `json:"userTypeStats"`
	DeviceTypeStats  *SimpleStats          `json:"deviceTypeStats"`
	SessionTypeStats *SimpleStats          `json:"sessionTypeStats"`
	SessionStats     *SessionStatsList     `json:"sessionStats"` // tracks duration into session, and prior total interaction stats
	PropertyStats    *NamedSimpleStatsList `json:"propertyStats"`
}

// EndpointStatsList --
type EndpointStatsList struct {
	index   map[uint32]*EndpointStats
	updated map[uint32]*EndpointStats
	*sync.Mutex
}

// NewEndpointStats --
func NewEndpointStats(event *Event, spanType string) (*EndpointStats, error) {
	i := event.Interaction
	profile, err := NewEndpointProfile(i)
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
	case Quarterly:
		start = temporal.QuarterStart(*i.CreatedAt)
		end = temporal.QuarterFinish(*i.CreatedAt)
	case Yearly:
		start = temporal.YearStart(*i.CreatedAt)
		end = temporal.YearFinish(*i.CreatedAt)
	}

	return &EndpointStats{
		ID:       &UUID{UUID: event.Endpoint},
		Start:    start,
		End:      end,
		SpanType: spanType,
		Profile:  profile,
	}, nil
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

// Load --
func (e *EndpointStatsList) Load(list []*EndpointStats) error {
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

// Apply --
func (e *EndpointStats) Apply(event *Event) error {
	cat := *event.Interaction.CreatedAt
	if cat.After(e.End) {
		newStats, err := NewEndpointStats(event, e.SpanType)
		if err != nil {
			return errors.New(err, nil)
		}
		e = newStats
		return nil
	}

	return e.Profile.Apply(event.Interaction, event.Session)
}

// Apply --
func (e *EndpointStatsList) Apply(event *Event) error {
	e.Lock()
	defer e.Unlock()

	for _, spantype := range Intervals {
		hasher := fnv.New32a()
		_, err := hasher.Write([]byte(fmt.Sprintf("%s-%s", event.Endpoint.String(), spantype)))
		if err != nil {
			return errors.New(err, map[string]interface{}{
				"spantype": spantype,
				"id":       event.Endpoint,
			})
		}
		hashedKey := hasher.Sum32()

		stats, ok := e.index[hashedKey]
		if !ok {
			newStats, err := NewEndpointStats(event, spantype)
			if err != nil {
				return errors.New(err, map[string]interface{}{
					"spantype": spantype,
					"id":       event.Endpoint,
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
				"id":       event.Endpoint,
			})
		}
		e.updated[hashedKey] = stats
	}
	return nil
}

// Get --
func (e *EndpointStatsList) Get(id, spanType string) (*EndpointStats, error) {
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

// Update --
func (e *EndpointStatsList) Update(updateFunc func(object interface{}) error) error {
	e.Lock()
	defer e.Unlock()

	for id, mab := range e.updated {
		err := updateFunc(mab)
		if err != nil {
			return errors.New(err, nil)
		}
		delete(e.updated, id)
	}

	return nil
}

// Apply --
func (i *EndpointProfile) Apply(interaction *Interaction, sess *UserSession) error {
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
