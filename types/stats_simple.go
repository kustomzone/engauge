package types

import (
	"math"
	"sort"
	"time"

	"github.com/JKhawaja/errors"
	"github.com/humilityai/sam"
)

// SimpleStats --
type SimpleStats struct {
	Type     string              `json:"type"`
	Total    int64               `json:"total"`
	Mean     interface{}         `json:"mean"`
	Mode     interface{}         `json:"mode"`
	Values   []*SimpleValueStats `json:"values"`
	Variance sam.SliceFloat64    `json:"variance"`
	StdDev   sam.SliceFloat64    `json:"std_dev"`
}

// SimpleValueStats --
type SimpleValueStats struct {
	Value      interface{} `json:"value"`
	Count      int64       `json:"count"`
	Percentage float64     `json:"percentage"`
}

// NewSimpleStats --
func NewSimpleStats(value interface{}) (*SimpleStats, error) {
	s := &SimpleStats{}

	switch value.(type) {
	case string:
		v := value.(string)
		s.Type = String
		s.Total = 1
		s.Mean = v
		s.Mode = v
		s.Values = make([]*SimpleValueStats, 0)
		s.add(value)
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
		s.Values = make([]*SimpleValueStats, 0)
		s.add(value)
		s.Variance = make(sam.SliceFloat64, 1, 1)
		s.StdDev = make(sam.SliceFloat64, 1, 1)
	case time.Duration:
		v := value.(time.Duration)
		s.Type = Duration
		s.Total = 1
		s.Mean = v
		s.Mode = v
		s.Values = make([]*SimpleValueStats, 0)
		s.add(value)
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
		s.Values = make([]*SimpleValueStats, 0)
		for _, val := range v {
			s.add(val)
		}
		s.Variance = make(sam.SliceFloat64, 1, 1)
		s.StdDev = make(sam.SliceFloat64, 1, 1)
	case []string:
		v := value.([]string)
		s.Type = StringArray
		s.Total = 1
		s.Mean = v
		s.Mode = v
		s.Values = make([]*SimpleValueStats, 0)
		for _, val := range v {
			s.add(val)
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

// NewSimpleValueStats --
func NewSimpleValueStats(value interface{}) (*SimpleValueStats, error) {
	return &SimpleValueStats{
		Value: value,
		Count: 1,
	}, nil
}

// Update --
func (s *SimpleValueStats) Update() {
	s.Count++
}

// Update --
func (s *SimpleStats) UpdateValue(value interface{}) error {
	// update value stats (if exists)
	var exists bool
	for _, stats := range s.Values {
		if stats.Value == value {
			stats.Update()
			exists = true
			break
		}
	}

	// if dne: append new stats
	if !exists {
		stats, err := NewSimpleValueStats(value)
		if err != nil {
			return errors.New(err, nil)
		}
		s.Values = append(s.Values, stats)
	}

	// set percentages
	sum := s.sum()
	for _, stats := range s.Values {
		stats.Percentage = float64(stats.Count) / float64(sum)
	}

	return nil
}

func (s *SimpleStats) sum() int64 {
	var sum int64
	for _, stats := range s.Values {
		sum += stats.Count
	}
	return sum
}

func (s *SimpleStats) add(value interface{}) error {
	stats, err := NewSimpleValueStats(value)
	if err != nil {
		return errors.New(err, nil)
	}

	s.Values = append(s.Values, stats)
	return nil
}

func (s *SimpleStats) get(idx int) *SimpleValueStats {
	return s.Values[idx]
}

// returns the index of the count which represents the categorical "mean"
// of all the unique count objects.
func (s *SimpleStats) meanIndex() int {
	sort.Sort(s)
	var sum, mean float64
	for i, stats := range s.Values {
		total := float64(stats.Count)
		sum += total
		mean += float64(i) * total
	}
	return int(math.Round(mean / sum))
}

func (s *SimpleStats) max() (value interface{}) {
	max := int64(math.MinInt64)
	for _, stats := range s.Values {
		if stats.Count > max {
			max = stats.Count
			value = stats.Value
		}
	}
	return
}

func (s *SimpleStats) contains(value interface{}) bool {
	for _, stats := range s.Values {
		if stats.Value == value {
			return true
		}
	}

	return false
}

// Update will update the total and mean for the stream stats
// with the new value.
func (s *SimpleStats) Update(value interface{}) error {
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

		s.UpdateValue(v)

		// update stats
		s.Total++
		s.Mean = s.get(s.meanIndex()).Value
		s.Mode = s.max()
		sum := s.sum()
		for i := 0; i < len(s.Values); i++ {
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

		s.UpdateValue(value)
	default:
		return errors.New(ErrDataType, nil)
	}

	return nil
}

// Len ...
func (s *SimpleStats) Len() int { return len(s.Values) }

// Swap ...
func (s *SimpleStats) Swap(i, j int) {
	s.Values[i], s.Values[j] = s.Values[j], s.Values[i]
}

// Less ....
func (s *SimpleStats) Less(i, j int) bool { return s.Values[i].Count < s.Values[j].Count }
