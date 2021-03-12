package local

import (
	"log"
	"sync"
	"time"

	"github.com/EngaugeAI/engauge/db"
	"github.com/EngaugeAI/engauge/types"

	"github.com/JKhawaja/cache"
	"github.com/JKhawaja/errors"
)

func (c *Client) InitCache() {
	log.Println("creating sessions cache")
	c.initSessionsCache()

	log.Println("setting global settings")
	db.GlobalSettings = types.NewSettings()
	settingsResult := c.Do(&db.Op{
		Resource: db.Settings,
		Type:     db.Read,
	})
	if settingsResult.Error == types.ErrDNE {
		createResults := c.Do(&db.Op{
			Resource: db.Settings,
			Type:     db.Create,
			Item:     db.GlobalSettings,
		})
		if createResults.Error != nil {
			panic(createResults.Error)
		}
	} else if settingsResult.Error != nil {
		panic(settingsResult.Error)
	} else {
		db.GlobalSettings = settingsResult.Item.(*types.Settings)
	}

	log.Println("loading summaries")
	db.SummaryCache = &sync.Map{}
	summariesResult := c.Do(&db.Op{
		Resource: db.Summaries,
		Type:     db.List,
	})
	if summariesResult.Error != nil {
		panic(summariesResult.Error)
	}
	for _, summary := range summariesResult.Item.([]*types.Summary) {
		db.SummaryCache.Store(summary.SpanType, summary)
	}

	log.Println("loading properties")
	db.PropertiesCache = types.NewProperties()
	propertiesResult := c.Do(&db.Op{
		Resource: db.Properties,
		Type:     db.List,
	})
	if propertiesResult.Error != nil {
		panic(propertiesResult.Error)
	}
	for _, property := range propertiesResult.Item.([]*types.Property) {
		db.PropertiesCache.Set(property.Name, property)
	}

	log.Println("loading property stats")
	db.PropertyStatsCache = types.NewPropertyStatsList()
	propertyStatsResult := c.Do(&db.Op{
		Resource: db.PropertyStats,
		Type:     db.List,
	})
	if propertyStatsResult.Error != nil {
		panic(propertyStatsResult.Error)
	}
	err := db.PropertyStatsCache.Load(propertyStatsResult.Item.([]*types.PropertyStats))
	if err != nil {
		panic(err)
	}

	log.Println("loading origins")
	db.OriginsCache = types.NewOrigins()
	originsResult := c.Do(&db.Op{
		Resource: db.Origins,
		Type:     db.List,
	})
	if originsResult.Error != nil {
		panic(originsResult.Error)
	}
	for _, origin := range originsResult.Item.([]*types.Origin) {
		err := db.OriginsCache.Set(origin.ID.UUID, origin)
		if err != nil {
			panic(err)
		}
	}

	log.Println("loading origin stats")
	db.OriginsStatsCache = types.NewOriginStatsList()
	originStatsResult := c.Do(&db.Op{
		Resource: db.OriginStats,
		Type:     db.List,
	})
	if originStatsResult.Error != nil {
		panic(originStatsResult.Error)
	}
	err = db.OriginsStatsCache.Load(originStatsResult.Item.([]*types.OriginStats))
	if err != nil {
		panic(err)
	}

	log.Println("loading entities")
	db.EntitiesCache = types.NewEntities()
	entitiesResult := c.Do(&db.Op{
		Resource: db.Entities,
		Type:     db.List,
	})
	if entitiesResult.Error != nil {
		panic(entitiesResult.Error)
	}
	for _, entity := range entitiesResult.Item.([]*types.Entity) {
		db.EntitiesCache.Set(entity)
	}

	log.Println("loading entity stats")
	db.EntityStatsCache = types.NewEntityStatsList()
	entityStatsResult := c.Do(&db.Op{
		Resource: db.EntityStats,
		Type:     db.List,
	})
	if entityStatsResult.Error != nil {
		panic(entityStatsResult.Error)
	}
	err = db.EntityStatsCache.Load(entityStatsResult.Item.([]*types.EntityStats))
	if err != nil {
		panic(err)
	}

	log.Println("loading endpoints")
	db.EndpointsCache = types.NewEndpoints()
	endpointsResult := c.Do(&db.Op{
		Resource: db.Endpoints,
		Type:     db.List,
	})
	if endpointsResult.Error != nil {
		panic(endpointsResult.Error)
	}

	for _, endpoint := range endpointsResult.Item.([]*types.Endpoint) {
		err := db.EndpointsCache.Set(endpoint.ID.UUID, endpoint)
		if err != nil {
			panic(err)
		}
	}

	log.Println("loading endpoint stats")
	db.EndpointsStatsCache = types.NewEndpointStatsList()
	endpointStatsResult := c.Do(&db.Op{
		Resource: db.EndpointStats,
		Type:     db.List,
	})
	if endpointStatsResult.Error != nil {
		panic(endpointStatsResult.Error)
	}
	err = db.EndpointsStatsCache.Load(endpointStatsResult.Item.([]*types.EndpointStats))
	if err != nil {
		panic(err)
	}
}

func (c *Client) initSessionsCache() {
	config := &cache.CacheConfig{
		OnExpires: func(item interface{}) {
			sess := item.(*types.UserSession)

			// update summaries
			db.SummaryCache.Range(func(key, value interface{}) bool {
				summary := value.(*types.Summary)
				spanType := summary.SpanType

				var toggle bool
				switch spanType {
				case types.Hourly:
					toggle = db.GlobalSettings.StatsToggles.Hourly
				case types.Daily:
					toggle = db.GlobalSettings.StatsToggles.Daily
				case types.Weekly:
					toggle = db.GlobalSettings.StatsToggles.Weekly
				case types.Monthly:
					toggle = db.GlobalSettings.StatsToggles.Monthly
				}

				if !toggle {
					return true
				}

				err := summary.SessionExpirationUpdate(sess)
				if err != nil {
					log.Println(errors.NewTrace(err).Error())
					return true
				}

				summaryResult := c.Do(&db.Op{
					Resource: db.Summaries,
					Type:     db.Update,
					Where: db.WhereMap{
						"spanType": spanType,
					},
					Item: summary,
				})
				if summaryResult.Error != nil {
					log.Println(errors.NewTrace(summaryResult.Error).Error())
				}

				return true
			})
		},
		Refresh:         true,
		RefreshDuration: types.SessionExpiryDuration,
		CleanDuration:   1 * time.Minute,
	}

	db.SessionsCache = types.NewUserSessions(cache.NewCache(config))
}
