package types

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"github.com/JKhawaja/errors"
	"github.com/gofrs/uuid"
	"github.com/humilityai/temporal"
)

const (
	OriginObjectType   = "origin"
	EndpointObjectType = "endpoint"
	EntityObjectType   = "entity"
)

// Updater --
type Updater interface {
	Update(event *Event) error
}

// IntervalStats --
type IntervalStats struct {
	ID       *UUID     `json:"id"`
	Interval string    `json:"interval"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Stats    Updater   `json:"stats"`
}

// IntervalStatsList --
type IntervalStatsList struct {
	Type    string
	index   map[uint32]*IntervalStats
	updated map[uint32]*IntervalStats
	*sync.Mutex
}

// AllIntervalStats --
type AllIntervalStats struct {
	Alltime   *IntervalStats `json:"allTime,omitempty"`
	Hourly    *IntervalStats `json:"hourly,omitempty"`
	Daily     *IntervalStats `json:"daily,omitempty"`
	Weekly    *IntervalStats `json:"weekly,omitempty"`
	Monthly   *IntervalStats `json:"monthly,omitempty"`
	Quarterly *IntervalStats `json:"quarterly,omitempty"`
	Yearly    *IntervalStats `json:"yearly,omitempty"`
}

// NewIntervalStatsList --
func NewIntervalStatsList(objectType string) *IntervalStatsList {
	return &IntervalStatsList{
		Type:    objectType,
		index:   make(map[uint32]*IntervalStats),
		updated: make(map[uint32]*IntervalStats),
		Mutex:   &sync.Mutex{},
	}
}

// NewIntervalStats --
func NewIntervalStats(event *Event, interval string, objectType string) (*IntervalStats, error) {
	var updater Updater
	var id uuid.UUID
	switch objectType {
	case OriginObjectType:
		prof, err := NewOriginProfile(event)
		if err != nil {
			return nil, errors.New(err, nil)
		}
		updater = prof
		id = event.Origin
	case EndpointObjectType:
		prof, err := NewEndpointProfile(event.Interaction)
		if err != nil {
			return nil, errors.New(err, nil)
		}
		updater = prof
		id = event.Endpoint
	case EntityObjectType:
		prof, err := NewEntityProfile(event.Interaction)
		if err != nil {
			return nil, errors.New(err, nil)
		}
		updater = prof
		id = event.Entity
	}

	i := event.Interaction

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

	return &IntervalStats{
		ID:       &UUID{id},
		Interval: interval,
		Start:    start,
		End:      end,
		Stats:    updater,
	}, nil
}

// Get --
func (i *IntervalStatsList) Get(id, interval string) (*IntervalStats, error) {
	hasher := fnv.New32a()
	_, err := hasher.Write([]byte(fmt.Sprintf("%s-%s", id, interval)))
	if err != nil {
		return nil, errors.New(err, map[string]interface{}{
			"interval": interval,
			"id":       id,
		})
	}
	hashedKey := hasher.Sum32()

	stats, ok := i.index[hashedKey]
	if !ok {
		return nil, ErrDNE
	}

	return stats, nil
}

// Apply --
func (i *IntervalStatsList) Apply(event *Event) error {
	i.Lock()
	defer i.Unlock()

	var id uuid.UUID
	switch i.Type {
	case OriginObjectType:
		id = event.Origin
	case EndpointObjectType:
		id = event.Endpoint
	case EntityObjectType:
		id = event.Entity
	}

	for _, interval := range Intervals {
		hasher := fnv.New32a()
		_, err := hasher.Write([]byte(fmt.Sprintf("%s-%s", id.String(), interval)))
		if err != nil {
			return errors.New(err, map[string]interface{}{
				"interval": interval,
				"id":       id,
			})
		}
		hashedKey := hasher.Sum32()

		stats, ok := i.index[hashedKey]
		if !ok {
			newStats, err := NewIntervalStats(event, interval, i.Type)
			if err != nil {
				return errors.New(err, map[string]interface{}{
					"interval": interval,
					"id":       id,
				})
			}

			i.index[hashedKey] = newStats
			i.updated[hashedKey] = newStats
			continue
		}

		// update
		cat := *event.Interaction.CreatedAt
		if cat.After(stats.End) && stats.Interval != AllTime {
			newStats, err := NewIntervalStats(event, stats.Interval, i.Type)
			if err != nil {
				return errors.New(err, nil)
			}
			stats = newStats
		} else {
			err = stats.Stats.Update(event)
			if err != nil {
				return errors.New(err, map[string]interface{}{
					"interval": interval,
					"id":       id,
				})
			}
		}
		i.updated[hashedKey] = stats
	}
	return nil
}

// Load --
func (i *IntervalStatsList) Load(list []*IntervalStats) error {
	i.Lock()
	defer i.Unlock()

	for _, stats := range list {
		hasher := fnv.New32a()
		_, err := hasher.Write([]byte(fmt.Sprintf("%s-%s", stats.ID.UUID.String(), stats.Interval)))
		if err != nil {
			return errors.New(err, map[string]interface{}{
				"interval": stats.Interval,
				"id":       stats.ID,
			})
		}
		hashedKey := hasher.Sum32()

		i.index[hashedKey] = stats
	}

	return nil
}

// Update --
func (i *IntervalStatsList) Update(updateFunc func(object interface{}) error) error {
	i.Lock()
	defer i.Unlock()

	for id, es := range i.updated {
		err := updateFunc(es)
		if err != nil {
			return errors.New(err, nil)
		}
		delete(i.updated, id)
	}

	return nil
}

// AllIntervalStats --
func (i *IntervalStatsList) AllIntervalStats(id string) (*AllIntervalStats, error) {
	ais := &AllIntervalStats{}

	for _, interval := range Intervals {
		switch interval {
		case Hourly:
			stats, err := i.Get(id, interval)
			if err != nil {
				return ais, errors.New(err, nil)
			}

			if time.Now().After(stats.End) {
				stats = &IntervalStats{}
			}

			ais.Hourly = stats
		case Daily:
			stats, err := i.Get(id, interval)
			if err != nil {
				return ais, errors.New(err, nil)
			}

			if time.Now().After(stats.End) {
				stats = &IntervalStats{}
			}

			ais.Daily = stats
		case Weekly:
			stats, err := i.Get(id, interval)
			if err != nil {
				return ais, errors.New(err, nil)
			}

			if time.Now().After(stats.End) {
				stats = &IntervalStats{}
			}

			ais.Weekly = stats
		case Monthly:
			stats, err := i.Get(id, interval)
			if err != nil {
				return ais, errors.New(err, nil)
			}

			if time.Now().After(stats.End) {
				stats = &IntervalStats{}
			}

			ais.Monthly = stats

		case Quarterly:
			stats, err := i.Get(id, interval)
			if err != nil {
				return ais, errors.New(err, nil)
			}

			if time.Now().After(stats.End) {
				stats = &IntervalStats{}
			}

			ais.Quarterly = stats
		case Yearly:
			stats, err := i.Get(id, interval)
			if err != nil {
				return ais, errors.New(err, nil)
			}

			if time.Now().After(stats.End) {
				stats = &IntervalStats{}
			}

			ais.Yearly = stats
		case AllTime:
			stats, err := i.Get(id, interval)
			if err != nil {
				return ais, errors.New(err, nil)
			}

			ais.Alltime = stats
		}
	}

	return ais, nil
}

// MarshalJSON --
func (i *IntervalStats) MarshalJSON() ([]byte, error) {
	eCopy := struct {
		ID       *UUID       `json:"id"`
		Interval string      `json:"interval"`
		Start    time.Time   `json:"start"`
		End      time.Time   `json:"end"`
		Stats    interface{} `json:"stats"`
	}{
		ID:       i.ID,
		Interval: i.Interval,
		Start:    i.Start,
		End:      i.End,
		Stats:    i.Stats,
	}

	if eCopy.Interval == AllTime {
		eCopy.End = time.Time{}
	}

	return json.Marshal(eCopy)
}
