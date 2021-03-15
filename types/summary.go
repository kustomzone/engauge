package types

import (
	"hash/fnv"
	"time"

	"github.com/JKhawaja/errors"
	"github.com/humilityai/temporal"
)

var (
	/* span type */

	// AllTime is a span type
	AllTime = "allTime"
	// Hourly is a span type
	Hourly = "hourly"
	// Daily is a span type
	Daily = "daily"
	// Weekly is a span type
	Weekly = "weekly"
	// Monthly is a span type
	Monthly = "monthly"
	// Spans is the full list of span types for summaries
	Spans = []string{AllTime, Hourly, Daily, Weekly, Monthly}
)

// Summary --
type Summary struct {
	SpanType         string
	Start            time.Time
	End              time.Time
	Total            int64
	ActionStats      *SimpleStats
	OriginTypeStats  *SimpleStats
	EntityTypeStats  *SimpleStats
	UserTypeStats    *SimpleStats
	DeviceTypeStats  *SimpleStats
	SessionTypeStats *SimpleStats
	Users            map[uint32]struct{}
	SessionStats     *SessionStatsList
	ConversionStats  *ConversionStatsList
	UnitMetrics      *UnitMetrics
}

// SummaryListView --
type SummaryListView struct {
	SpanType string    `json:"id"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
}

// SummaryResponse --
type SummaryResponse struct {
	ID               string             `json:"id"`
	Start            time.Time          `json:"start"`
	End              time.Time          `json:"end"`
	Total            int64              `json:"total"`
	ActionStats      *SimpleStats       `json:"actionStats,omitempty"`
	OriginTypeStats  *SimpleStats       `json:"originStats,omitempty"`
	EntityTypeStats  *SimpleStats       `json:"entityTypeStats,omitempty"`
	UserTypeStats    *SimpleStats       `json:"userTypeStats,omitempty"`
	DeviceTypeStats  *SimpleStats       `json:"deviceTypeStats,omitempty"`
	SessionTypeStats *SimpleStats       `json:"sessionTypeStats,omitempty"`
	SessionStats     []*SessionStats    `json:"sessionStats,omitempty"` // seu
	ConversionStats  []*ConversionStats `json:"conversionStats,omitempty"`
	UnitMetrics      *UnitMetrics       `json:"unitMetrics,omitempty"`
}

// NewSummary will generate a new summary for a specific spantype (daily, weekly, monthly)
// and for a specific interaction
func NewSummary(spanType string, event *Event) (*Summary, error) {
	i := event.Interaction
	sess := event.Session

	// span
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
		// all-time summary is a simplified summary object
		start = time.Time{}
		end = time.Unix(1<<63-1, 0)

		unitMetrics := &UnitMetrics{}
		unitMetrics.SimpleUpdate(i)

		sessionStats := NewSessionStatsList()

		return &Summary{
			Start:        start,
			End:          end,
			SpanType:     spanType,
			Total:        1,
			UnitMetrics:  unitMetrics,
			SessionStats: sessionStats,
		}, nil
	}

	actionStats, err := NewSimpleStats(*i.Action)
	if err != nil {
		return nil, errors.New(err, nil)
	}

	var originTypeStats, entityTypeStats, userTypeStats, deviceTypeStats, sessionTypeStats *SimpleStats
	if i.OriginType != nil {
		ets, err := NewSimpleStats(*i.OriginType)
		if err != nil {
			return nil, errors.New(err, nil)
		}
		originTypeStats = ets
	}

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

	// user
	users := make(map[uint32]struct{})

	hasher := fnv.New32a()
	_, err = hasher.Write([]byte(i.User().String()))
	if err != nil {
		return nil, errors.New(err, nil)
	}
	hashedKey := hasher.Sum32()
	users[hashedKey] = struct{}{}

	sessionStats := NewSessionStatsList()
	err = sessionStats.Update(sess)
	if err != nil {
		return nil, errors.New(err, nil)
	}

	// conversions
	conversionStats := NewConversionStatsList()
	err = conversionStats.Update(event)
	if err != nil {
		return nil, errors.New(err, nil)
	}

	unitMetrics := &UnitMetrics{}
	err = unitMetrics.Update(i, int64(len(users)))
	if err != nil {
		return nil, errors.New(err, nil)
	}

	return &Summary{
		Start:            start,
		End:              end,
		SpanType:         spanType,
		Total:            1,
		ActionStats:      actionStats,
		OriginTypeStats:  originTypeStats,
		EntityTypeStats:  entityTypeStats,
		UserTypeStats:    userTypeStats,
		DeviceTypeStats:  deviceTypeStats,
		SessionTypeStats: sessionTypeStats,
		Users:            users,
		SessionStats:     sessionStats,
		ConversionStats:  conversionStats,
		UnitMetrics:      unitMetrics,
	}, nil
}

// ListView --
func (s *Summary) ListView() *SummaryListView {
	return &SummaryListView{
		SpanType: s.SpanType,
		Start:    s.Start,
		End:      s.End,
	}
}

// Response --
func (s *Summary) Response() *SummaryResponse {
	return &SummaryResponse{
		ID:               s.SpanType,
		Start:            s.Start,
		End:              s.End,
		Total:            s.Total,
		ActionStats:      s.ActionStats,
		OriginTypeStats:  s.OriginTypeStats,
		EntityTypeStats:  s.EntityTypeStats,
		UserTypeStats:    s.UserTypeStats,
		DeviceTypeStats:  s.DeviceTypeStats,
		SessionTypeStats: s.SessionTypeStats,
		SessionStats:     s.SessionStats.List,
		ConversionStats:  s.ConversionStats.List,
		UnitMetrics:      s.UnitMetrics,
	}
}

// SessionExpirationUpdate --
func (s *Summary) SessionExpirationUpdate(session *UserSession) error {
	if s.SpanType == AllTime {
		return s.SessionStats.SimpleUpdate(session)
	}

	return s.SessionStats.Update(session)
}

// Expired will return whether or not the interaction is past the
// end time of the summary or not.
func (s *Summary) Expired(i *Interaction) bool {
	return (*i.CreatedAt).After(s.End)
}

// Apply --
func (s *Summary) Apply(event *Event) error {
	s.Total++

	i := event.Interaction

	err := s.ActionStats.Update(*i.Action)
	if err != nil {
		return errors.New(err, nil)
	}

	if i.OriginType != nil {
		if s.OriginTypeStats != nil {
			err := s.OriginTypeStats.Update(*i.OriginType)
			if err != nil {
				return errors.New(err, nil)
			}
		} else {
			ets, err := NewSimpleStats(*i.OriginType)
			if err != nil {
				return errors.New(err, nil)
			}
			s.OriginTypeStats = ets
		}
	}

	if i.EntityType != nil {
		if s.EntityTypeStats != nil {
			err := s.EntityTypeStats.Update(*i.EntityType)
			if err != nil {
				return errors.New(err, nil)
			}
		} else {
			ets, err := NewSimpleStats(*i.EntityType)
			if err != nil {
				return errors.New(err, nil)
			}
			s.EntityTypeStats = ets
		}
	}

	if i.UserType != nil {
		if s.UserTypeStats != nil {
			err := s.UserTypeStats.Update(*i.UserType)
			if err != nil {
				return errors.New(err, nil)
			}
		} else {
			ets, err := NewSimpleStats(*i.UserType)
			if err != nil {
				return errors.New(err, nil)
			}
			s.UserTypeStats = ets
		}
	}

	if i.DeviceType != nil {
		if s.DeviceTypeStats != nil {
			err := s.DeviceTypeStats.Update(*i.DeviceType)
			if err != nil {
				return errors.New(err, nil)
			}
		} else {
			ets, err := NewSimpleStats(*i.DeviceType)
			if err != nil {
				return errors.New(err, nil)
			}
			s.DeviceTypeStats = ets
		}
	}

	if i.SessionType != nil {
		if s.SessionTypeStats != nil {
			err := s.SessionTypeStats.Update(*i.SessionType)
			if err != nil {
				return errors.New(err, nil)
			}
		} else {
			ets, err := NewSimpleStats(*i.SessionType)
			if err != nil {
				return errors.New(err, nil)
			}
			s.SessionTypeStats = ets
		}
	}

	if s.SpanType == AllTime {
		s.UnitMetrics.SimpleUpdate(i)
		return nil
	}

	err = s.ConversionStats.Update(event)
	if err != nil {
		return err
	}

	if s.UnitMetrics == nil {
		s.UnitMetrics = &UnitMetrics{}
		err := s.UnitMetrics.Update(i, int64(len(s.Users)))
		if err != nil {
			return err
		}
	} else {
		err := s.UnitMetrics.Update(i, int64(len(s.Users)))
		if err != nil {
			return err
		}
	}

	return nil
}
