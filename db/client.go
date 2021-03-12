package db

const (
	/*  resource types */

	// Interactions is a resource type
	Interactions = "interactions"
	// Endpoints is a resource type
	Endpoints = "endpoints"
	// EndpointStats is a resource type
	EndpointStats = "endpointStats"
	// Origins is a resource type
	Origins = "origins"
	// OriginStats is a resource type
	OriginStats = "originStats"
	// Entities is a resource type
	Entities = "entities"
	// EntityStats is a resource type
	EntityStats = "entityStats"
	// Properties is a resource type
	Properties = "properties"
	// PropertyStats is a resource type
	PropertyStats = "propertyStats"
	// Summaries is a resource type (current summaries)
	Summaries = "summaries"
	// Settings is a resource type
	Settings = "settings"

	/*  operation types */

	// Create is an operation type
	Create = "create"
	// Update is an operation type
	Update = "update"
	// Read is an operation type
	Read = "read"
	// Delete is an operation type
	Delete = "delete"
	// List is an operation type
	List = "list"
	// Count will return the count of documents in the resource
	Count = "count"
)

// Client --
type Client interface {
	Do(operation *Op) Result
}

// Where --
type Where interface {
	Ordered() bool
	Keys() []string
	Values() []interface{}
	Pairs() WhereList
}

// Pair --
type Pair struct {
	Key   string
	Value interface{}
}

// WhereList --
type WhereList []Pair

// WhereMap --
type WhereMap map[string]interface{}

// Ordered --
func (w WhereMap) Ordered() bool {
	return false
}

// Keys --
func (w WhereMap) Keys() []string {
	keys := make([]string, 0, len(w))
	for key := range w {
		keys = append(keys, key)
	}
	return keys
}

// Values --
func (w WhereMap) Values() []interface{} {
	values := make([]interface{}, 0, len(w))
	for _, value := range w {
		values = append(values, value)
	}
	return values
}

// Pairs --
func (w WhereMap) Pairs() WhereList {
	pairs := make(WhereList, 0, len(w))

	for key, value := range w {
		pairs = append(pairs, Pair{
			Key:   key,
			Value: value,
		})
	}

	return pairs
}

// Ordered --
func (w WhereList) Ordered() bool {
	return true
}

// Keys --
func (w WhereList) Keys() []string {
	keys := make([]string, 0, len(w))
	for _, v := range w {
		keys = append(keys, v.Key)
	}
	return keys
}

// Values --
func (w WhereList) Values() []interface{} {
	values := make([]interface{}, 0, len(w))
	for _, v := range w {
		values = append(values, v.Value)
	}
	return values
}

// Pairs --
func (w WhereList) Pairs() WhereList {
	return w
}

// Op represents a generic database operation.
type Op struct {
	Resource      string
	Type          string      // operation-type
	Item          interface{} // e.g., document
	Where         Where       // optional, but should most likely contain an id
	Upsert        bool
	Limit, Offset *int64 // optional
}

// Result represents a generic database result after executing an Op
type Result struct {
	Item  interface{} // e.g., document
	Error error
}
