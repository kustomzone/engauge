package types

import (
	"bytes"
	"encoding/gob"

	"github.com/JKhawaja/errors"
)

const (
	/* toggle type */

	// OnOffToggle is a toggle type
	OnOffToggle = "onoff"
	// HistoryToggle is a toggle type
	HistoryToggle = "history"

	/* interval types */

	// AllTime is an interval type
	AllTime = "allTime"
	// Hourly is an interval type
	Hourly = "hourly"
	// Daily is an interval type
	Daily = "daily"
	// Weekly is an interval type
	Weekly = "weekly"
	// Monthly is an interval type
	Monthly = "monthly"
	// Quarterly is an interval type
	Quarterly = "quarterly"
	// Yearly is an interval type
	Yearly = "yearly"
)

var (
	// Intervals is the full list of interval types for summaries
	Intervals = []string{AllTime, Hourly, Daily, Weekly, Monthly, Quarterly, Yearly}
)

// Settings is the settings for the entire system
type Settings struct {
	ID                  string        `json:"id"`
	StatsToggles        *StatsToggles `json:"statsToggles"`
	InteractionsStorage bool          `json:"interactions"`
	User                string        `json:"-"`
	Password            string        `json:"-"`
	APIKey              string        `json:"-"`
	JWTSecret           string        `json:"-"`
}

// StatsToggles is the set of all summary toggles
type StatsToggles struct {
	Hourly    bool `json:"hourly"`
	Daily     bool `json:"daily"`
	Weekly    bool `json:"weekly"`
	Monthly   bool `json:"monthly"`
	Quarterly bool `json:"quarterly"`
	Yearly    bool `json:"yearly"`
}

// NewSettings --
func NewSettings() *Settings {
	return &Settings{
		StatsToggles:        NewStatsToggles(),
		InteractionsStorage: true,
	}
}

// NewStatsToggles will return a pointer to a new `StatsToggles` object.
func NewStatsToggles() *StatsToggles {
	return &StatsToggles{
		Hourly:    true,
		Daily:     true,
		Weekly:    true,
		Monthly:   true,
		Quarterly: true,
		Yearly:    true,
	}
}

// Update will toggle the value on a summary toggle field
func (s *StatsToggles) Update(spanType, toggleType string) {
	switch spanType {
	case Hourly:
		s.Hourly = toggleBool(s.Hourly)
	case Daily:
		s.Daily = toggleBool(s.Daily)
	case Weekly:
		s.Weekly = toggleBool(s.Weekly)
	case Monthly:
		s.Monthly = toggleBool(s.Monthly)
	case Quarterly:
		s.Quarterly = toggleBool(s.Quarterly)
	case Yearly:
		s.Yearly = toggleBool(s.Yearly)
	}
}

func toggleBool(value bool) bool {
	if value {
		return false
	}

	return true
}

// GobEncode --
func (s *Settings) GobEncode() ([]byte, error) {
	sCopy := struct {
		StatsToggles        *StatsToggles
		InteractionsStorage bool
	}{
		StatsToggles:        s.StatsToggles,
		InteractionsStorage: s.InteractionsStorage,
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(sCopy)
	if err != nil {
		return []byte{}, errors.New(err, nil)
	}

	return buf.Bytes(), nil
}

// GobDecode --
func (s *Settings) GobDecode(data []byte) error {
	type settings struct {
		StatsToggles           *StatsToggles
		InteractionsStorage    bool
		ConversionsStorageOnly bool
		InteractionsRetention  int
	}
	var sCopy *settings
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	err := dec.Decode(sCopy)
	if err != nil {
		return errors.New(err, nil)
	}

	s.StatsToggles = sCopy.StatsToggles
	s.InteractionsStorage = sCopy.InteractionsStorage
	return nil
}
