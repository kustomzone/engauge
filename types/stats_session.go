package types

import "github.com/JKhawaja/errors"

// SessionStats --
type SessionStats struct {
	UserType        string       `json:"userType"`
	DeviceType      string       `json:"deviceType"`
	SessionType     string       `json:"sessionType"`
	Count           int64        `json:"count"`
	Percentage      float64      `json:"percentage"`
	Duration        *SimpleStats `json:"durationStats"`
	Interactions    *SimpleStats `json:"interactionStats"`
	Conversions     int64        `json:"conversions"`
	ConversionRate  float64      `json:"conversionRate"`
	BouncedSessions int64        `json:"bouncedSessions"`
	BounceRate      float64      `json:"bounceRate"`
}

// SessionStatsList --
type SessionStatsList struct {
	List []*SessionStats
}

// NewSessionStatsList --
func NewSessionStatsList() *SessionStatsList {
	return &SessionStatsList{
		List: make([]*SessionStats, 0),
	}
}

// NewSessionStats --
func NewSessionStats(sess *UserSession) (*SessionStats, error) {
	durationStats, err := NewSimpleStats(sess.Duration(sess.Expired()))
	if err != nil {
		return nil, errors.New(err, nil)
	}

	interactionStats, err := NewSimpleStats(sess.Total)
	if err != nil {
		return nil, errors.New(err, nil)
	}

	var bouncedSessions int64
	if sess.Bounced() {
		bouncedSessions++
	}

	return &SessionStats{
		UserType:        sess.UserType,
		DeviceType:      sess.DeviceType,
		SessionType:     sess.Type,
		Count:           1,
		Duration:        durationStats,
		Interactions:    interactionStats,
		BouncedSessions: bouncedSessions,
		BounceRate:      float64(bouncedSessions) / 1.0,
	}, nil
}

// Update --
func (u *SessionStats) Update(sess *UserSession) {
	u.Count++
	u.Duration.Update(sess.Duration(sess.Expired()))
	u.Interactions.Update(sess.Total)
	u.Conversions += sess.Conversions
	if sess.Bounced() {
		u.BouncedSessions++
	}
	u.BounceRate = float64(u.BouncedSessions) / float64(u.Count)
}

// SimpleUpdate --
func (u *SessionStats) SimpleUpdate(sess *UserSession) {
	u.Count++
	u.Conversions += sess.Conversions
	if sess.Bounced() {
		u.BouncedSessions++
	}
	u.BounceRate = float64(u.BouncedSessions) / float64(u.Count)
}

// SimpleUpdate --
func (s *SessionStatsList) SimpleUpdate(sess *UserSession) error {
	// update stats if exists
	for _, stats := range s.List {
		if stats.UserType == sess.UserType && stats.SessionType == sess.Type && sess.DeviceType == sess.DeviceType {
			stats.SimpleUpdate(sess)
			return nil
		}
	}

	// else: create new stats
	newStats, err := NewSessionStats(sess)
	if err != nil {
		return errors.New(err, nil)
	}
	s.List = append(s.List, newStats)

	return nil
}

// Update --
func (u *SessionStatsList) Update(sess *UserSession) error {
	// update stats if exists
	for _, stats := range u.List {
		if stats.UserType == sess.UserType && stats.SessionType == sess.Type && stats.DeviceType == sess.DeviceType {
			stats.Update(sess)
			return nil
		}
	}

	// else: create new stats
	newSession, err := NewSessionStats(sess)
	if err != nil {
		return errors.New(err, nil)
	}
	u.List = append(u.List, newSession)

	return nil
}

// Count --
func (u *SessionStatsList) Count(usertype, sessiontype, devicetype string) int64 {
	for _, stats := range u.List {
		if stats.UserType == usertype && stats.SessionType == sessiontype && stats.DeviceType == devicetype {
			return stats.Count
		}
	}

	return 0.0
}
