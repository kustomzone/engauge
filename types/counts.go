package types

import (
	"math"
	"sort"

	"github.com/humilityai/sam"
)

// Count holds a value name and total count value
type Count struct {
	Name  interface{} `json:"name"`
	Value int64       `json:"value"`
}

// Counts is a list of Count objects
type Counts []Count

// Rate is key-value struct for float64 values
type Rate struct {
	Name  interface{} `json:"name"`
	Value float64     `json:"value"`
}

// Rates is a list of Rate objects
type Rates []Rate

// Keys returns all names in the Counts list
func (c Counts) Keys() sam.Slice {
	if len(c) == 0 {
		return nil
	}

	switch c[0].Name.(type) {
	case string:
		var s sam.SliceString
		for _, count := range c {
			s = append(s, count.Name.(string))
		}
		return s
	case float64:
		var s sam.SliceFloat64
		for _, count := range c {
			s = append(s, count.Name.(float64))
		}
		return s
	}

	return nil
}

// Contains specifies in a name exists in Counts list
func (c Counts) Contains(key interface{}) bool {
	for _, count := range c {
		if count.Name == key {
			return true
		}
	}

	return false
}

// Get will get a count by it's index in the list
func (c Counts) Get(idx int) Count {
	if idx < 0 || idx > len(c)-1 {
		return Count{}
	}

	return c[idx]
}

// GetValue will get a count value by the count's name
func (c Counts) GetValue(name interface{}) int64 {
	for _, count := range c {
		if count.Name == name {
			return count.Value
		}
	}
	return 0
}

// Index will find the index of the count with the given name.
// It returns 0 if the count is not found in the list.
func (c Counts) Index(name interface{}) int {
	for i, count := range c {
		if count.Name == name {
			return i
		}
	}
	return 0
}

// Set will set the value of count the specified name to the specified value.
// If the count is not found in the Counts list then it will be created and appended
// to the list.
func (c Counts) Set(name interface{}, value int64) Counts {
	var exists bool
	for i, count := range c {
		if count.Name == name {
			c[i].Value = value
			exists = true
		}
	}

	if !exists {
		c = append(c, Count{
			Name:  name,
			Value: value,
		})
	}

	return c
}

// Increment will increase the count of the given
// name by 1.
// If the provided name does not exist in the list of count
// objects, then a new count object will be created with the
// provided name and initialized with a count of 1.
func (c Counts) Increment(name interface{}) Counts {
	var exists bool
	for i := range c {
		if c[i].Name == name {
			c[i].Value++
			exists = true
			break
		}
	}

	if !exists {
		c = append(c, Count{
			Name:  name,
			Value: 1,
		})
	}

	return c
}

// MeanIndex returns the index of the count which represents the categorical "mean"
// of all the unique count objects.
func (c Counts) MeanIndex() int {
	sort.Sort(c)

	var sum, mean float64
	for i, count := range c {
		total := float64(count.Value)
		sum += total
		mean += float64(i) * total
	}
	return int(math.Round(mean / sum))
}

// Sum returns the sum value of all the count values in the list.
func (c Counts) Sum() float64 {
	var total int64
	for _, count := range c {
		total += count.Value
	}
	return float64(total)
}

// Len ...
func (c Counts) Len() int { return len(c) }

// Swap ...
func (c Counts) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

// Less ....
func (c Counts) Less(i, j int) bool { return c[i].Value < c[j].Value }
