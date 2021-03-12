package types

import (
	"math"
	"sort"
	"time"

	"github.com/JKhawaja/errors"
	"github.com/humilityai/temporal"
)

// ValueStats --
type ValueStats struct {
	Value      interface{} `json:"value"`
	Count      int64       `json:"count"`
	Percentage float64     `json:"percentage"`

	// time stats
	Timezones   *SimpleStats         `json:"timezones"`
	Minutes     *SimpleStats         `json:"minutes"`
	Hours       *SimpleStats         `json:"hours"`
	Days        *SimpleStats         `json:"day"`
	Weekdays    *SimpleStats         `json:"weekdays"`
	Weeks       *SimpleStats         `json:"weeks"`
	Months      *SimpleStats         `json:"months"`
	Years       *SimpleStats         `json:"years"`
	TimeOfDay   *SimpleStats         `json:"timeOfDay"`
	Seasons     *SimpleStats         `json:"seasons"`
	Quarters    *SimpleStats         `json:"quarters"`
	YearSeasons temporal.YearSeasons `json:"-"`
}

// ValueStatsList --
type ValueStatsList []*ValueStats

// NewValueStats --
func NewValueStats(value interface{}, timestamp *time.Time) (*ValueStats, error) {
	y, w := timestamp.ISOWeek()
	yearSeasons := make(temporal.YearSeasons)

	tzStats, err := NewSimpleStats(timestamp.Location().String())
	if err != nil {
		return nil, errors.New(err, nil)
	}

	minuteStats, err := NewSimpleStats(timestamp.Minute())
	if err != nil {
		return nil, errors.New(err, nil)
	}

	hourStats, err := NewSimpleStats(timestamp.Hour())
	if err != nil {
		return nil, errors.New(err, nil)
	}

	dayStats, err := NewSimpleStats(timestamp.Day())
	if err != nil {
		return nil, errors.New(err, nil)
	}

	weekdayStats, err := NewSimpleStats(timestamp.Weekday().String())
	if err != nil {
		return nil, errors.New(err, nil)
	}

	weekStats, err := NewSimpleStats(w)
	if err != nil {
		return nil, errors.New(err, nil)
	}

	monthStats, err := NewSimpleStats(timestamp.Month().String())
	if err != nil {
		return nil, errors.New(err, nil)
	}

	yearStats, err := NewSimpleStats(y)
	if err != nil {
		return nil, errors.New(err, nil)
	}

	timeOfDayStats, err := NewSimpleStats(timeOfDay(*timestamp))
	if err != nil {
		return nil, errors.New(err, nil)
	}

	seasonStats, err := NewSimpleStats(yearSeasons.GetSeason(*timestamp).Name)
	if err != nil {
		return nil, errors.New(err, nil)
	}

	quarterStats, err := NewSimpleStats(temporal.Quarter(*timestamp))
	if err != nil {
		return nil, errors.New(err, nil)
	}

	return &ValueStats{
		Value:       value,
		Count:       1,
		Timezones:   tzStats,
		Minutes:     minuteStats,
		Hours:       hourStats,
		Days:        dayStats,
		Weekdays:    weekdayStats,
		Weeks:       weekStats,
		Months:      monthStats,
		Years:       yearStats,
		TimeOfDay:   timeOfDayStats,
		Seasons:     seasonStats,
		Quarters:    quarterStats,
		YearSeasons: yearSeasons,
	}, nil
}

// Update --
func (v *ValueStats) Update(timestamp *time.Time) {
	v.Count++

	if timestamp != nil {
		y, w := timestamp.ISOWeek()
		v.Timezones.Update(timestamp.Location().String())
		v.Minutes.Update(timestamp.Minute())
		v.Hours.Update(timestamp.Hour())
		v.Days.Update(timestamp.Day())
		v.Weekdays.Update(timestamp.Weekday().String())
		v.Weeks.Update(w)
		v.Months.Update(timestamp.Month().String())
		v.Years.Update(y)
		v.TimeOfDay.Update(temporal.TimeOfDay(*timestamp))
		v.Seasons.Update(v.YearSeasons.GetSeason(*timestamp).Name)
	}
}

func (s *Stats) get(idx int) *ValueStats {
	return s.Values[idx]
}

func (s *Stats) add(value interface{}, timestamp *time.Time) error {
	timeStats, err := NewValueStats(value, timestamp)
	if err != nil {
		return errors.New(err, nil)
	}

	s.Values = append(s.Values, timeStats)
	return nil
}

func (s *Stats) updateValue(value interface{}, timestamp *time.Time) error {
	// update value stats (if exists)
	var exists bool
	for _, stats := range s.Values {
		if stats.Value == value {
			stats.Update(timestamp)
			exists = true
			break
		}
	}

	// if dne: append new stats
	if !exists {
		valueStats, err := NewValueStats(value, timestamp)
		if err != nil {
			return errors.New(err, nil)
		}

		s.Values = append(s.Values, valueStats)
	}

	// set percentages
	sum := s.sum()
	for _, stats := range s.Values {
		stats.Percentage = float64(stats.Count) / float64(sum)
	}

	return nil
}

// returns the index of the count which represents the categorical "mean"
// of all the unique count objects.
func (s *Stats) meanIndex() int {
	sort.Sort(s)
	var sum, mean float64
	for i, stats := range s.Values {
		total := float64(stats.Count)
		sum += total
		mean += float64(i) * total
	}
	return int(math.Round(mean / sum))
}

func (s *Stats) max() (value interface{}) {
	max := int64(math.MinInt64)
	for _, stats := range s.Values {
		if stats.Count > max {
			max = stats.Count
			value = stats.Value
		}
	}
	return
}

func (s *Stats) sum() int64 {
	var sum int64
	for _, stats := range s.Values {
		sum += stats.Count
	}
	return sum
}

func (s *Stats) contains(value interface{}) bool {
	for _, stats := range s.Values {
		if stats.Value == value {
			return true
		}
	}

	return false
}

// Len ...
func (s *Stats) Len() int { return len(s.Values) }

// Swap ...
func (s *Stats) Swap(i, j int) { s.Values[i], s.Values[j] = s.Values[j], s.Values[i] }

// Less ....
func (s *Stats) Less(i, j int) bool { return s.Values[i].Count < s.Values[j].Count }
