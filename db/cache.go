package db

import (
	"sync"

	"github.com/EngaugeAI/engauge/types"
)

var (
	// GlobalSettings --
	GlobalSettings *types.Settings
	// SummaryCache --
	SummaryCache *sync.Map
	// SessionsCache holds the current active sessions
	SessionsCache *types.UserSessions
	// EndpointsCache --
	EndpointsCache *types.Endpoints
	// EndpointsStatsCache --
	EndpointsStatsCache *types.IntervalStatsList
	// OriginsCache --
	OriginsCache *types.Origins
	// OriginsStatsCache --
	OriginsStatsCache *types.IntervalStatsList
	// EntitiesCache --
	EntitiesCache *types.Entities
	// EntityStatsCache --
	EntityStatsCache *types.IntervalStatsList
	// PropertiesCache --
	PropertiesCache *types.Properties
	// PropertyStatsCache --
	PropertyStatsCache *types.PropertyStatsList
	// TimestampFormats holds the list of timestamp formats that have been seen across the system
	TimestampFormats []string
)
