package types

import (
	"github.com/JKhawaja/errors"
	"github.com/gofrs/uuid"
)

// UnitMetrics is calculated from conversion event types
// and revenue is calculated from amount property on conversion event types.
// If there are no conversion event types then UnitMetrics will be zeroed.
// If there are no amount property keys then revenue results will be zeroed.
type UnitMetrics struct {
	TotalConversions      int64        `json:"totalConversions"`
	TotalRevenue          float64      `json:"totalRevenue"`
	AverageRevenuePerUser *float64     `json:"averageRevenuePerUser,omitempty"`
	AmountStats           *SimpleStats `json:"amountStats,omitempty"`
}

// ConversionStats --
type ConversionStats struct {
	Endpoint    uuid.UUID    `json:"endpoint"`
	TotalValue  float64      `json:"value"`
	AmountStats *SimpleStats `json:"amountStats"`
}

// NewConversionStats --
func NewConversionStats(endpoint uuid.UUID, amount float64) (*ConversionStats, error) {
	stats, err := NewSimpleStats(amount)
	if err != nil {
		return nil, errors.New(err, nil)
	}

	return &ConversionStats{
		Endpoint:    endpoint,
		TotalValue:  amount,
		AmountStats: stats,
	}, nil
}

// ConversionStatsList --
type ConversionStatsList struct {
	List []*ConversionStats
}

// NewConversionStatsList --
func NewConversionStatsList() *ConversionStatsList {
	return &ConversionStatsList{
		List: make([]*ConversionStats, 0),
	}
}

// NewUnitMetrics --
func NewUnitMetrics(i *Interaction) (*UnitMetrics, error) {
	if *i.Action != Conversion {
		return nil, nil
	}

	var amount float64
	if i.Properties != nil {
		a, ok := i.Properties["amount"]
		if ok {
			f, ok := a.(float64)
			if ok {
				amount = f
			}
		}
	}

	amountStats, err := NewSimpleStats(amount)
	if err != nil {
		return nil, errors.New(err, map[string]interface{}{
			"amount":    amount,
			"timestamp": i.CreatedAt,
		})
	}

	return &UnitMetrics{
		TotalConversions:      1,
		TotalRevenue:          amount,
		AverageRevenuePerUser: &amount,
		AmountStats:           amountStats,
	}, nil
}

// Update --
func (c *ConversionStats) Update(amount float64) error {
	c.TotalValue += amount
	err := c.AmountStats.Update(amount)
	if err != nil {
		return errors.New(err, nil)
	}

	return nil
}

// Update --
func (c *ConversionStatsList) Update(event *Event) error {
	i := event.Interaction
	if *i.Action == Conversion {
		if i.Properties != nil {
			a, ok := i.Properties["amount"]
			if ok {
				amount, ok := a.(float64)
				if ok {
					for _, stats := range c.List {
						if stats.Endpoint == event.Endpoint {
							err := stats.Update(amount)
							if err != nil {
								return errors.New(err, nil)
							}
							return nil
						}
					}

					newConversionStats, err := NewConversionStats(event.Endpoint, amount)
					if err != nil {
						return errors.New(err, nil)
					}
					c.List = append(c.List, newConversionStats)
				}
			}
		}
	}

	return nil
}

// SimpleUpdate --
func (u *UnitMetrics) SimpleUpdate(i *Interaction) {
	if *i.Action == Conversion {
		u.TotalConversions++

		a, ok := i.Properties["amount"]
		if ok {
			amount, ok := a.(float64)
			if ok {
				u.TotalRevenue += amount
				u.AmountStats.Update(amount)
			}
		}
	}
}

// SessionUpdate --
func (u *UnitMetrics) SessionUpdate(s *UserSession) error {
	u.TotalConversions += s.Conversions
	u.TotalRevenue += s.Value
	return u.AmountStats.Update(s.Value)
}

// Update --
func (u *UnitMetrics) Update(i *Interaction, users int64) error {
	if *i.Action == Conversion {
		u.TotalConversions++

		if i.Properties != nil {
			a, ok := i.Properties["amount"]
			if ok {
				amount, ok := a.(float64)
				if ok {
					u.TotalRevenue += amount
					arpu := u.TotalRevenue / float64(users)
					u.AverageRevenuePerUser = &arpu
					if u.AmountStats == nil {
						amountStats, err := NewSimpleStats(amount)
						if err != nil {
							return errors.New(err, nil)
						}
						u.AmountStats = amountStats
					} else {
						u.AmountStats.Update(amount)
					}
				}
			}
		}
	}

	return nil
}
