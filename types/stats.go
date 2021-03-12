package types

import (
	"math"
	"time"

	"github.com/JKhawaja/errors"
	"github.com/humilityai/sam"
	"gonum.org/v1/gonum/stat/distuv"
)

/*
	The variances/standard-deviations on categorical/multinomial distributions represent the reasonable
	amount of expected variation in counts for each unique value in any given arbitrary sample of values.
*/

// NamedStats --
type NamedStats struct {
	Name  string `json:"name"`
	Stats *Stats `json:"stats"`
}

// NamedSimpleStats --
type NamedSimpleStats struct {
	Name  string       `json:"name"`
	Stats *SimpleStats `json:"stats"`
}

// NamedSimpleStatsList --
type NamedSimpleStatsList struct {
	List []*NamedSimpleStats
}

// NamedStatsList --
type NamedStatsList struct {
	List []*NamedStats
}

// Stats holds statistical data for a sample.
type Stats struct {
	Type     string           `json:"type"`
	Total    int64            `json:"total"`
	Mean     interface{}      `json:"mean"`
	Mode     interface{}      `json:"mode"`
	Values   ValueStatsList   `json:"values"`
	Variance sam.SliceFloat64 `json:"variance"`
	StdDev   sam.SliceFloat64 `json:"stdDev"`
}

// NewNamedStatsList --
func NewNamedStatsList() *NamedStatsList {
	return &NamedStatsList{
		List: make([]*NamedStats, 0),
	}
}

// NewNamedSimpleStatsList --
func NewNamedSimpleStatsList() *NamedSimpleStatsList {
	return &NamedSimpleStatsList{
		List: make([]*NamedSimpleStats, 0),
	}
}

// NewNamedStats --
func NewNamedStats(name string, value interface{}, timestamp *time.Time) (*NamedStats, error) {
	stats, err := NewStats(value, timestamp)
	if err != nil {
		return nil, errors.New(err, nil)
	}

	return &NamedStats{
		Name:  name,
		Stats: stats,
	}, nil
}

// NewNamedSimpleStats --
func NewNamedSimpleStats(name string, value interface{}) (*NamedSimpleStats, error) {
	stats, err := NewSimpleStats(value)
	if err != nil {
		return nil, errors.New(err, nil)
	}

	return &NamedSimpleStats{
		Name:  name,
		Stats: stats,
	}, nil
}

// Update --
func (n *NamedSimpleStatsList) Update(name string, value interface{}) error {
	var exists bool
	for _, ns := range n.List {
		if ns.Name == name {
			ns.Stats.Update(value)
			exists = true
			break
		}
	}

	if !exists {
		newStats, err := NewNamedSimpleStats(name, value)
		if err != nil {
			return errors.New(err, nil)
		}

		n.List = append(n.List, newStats)
	}

	return nil
}

// Update will update the named stats for the provided name using
// the new given value. If the named stats does not already exist
// then a new statistics object will be created.
func (s *NamedStatsList) Update(name string, value interface{}, timestamp *time.Time) error {
	var exists bool
	for _, ns := range s.List {
		if ns.Name == name {
			ns.Stats.Update(value, timestamp)
			exists = true
			break
		}
	}

	if !exists {
		newStats, err := NewNamedStats(name, value, timestamp)
		if err != nil {
			return errors.New(err, nil)
		}

		s.List = append(s.List, newStats)
	}

	return nil
}

// NewStats generates a new statistical object dependent on the provided type
func NewStats(value interface{}, timestamp *time.Time) (*Stats, error) {
	s := &Stats{}

	switch value.(type) {
	case string:
		v := value.(string)
		s.Type = String
		s.Total = 1
		s.Mean = v
		s.Mode = v
		s.Values = make([]*ValueStats, 0)
		s.add(value, timestamp)
		s.Variance = make(sam.SliceFloat64, 1, 1)
		s.StdDev = make(sam.SliceFloat64, 1, 1)
	case float64, int, int64:
		var v float64
		switch value.(type) {
		case float64:
			v = value.(float64)
		case int:
			iv := value.(int)
			v = float64(iv)
		case int64:
			iv := value.(int64)
			v = float64(iv)
		}

		s.Type = Number
		s.Total = 1
		s.Mean = v
		s.Mode = v
		s.Values = make([]*ValueStats, 0)
		s.add(value, timestamp)
		s.Variance = make(sam.SliceFloat64, 1, 1)
		s.StdDev = make(sam.SliceFloat64, 1, 1)
	case time.Duration:
		v := value.(time.Duration)
		s.Type = Duration
		s.Total = 1
		s.Mean = v
		s.Mode = v
		s.Values = make([]*ValueStats, 0)
		s.add(value, timestamp)
		s.Variance = make(sam.SliceFloat64, 1, 1)
		s.StdDev = make(sam.SliceFloat64, 1, 1)
	case []float64, []int:
		var v []float64
		switch value.(type) {
		case float64:
			v = value.([]float64)
		case int:
			iv := value.([]int)
			for _, val := range iv {
				v = append(v, float64(val))
			}
		}

		s.Type = NumberArray
		s.Total = 1
		s.Mean = v
		s.Mode = v
		s.Values = make([]*ValueStats, 0)
		for _, val := range v {
			s.add(val, timestamp)
		}
		s.Variance = make(sam.SliceFloat64, 1, 1)
		s.StdDev = make(sam.SliceFloat64, 1, 1)
	case []string:
		v := value.([]string)
		s.Type = StringArray
		s.Total = 1
		s.Mean = v
		s.Mode = v
		s.Values = make([]*ValueStats, 0)
		for _, val := range v {
			s.add(val, timestamp)
		}
		s.Variance = make(sam.SliceFloat64, len(s.Values), len(s.Values))
		s.StdDev = make(sam.SliceFloat64, len(s.Values), len(s.Values))
	default:
		return s, errors.New(ErrDataType, map[string]interface{}{
			"value": value,
		})
	}

	return s, nil
}

// Update will update the total and mean for the stream stats
// with the new value.
func (s *Stats) Update(value interface{}, timestamp *time.Time) error {
	switch value.(type) {
	case string:
		if s.Type != String {
			return errors.New(ErrDataType, nil)
		}

		v := value.(string)

		// if value is a new value then increase variance & stddev array size
		if !s.contains(v) {
			s.Variance = append(s.Variance, 0)
			s.StdDev = append(s.StdDev, 0)
		}

		s.updateValue(v, timestamp)

		// update stats
		s.Total++
		s.Mean = s.get(s.meanIndex()).Value
		s.Mode = s.max()
		sum := s.sum()
		for i := 0; i < s.Len(); i++ {
			val := s.get(i).Count
			prob := float64(val) / float64(sum)
			variance := float64(val) * (1 - prob)
			s.Variance[i] = variance
			s.StdDev[i] = math.Sqrt(variance)
		}
	case float64, int:
		if s.Type != Number {
			return errors.New(ErrDataType, nil)
		}

		var v float64
		switch value.(type) {
		case float64:
			v = value.(float64)
		case int:
			iv := value.(int)
			v = float64(iv)
		}

		oldMean := s.Mean.(float64)
		s.Mean = (float64(s.Total)*oldMean + v) / (float64(s.Total) + 1)

		s.Total++
		s.Variance[0] = ((float64(s.Total)-2)*s.Variance[0] + (v-s.Mean.(float64))*(v-oldMean)) / (float64(s.Total) - 1)
		s.StdDev[0] = math.Sqrt(s.Variance[0])

		s.updateValue(value, timestamp)
	case []string:
		if s.Type != StringArray {
			return ErrDataType
		}

		v := value.([]string)

		for _, str := range v {
			if !s.contains(str) {
				s.Variance = append(s.Variance, 0)
				s.StdDev = append(s.StdDev, 0)
			}

			// increment
			s.updateValue(str, timestamp)

			// update stats
			s.Total++
			s.Mean = s.get(s.meanIndex()).Value
			s.Mode = s.max()
			sum := s.sum()
			for i := 0; i < s.Len(); i++ {
				val := s.get(i).Count
				prob := float64(val) / float64(sum)
				variance := float64(val) * (1 - prob)
				s.Variance[i] = variance
				s.StdDev[i] = math.Sqrt(variance)
			}
		}
	case []float64, []int:
		if s.Type != NumberArray {
			return errors.New(ErrDataType, nil)
		}

		var v []float64
		switch value.(type) {
		case []float64:
			v = value.([]float64)
		case []int:
			iv := value.([]int)
			for _, val := range iv {
				v = append(v, float64(val))
			}
		}

		for _, f := range v {
			oldMean := s.Mean.(float64)
			s.Mean = (float64(s.Total)*oldMean + f) / (float64(s.Total) + 1)

			s.Total++
			s.Variance[0] = ((float64(s.Total)-2)*s.Variance[0] + (f-s.Mean.(float64))*(f-oldMean)) / (float64(s.Total) - 1)
			s.StdDev[0] = math.Sqrt(s.Variance[0])

			s.updateValue(f, timestamp)
		}
	case time.Duration:
		if s.Type != Duration {
			return errors.New(ErrDataType, nil)
		}
		v := float64(value.(time.Duration))

		oldMean, ok := s.Mean.(float64)
		if !ok {
			return errors.New(ErrAssertion, map[string]interface{}{
				"mean":  s.Mean,
				"value": value,
				"total": s.Total,
			})
		}
		s.Mean = (float64(s.Total)*oldMean + v) / (float64(s.Total) + 1)

		s.Total++
		s.Variance[0] = ((float64(s.Total)-2)*s.Variance[0] + (v-s.Mean.(float64))*(v-oldMean)) / (float64(s.Total) - 1)
		s.StdDev[0] = math.Sqrt(s.Variance[0])

		s.updateValue(value, timestamp)

	default:
		return errors.New(ErrDataType, nil)
	}

	return nil
}

func timeOfDay(timestamp time.Time) string {
	h := timestamp.Hour()
	switch {
	case h > 2 && h < 6:
		return "late night" // 3, 4, 5
	case h > 5 && h < 9:
		return "early morning" // 6, 7, 8
	case h > 8 && h < 12:
		return "late morning" // 9, 10, 11
	case h > 11 && h < 15:
		return "early afternoon" // 12, 13, 14
	case h > 14 && h < 18:
		return "late afternoon" // 15, 16, 17
	case h > 17 && h < 21:
		return "early evening" // 18, 19, 20
	case h > 20 && h < 24:
		return "late evening" // 21, 22, 23
	default:
		return "early night" // 24, 1, 2
	}
}

func pValue(rewards, counts sam.SliceFloat64) float64 {
	totalScore := rewards.Sum()
	var se0, se1, seDiff, zScore float64
	if counts[0] > 0 {
		se0 = math.Sqrt((rewards[0] * (totalScore - rewards[0])) / counts[0])
	}
	if counts[1] > 0 {
		se1 = math.Sqrt((rewards[1] * (totalScore - rewards[1])) / counts[1])
	}
	seDiff = math.Sqrt(se0*se0 + se1*se1)
	if seDiff != 0 {
		zScore = (rewards[1] - rewards[0]) / seDiff
	}

	pnorm := distuv.Normal{
		Mu:    0,
		Sigma: 1,
	}

	return pnorm.CDF(-math.Abs(zScore))
}
