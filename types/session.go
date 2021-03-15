package types

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/JKhawaja/cache"
	"github.com/JKhawaja/errors"
)

var (
	// AutomatedSessionDetectionType is the session-type default name for when autoamted session detection is on
	AutomatedSessionDetectionType = "asd"
	// SessionExpiryDuration sets the basic maximum session period.
	SessionExpiryDuration = 1 * time.Hour
)

// Session hold the type and id of a session object
type Session struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// UserSession can track the latest session for a user
type UserSession struct {
	// session
	Type string `json:"type"`
	ID   string `json:"id"`

	// user
	UserType   string `json:"userType"`
	UserID     string `json:"userID"`
	DeviceType string `json:"deviceType"`
	DeviceID   string `json:"deviceID"`

	// interactions
	Total int64 `json:"total"`

	// conversions
	Conversions  int64         `json:"totalConversions"`
	Value        float64       `json:"value"` // total conversion amount for session
	OriginCounts *OriginCounts `json:"originCounts"`
	PrevEndpoint *Endpoint     `json:"prevEndpoint"`

	// current-origin
	VisitTotal     int           `json:"visit_total"`
	OriginDuration time.Duration `json:"origin_duration"`
	CurrentOrigin  *Origin       `json:"current_origin"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Device --
type Device struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// User represents the id and latest session for
// a user.
type User struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// UserSessions is a cache of current ongoing sessions
type UserSessions struct {
	*cache.Cache
}

// NewUserSessions --
func NewUserSessions(c *cache.Cache) *UserSessions {
	return &UserSessions{
		Cache: c,
	}
}

// Session --
func (s *UserSession) Session() Session {
	return Session{
		Type: s.Type,
		ID:   s.ID,
	}
}

// String --
func (s Session) String() string {
	return s.Type + "+" + s.ID
}

// GetSession will retrieve the session associated with the interaction.
// If no session is associated with the interaction, then a new session
// is started and that interaction is considered the first interaction of the
// session.
func (u *UserSessions) GetSession(i *Interaction) (*UserSession, error) {
	item, err := u.Get(i.User().String())
	if err == cache.ErrDNE {
		s := NewSession(i)
		err := u.Add(i.User().String(), s, 1*time.Hour)
		if err != nil {
			return s, errors.New(err, nil)
		}
		return s, nil
	} else if err != nil {
		return nil, errors.New(err, nil)
	}

	return item.(*UserSession), nil
}

// String will return a unique string which
// represents the User object.
func (u User) String() string {
	return u.Type + "," + u.ID
}

// User method returns a user object based on the interaction
// user type and user id.
func (i *Interaction) User() User {
	var ut, uid string

	if i.UserType != nil {
		ut = *i.UserType
	}

	if i.UserID != nil {
		uid = *i.UserID
	}

	return User{
		Type: ut,
		ID:   uid,
	}
}

// Session --
func (i *Interaction) Session() *Session {
	var st, sid string
	if i.SessionType != nil {
		st = *i.SessionType
	}
	if i.SessionID != nil {
		sid = *i.SessionID
	}
	return &Session{
		Type: st,
		ID:   sid,
	}
}

// Device --
func (i *Interaction) Device() Device {
	var dt, did string

	if i.DeviceType != nil {
		dt = *i.DeviceType
	}

	if i.DeviceID != nil {
		did = *i.DeviceID
	}

	return Device{
		Type: dt,
		ID:   did,
	}
}

// SetSession should be called to set the session-id for the interaction
// based on the automated session detection system.
func (u *UserSessions) SetSession(i *Interaction) error {
	session, err := u.GetSession(i)
	if err != nil {
		return errors.New(err, nil)
	}

	if i.SessionType == nil {
		i.SessionType = &AutomatedSessionDetectionType
	}

	if i.SessionID == nil {
		i.SessionID = &session.ID
	}

	return nil
}

// NewSession will create a new session object.
func NewSession(i *Interaction) *UserSession {
	sessType := AutomatedSessionDetectionType
	if i.SessionType != nil {
		sessType = *i.SessionType
	}

	user := i.User()
	device := i.Device()
	origin := i.Origin()
	origins := NewOriginCounts()
	origins.AddUnique(origin)

	return &UserSession{
		Type:          sessType,
		Total:         1,
		UserType:      user.Type,
		UserID:        user.ID,
		DeviceType:    device.Type,
		DeviceID:      device.ID,
		OriginCounts:  origins,
		CurrentOrigin: origin,
		ID:            NewUUID().String(),
		CreatedAt:     *i.CreatedAt,
		UpdatedAt:     *i.CreatedAt,
	}
}

// Renew will reset the values of the session using the information
// from the interaction.
// Note: the interaction will *not* be used to update the session count
// unless the Update() method is also called using this interaction.
func (s *UserSession) Renew(i *Interaction) {
	s.Total = 0
	s.ID = NewUUID().String()
	s.OriginCounts = NewOriginCounts()
	s.OriginDuration = time.Duration(0)
	s.CurrentOrigin = i.Origin()
	s.CreatedAt = *i.CreatedAt
	s.UpdatedAt = *i.CreatedAt
}

// Expired will check if the session has expired or not.
func (s *UserSession) Expired() bool {
	return time.Now().UTC().After(s.UpdatedAt.Add(SessionExpiryDuration))
}

// Bounced returns whether or not the session was a bounced session or not
func (s *UserSession) Bounced() bool {
	if s.UpdatedAt.Equal(s.CreatedAt) {
		return true
	}

	if len(s.OriginCounts.List) == 1 {
		return true
	}

	return false
}

// OriginChange will return if the interactions origin matches
// the session's current origin or not.
func (s *UserSession) OriginChange(i *Interaction) bool {
	return !s.CurrentOrigin.Eq(i.Origin())
}

// Origin will return the session's current origin and the current duration
// that the session has been at that origin.
func (s *UserSession) Origin(i *Interaction) (*Origin, time.Duration) {
	return s.CurrentOrigin, s.OriginDuration + (*i.CreatedAt).Sub(s.UpdatedAt)
}

// Update wil return the duration since the last interaction
// and whether or not the origin has changed.
// SessionUpdate should occur last when processing an interaction.
func (s *UserSession) Update(i *Interaction) {
	s.Total++
	defer func(s *UserSession, i *Interaction) {
		s.UpdatedAt = *i.CreatedAt
	}(s, i)

	// set if converted session or not
	if *i.Action == Conversion {
		s.Conversions++
	}

	s.PrevEndpoint = i.Endpoint()

	// same origin
	origin := i.Origin()
	if s.CurrentOrigin.Eq(origin) {
		s.OriginDuration += (*i.CreatedAt).Sub(s.UpdatedAt)
		s.VisitTotal++
		return
	} else {
		s.OriginCounts.IncrementVisit(origin)
	}

	if !s.OriginCounts.AddUnique(origin) {
		s.OriginCounts.Increment(origin)
	}

	s.CurrentOrigin = origin
	s.OriginDuration = 0
	s.VisitTotal = 1
}

// Duration will return the duration of the entire session
func (s *UserSession) Duration(expired bool) float64 {
	if expired {
		return s.UpdatedAt.Sub(s.CreatedAt).Minutes()
	}

	return time.Now().UTC().Sub(s.CreatedAt).Minutes()
}

// Encode will gop encode a session object into a
// slice of bytes.
// This will primarily be used for caching.
func (s *UserSession) Encode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(s.encoded())
	if err != nil {
		return []byte{}, errors.New(err, nil)
	}

	return buf.Bytes(), nil
}

// Decode will gob decode a slice of bytes into a
// session object.
// This will primarily be used for caching.
func (s *UserSession) Decode(data []byte) error {
	encoded := &sessionEncoded{}
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(encoded)
	if err != nil {
		return errors.New(err, nil)
	}

	s.decoded(encoded)

	return nil
}

// sessionEncoded is the version of a session that is the gob encoded standard
// it is used to translate a session into a gob-encoded byte slice and in reverse.
type sessionEncoded struct {
	Type         string         `json:"type"`
	ID           string         `json:"id"`
	UserType     string         `json:"userType"`
	UserID       string         `json:"userID"`
	Total        int64          `json:"total"`
	Conversions  int64          `json:"conversions"`
	Value        float64        `json:"value"`
	OriginCounts []*OriginCount `json:"originCounts"`
	PrevEndpoint *Endpoint      `json:"prevEndpoint"`

	// current-origin
	VisitTotal     int           `json:"visit_total"`
	OriginDuration time.Duration `json:"origin_duration"`
	CurrentOrigin  *Origin       `json:"current_origin"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *UserSession) encoded() *sessionEncoded {
	return &sessionEncoded{
		Type:           s.Type,
		ID:             s.ID,
		UserType:       s.UserType,
		UserID:         s.UserID,
		Total:          s.Total,
		Conversions:    s.Conversions,
		Value:          s.Value,
		OriginCounts:   s.OriginCounts.List,
		PrevEndpoint:   s.PrevEndpoint,
		VisitTotal:     s.VisitTotal,
		OriginDuration: s.OriginDuration,
		CurrentOrigin:  s.CurrentOrigin,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
	}
}

func (s *UserSession) decoded(encoded *sessionEncoded) {
	originCounts := NewOriginCounts()
	originCounts.List = encoded.OriginCounts

	s.Type = encoded.Type
	s.ID = encoded.ID
	s.UserType = encoded.UserType
	s.UserID = encoded.UserID
	s.Total = encoded.Total
	s.Conversions = encoded.Conversions
	s.Value = encoded.Value
	s.OriginCounts = originCounts
	s.PrevEndpoint = encoded.PrevEndpoint
	s.VisitTotal = encoded.VisitTotal
	s.OriginDuration = encoded.OriginDuration
	s.CurrentOrigin = encoded.CurrentOrigin
	s.CreatedAt = encoded.CreatedAt
	s.UpdatedAt = encoded.UpdatedAt
}
