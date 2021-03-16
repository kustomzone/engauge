package types

import (
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"github.com/JKhawaja/errors"
	"github.com/humilityai/temporal"
)

var (
	/* property types */

	// String is a property type
	String = "string"
	// Number is a property type
	Number = "number"
	// NumberArray is a property type
	NumberArray = "number-array"
	// StringArray is a property type
	StringArray = "string-array"

	// Duration is a value type
	Duration = "duration"
)

// Property holds the information for a property
type Property struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Stats *Stats `json:"stats"`
}

// PropertyResponse --
type PropertyResponse struct {
	ID             string         `json:"id"`
	Type           string         `json:"type"`
	Stats          *Stats         `json:"stats"`
	HourlyStats    *PropertyStats `json:"hourlyStats,omitempty"`
	DailyStats     *PropertyStats `json:"dailyStats,omitempty"`
	WeeklyStats    *PropertyStats `json:"weeklyStats,omitempty"`
	MonthlyStats   *PropertyStats `json:"monthlyStats,omitempty"`
	QuarterlyStats *PropertyStats `json:"quarterlyStats,omitempty"`
	YearlyStats    *PropertyStats `json:"yearlyStats,omitempty"`
}

// PropertyListView --
type PropertyListView struct {
	Name string `json:"id"`
	Type string `json:"type"`
}

// PropertyListViews --
type PropertyListViews []*PropertyListView

// ListView --
func (p *Property) ListView() *PropertyListView {
	return &PropertyListView{
		Name: p.Name,
		Type: p.Type,
	}
}

// PropertyStats --
type PropertyStats struct {
	Name     string       `json:"name"`
	SpanType string       `json:"spanType"`
	Start    time.Time    `json:"start"`
	End      time.Time    `json:"end"`
	Stats    *SimpleStats `json:"stats"`
}

// PropertyStatsList --
type PropertyStatsList struct {
	index   map[uint32]*PropertyStats
	updated map[uint32]*PropertyStats
	*sync.Mutex
}

// Properties --
type Properties struct {
	List    map[string]*Property
	updated map[string]*Property
	*sync.Mutex
}

// NewProperty --
func NewProperty(name string, value interface{}, timestamp *time.Time) (*Property, error) {
	propType := PropertyType(value)
	propStats, err := NewStats(value, timestamp)
	if err != nil {
		return nil, errors.New(err, nil)
	}

	return &Property{
		Name:  name,
		Type:  propType,
		Stats: propStats,
	}, nil
}

// NewProperties --
func NewProperties() *Properties {
	return &Properties{
		List:    make(map[string]*Property),
		updated: make(map[string]*Property),
		Mutex:   &sync.Mutex{},
	}
}

// NewPropertyStatsList --
func NewPropertyStatsList() *PropertyStatsList {
	return &PropertyStatsList{
		index:   make(map[uint32]*PropertyStats),
		updated: make(map[uint32]*PropertyStats),
		Mutex:   &sync.Mutex{},
	}
}

// NewPropertyStats --
func NewPropertyStats(name string, spanType string, value interface{}, timestamp *time.Time) (*PropertyStats, error) {
	stats, err := NewSimpleStats(value)
	if err != nil {
		return nil, errors.New(err, nil)
	}

	var start, end time.Time
	switch spanType {
	case Hourly:
		start = temporal.HourStart(*timestamp)
		end = temporal.HourFinish(*timestamp)
	case Daily:
		start = temporal.DayStart(*timestamp)
		end = temporal.DayFinish(*timestamp)
	case Weekly:
		start = temporal.WeekStart(*timestamp)
		end = temporal.WeekFinish(*timestamp)
	case Monthly:
		start = temporal.MonthStart(*timestamp)
		end = temporal.MonthFinish(*timestamp)
	case Quarterly:
		start = temporal.QuarterStart(*timestamp)
		end = temporal.QuarterFinish(*timestamp)
	case Yearly:
		start = temporal.YearStart(*timestamp)
		end = temporal.YearFinish(*timestamp)
	}

	return &PropertyStats{
		Name:     name,
		SpanType: spanType,
		Start:    start,
		End:      end,
		Stats:    stats,
	}, nil
}

// Len --
func (p *Properties) Len() int {
	return len(p.List)
}

// Response --
func (p *Property) Response() *PropertyResponse {
	return &PropertyResponse{
		ID:    p.Name,
		Type:  p.Type,
		Stats: p.Stats,
	}
}

// Apply --
func (p *Properties) Apply(event *Event) error {
	p.Lock()
	defer p.Unlock()
	i := event.Interaction

	if i.Properties != nil {
		for name, value := range i.Properties {
			prop, ok := p.List[name]
			if !ok {
				newProp, err := NewProperty(name, value, i.CreatedAt)
				if err != nil {
					return errors.New(err, nil)
				}
				p.List[name] = newProp
			} else {
				err := prop.Apply(value, i.CreatedAt)
				if err != nil {
					return errors.New(err, nil)
				}
			}

			p.updated[name] = p.List[name]
		}
	}

	return nil
}

// Update --
func (p *Properties) Update(updateFunc func(object interface{}) error) error {
	p.Lock()
	defer p.Unlock()

	for key, property := range p.updated {
		err := updateFunc(property)
		if err != nil {
			return errors.New(err, nil)
		}
		delete(p.updated, key)
	}

	return nil
}

// Remove --
func (p *Properties) Remove(key interface{}) error {
	p.Lock()
	defer p.Unlock()

	name, ok := key.(string)
	if !ok {
		return ErrKeyType
	}

	delete(p.List, name)
	return nil
}

// Set --
func (p *Properties) Set(key, value interface{}) error {
	p.Lock()
	defer p.Unlock()

	property, ok := value.(*Property)
	if !ok {
		return ErrValueType
	}

	name, ok := key.(string)
	if !ok {
		return ErrKeyType
	}

	p.List[name] = property
	return nil
}

// Get --
func (p *Properties) Get(key interface{}) interface{} {
	p.Lock()
	defer p.Unlock()

	name, ok := key.(string)
	if !ok {
		return nil
	}

	prop, ok := p.List[name]
	if !ok {
		return nil
	}

	return prop
}

// Apply --
func (p *Property) Apply(value interface{}, timestamp *time.Time) error {
	err := p.Stats.Update(value, timestamp)
	if err != nil {
		return errors.New(err, nil)
	}
	return nil
}

// PropertyType will return the property type of the given value
func PropertyType(v interface{}) string {
	switch v.(type) {
	case string:
		return String
	case float64, int:
		return Number
	case []string:
		return StringArray
	case []float64, []int:
		return NumberArray
	}

	return ""
}

// ValidPropertyType returns whether or not the type name is a valid
// property type name or not.
func ValidPropertyType(t string) bool {
	switch t {
	case String:
		return true
	case Number:
		return true
	case NumberArray:
		return true
	case StringArray:
		return true
	default:
		return false
	}
}

// Apply --
func (p *PropertyStats) Apply(value interface{}, timestamp *time.Time) error {
	t := *timestamp
	if t.After(p.End) {
		newStats, err := NewPropertyStats(p.Name, p.SpanType, value, timestamp)
		if err != nil {
			return errors.New(err, nil)
		}
		p = newStats
		return nil
	}

	return p.Stats.Update(value)
}

// Apply --
func (p *PropertyStatsList) Apply(event *Event) error {
	p.Lock()
	defer p.Unlock()

	i := event.Interaction

	if i.Properties != nil {
		for name, value := range i.Properties {
			for _, spantype := range Intervals {
				hasher := fnv.New32a()
				_, err := hasher.Write([]byte(fmt.Sprintf("%s-%s", name, spantype)))
				if err != nil {
					return errors.New(err, map[string]interface{}{
						"spantype": spantype,
						"name":     name,
					})
				}
				hashedKey := hasher.Sum32()

				stats, ok := p.index[hashedKey]
				if !ok {
					newStats, err := NewPropertyStats(name, spantype, value, i.CreatedAt)
					if err != nil {
						return errors.New(err, map[string]interface{}{
							"spantype": spantype,
							"name":     name,
						})
					}
					p.index[hashedKey] = newStats
					p.updated[hashedKey] = newStats
					continue
				}

				err = stats.Apply(value, i.CreatedAt)
				if err != nil {
					return errors.New(err, map[string]interface{}{
						"spantype": spantype,
						"name":     name,
					})
				}
				p.updated[hashedKey] = stats
			}
		}
	}

	return nil
}

// Load --
func (p *PropertyStatsList) Load(list []*PropertyStats) error {
	p.Lock()
	defer p.Unlock()

	for _, stats := range list {
		hasher := fnv.New32a()
		_, err := hasher.Write([]byte(fmt.Sprintf("%s-%s", stats.Name, stats.SpanType)))
		if err != nil {
			return errors.New(err, map[string]interface{}{
				"spantype": stats.SpanType,
				"name":     stats.Name,
			})
		}
		hashedKey := hasher.Sum32()

		p.index[hashedKey] = stats
	}

	return nil
}

// Get --
func (p *PropertyStatsList) Get(name, spanType string) (*PropertyStats, error) {
	hasher := fnv.New32a()
	_, err := hasher.Write([]byte(fmt.Sprintf("%s-%s", name, spanType)))
	if err != nil {
		return nil, errors.New(err, map[string]interface{}{
			"spantype": spanType,
			"name":     name,
		})
	}
	hashedKey := hasher.Sum32()

	stats, ok := p.index[hashedKey]
	if !ok {
		return nil, ErrDNE
	}

	return stats, nil
}

// Update --
func (p *PropertyStatsList) Update(updateFunc func(object interface{}) error) error {
	p.Lock()
	defer p.Unlock()

	for id, mab := range p.updated {
		err := updateFunc(mab)
		if err != nil {
			return errors.New(err, nil)
		}
		delete(p.updated, id)
	}

	return nil
}
