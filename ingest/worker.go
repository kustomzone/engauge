package ingest

import (
	"fmt"
	"sync"
	"time"

	"github.com/EngaugeAI/engauge/db"
	"github.com/EngaugeAI/engauge/types"

	"github.com/JKhawaja/errors"
)

var (
	// MaxProcWait is the maximum time to wait before processing interactions
	// in the buffer. This is useful for when incoming messages are slow and the buffer
	// is not filling up very quickly. Default time is 10 seconds.
	MaxProcWait = 10 * time.Second
	// MinBatchSize is the maximum number of interactions that will be processed in batch.
	// The larger this number is the fewer database calls will be processed.
	// The tradeoff is that the buffer size defines how many interactions could be potentially
	// lost in the event of the service going down.
	MinBatchSize = 10

	bufferChan      = make(chan *types.Interaction, MinBatchSize*2)
	buffer          = make([]*types.Interaction, 0, MinBatchSize)
	bufferMutex     = &sync.Mutex{}
	bufferUpdatedAt time.Time
)

// Init will intialize the worker that moves interactions
// from the buffer to the database, and updates all relevant analytical
// entities that are currently in-progress (started).
func Init(client db.Client) {
	initCache()
	clock(client)
	worker(client)
}

func clock(client db.Client) {
	bufferUpdatedAt = time.Now().UTC()
	go func(client db.Client) {
		for {
			tsu := time.Since(bufferUpdatedAt)
			if tsu >= MaxProcWait && len(buffer) > 0 {
				bufferMutex.Lock()
				copyBuf := make([]*types.Interaction, len(buffer))
				copy(copyBuf, buffer)
				processInteractions(client, copyBuf)
				buffer = make([]*types.Interaction, 0, MinBatchSize)
				bufferUpdatedAt = time.Now().UTC()
				bufferMutex.Unlock()
			} else {
				if tsu < MaxProcWait {
					time.Sleep(MaxProcWait - tsu)
				} else {
					time.Sleep(MaxProcWait)
				}
			}
		}
	}(client)
}

func worker(client db.Client) {
	go func(client db.Client) {
		for v := range bufferChan {
			buffer = append(buffer, v)

			// if min batch size requirement met & no interactions left in chan
			if len(buffer) >= MinBatchSize && len(bufferChan) == 0 {
				bufferMutex.Lock()
				copyBuf := make([]*types.Interaction, len(buffer))
				copy(copyBuf, buffer)
				processInteractions(client, copyBuf)
				buffer = make([]*types.Interaction, 0, MinBatchSize)
				bufferUpdatedAt = time.Now().UTC()
				bufferMutex.Unlock()
			}
		}
	}(client)
}

func processInteractions(client db.Client, interactions []*types.Interaction) {
	// process each interaction
	for _, interaction := range interactions {
		// event
		session, err := db.SessionsCache.GetSession(interaction)
		if err != nil {
			fmt.Println(errors.NewTrace(err).Error())
			continue
		}
		interaction.SessionID = &session.ID
		event := &types.Event{
			Interaction: interaction,
			Session:     session,
			Origin:      db.OriginsCache.ID(interaction.Origin()),
			Entity:      db.EntitiesCache.ID(interaction.Entity()),
			Endpoint:    db.EndpointsCache.ID(interaction.Endpoint()),
		}

		// endpoints
		err = db.EndpointsCache.Apply(event)
		if err != nil {
			fmt.Println(errors.NewTrace(err).Error())
		}

		err = db.EndpointsStatsCache.Apply(event)
		if err != nil {
			fmt.Println(errors.NewTrace(err).Error())
		}

		// origins
		db.OriginsCache.Apply(event)
		err = db.OriginsStatsCache.Apply(event)
		if err != nil {
			fmt.Println(errors.NewTrace(err).Error())
		}

		// entities
		db.EntitiesCache.Apply(event)
		err = db.EntityStatsCache.Apply(event)
		if err != nil {
			fmt.Println(errors.NewTrace(err).Error())
		}

		// properties
		err = db.PropertiesCache.Apply(event)
		if err != nil {
			fmt.Println(errors.NewTrace(err).Error())
		}

		err = db.PropertyStatsCache.Apply(event)
		if err != nil {
			fmt.Println(errors.NewTrace(err).Error())
		}

		// summaries
		for _, interval := range types.Intervals {
			var toggle bool
			switch interval {
			case types.Hourly:
				toggle = db.GlobalSettings.StatsToggles.Hourly
			case types.Daily:
				toggle = db.GlobalSettings.StatsToggles.Daily
			case types.Weekly:
				toggle = db.GlobalSettings.StatsToggles.Weekly
			case types.Monthly:
				toggle = db.GlobalSettings.StatsToggles.Monthly
			case types.Quarterly:
				toggle = db.GlobalSettings.StatsToggles.Quarterly
			case types.Yearly:
				toggle = db.GlobalSettings.StatsToggles.Yearly
			}

			if !toggle {
				continue
			}

			s, ok := db.SummaryCache.Load(interval)

			// create new summary if dne
			if !ok {
				newSummary, err := types.NewSummary(interval, event)
				if err != nil {
					fmt.Println(errors.NewTrace(err).Error())
					continue
				}
				db.SummaryCache.Store(interval, newSummary)
				continue
			}

			// update summary
			summary := s.(*types.Summary)
			if summary.Expired(interaction) {
				// new summary
				newSummary, err := types.NewSummary(interval, event)
				if err != nil {
					fmt.Println(errors.NewTrace(err).Error())
					continue
				}

				db.SummaryCache.Store(interval, newSummary)
			} else {
				err = summary.Apply(event)
				if err != nil {
					fmt.Println(errors.NewTrace(err).Error())
				}
			}
		}

		// session
		session.Update(interaction)

		// store interaction
		if db.GlobalSettings.InteractionsStorage {
			interactionResult := client.Do(&db.Op{
				Resource: db.Interactions,
				Type:     db.Create,
				Item:     interaction,
			})
			if interactionResult.Error != nil {
				fmt.Println(errors.NewTrace(err).Error())
			}
		}
	} // end process interactions loop

	/* update in db */
	updateDB(client)
}

func updateDB(client db.Client) {
	err := db.EndpointsCache.Update(func(object interface{}) error {
		endpoint, ok := object.(*types.Endpoint)
		if !ok {
			return errors.New(types.ErrAssertion, nil)
		}

		endpointUpdate := client.Do(&db.Op{
			Resource: db.Endpoints,
			Type:     db.Update,
			Where: db.WhereMap{
				"item.id": endpoint.ID,
			},
			Item:   endpoint,
			Upsert: true,
		})

		if endpointUpdate.Error != nil {
			return errors.New(endpointUpdate.Error, nil)
		}

		return nil
	})
	if err != nil {
		fmt.Println(errors.NewTrace(err).Error())
	}

	err = db.EndpointsStatsCache.Update(func(object interface{}) error {
		endpointStats, ok := object.(*types.IntervalStats)
		if !ok {
			return errors.New(types.ErrAssertion, nil)
		}

		endpointStatsUpdate := client.Do(&db.Op{
			Resource: db.EndpointStats,
			Type:     db.Update,
			Where: db.WhereMap{
				"item.id":       endpointStats.ID,
				"item.interval": endpointStats.Interval,
			},
			Item:   endpointStats,
			Upsert: true,
		})

		if endpointStatsUpdate.Error != nil {
			return errors.New(endpointStatsUpdate.Error, nil)
		}

		return nil
	})
	if err != nil {
		fmt.Println(errors.NewTrace(err).Error())
	}

	err = db.OriginsCache.Update(func(object interface{}) error {
		origin, ok := object.(*types.Origin)
		if !ok {
			return errors.New(types.ErrAssertion, nil)
		}

		originUpdate := client.Do(&db.Op{
			Resource: db.Origins,
			Type:     db.Update,
			Where: db.WhereMap{
				"item.id": origin.ID,
			},
			Item:   origin,
			Upsert: true,
		})

		if originUpdate.Error != nil {
			return errors.New(originUpdate.Error, nil)
		}

		return nil
	})
	if err != nil {
		fmt.Println(errors.NewTrace(err).Error())
	}

	err = db.OriginsStatsCache.Update(func(object interface{}) error {
		originStats, ok := object.(*types.IntervalStats)
		if !ok {
			return errors.New(types.ErrAssertion, nil)
		}

		originStatsUpdate := client.Do(&db.Op{
			Resource: db.OriginStats,
			Type:     db.Update,
			Where: db.WhereMap{
				"item.id":       originStats.ID,
				"item.interval": originStats.Interval,
			},
			Item:   originStats,
			Upsert: true,
		})

		if originStatsUpdate.Error != nil {
			return errors.New(originStatsUpdate.Error, nil)
		}

		return nil
	})
	if err != nil {
		fmt.Println(errors.NewTrace(err).Error())
	}

	err = db.EntitiesCache.Update(func(object interface{}) error {
		entity, ok := object.(*types.Entity)
		if !ok {
			return errors.New(types.ErrAssertion, nil)
		}

		entityUpdate := client.Do(&db.Op{
			Resource: db.Entities,
			Type:     db.Update,
			Where: db.WhereMap{
				"item.id": entity.ID,
			},
			Item:   entity,
			Upsert: true,
		})

		if entityUpdate.Error != nil {
			return errors.New(entityUpdate.Error, nil)
		}

		return nil
	})
	if err != nil {
		fmt.Println(errors.NewTrace(err).Error())
	}

	err = db.EntityStatsCache.Update(func(object interface{}) error {
		entityStats, ok := object.(*types.IntervalStats)
		if !ok {
			return errors.New(types.ErrAssertion, nil)
		}

		entityStatsUpdate := client.Do(&db.Op{
			Resource: db.EntityStats,
			Type:     db.Update,
			Where: db.WhereMap{
				"item.id":       entityStats.ID,
				"item.interval": entityStats.Interval,
			},
			Item:   entityStats,
			Upsert: true,
		})

		if entityStatsUpdate.Error != nil {
			return errors.New(entityStatsUpdate.Error, nil)
		}

		return nil
	})
	if err != nil {
		fmt.Println(errors.NewTrace(err).Error())
	}

	err = db.PropertiesCache.Update(func(object interface{}) error {
		property, ok := object.(*types.Property)
		if !ok {
			return errors.New(types.ErrAssertion, nil)
		}
		propertyUpdate := client.Do(&db.Op{
			Resource: db.Properties,
			Type:     db.Update,
			Where: db.WhereMap{
				"item.name": property.Name,
			},
			Item:   property,
			Upsert: true,
		})

		if propertyUpdate.Error != nil {
			return errors.New(propertyUpdate.Error, nil)
		}

		return nil
	})
	if err != nil {
		fmt.Println(errors.NewTrace(err).Error())
	}

	err = db.PropertyStatsCache.Update(func(object interface{}) error {
		propertyStats, ok := object.(*types.PropertyStats)
		if !ok {
			return errors.New(types.ErrAssertion, nil)
		}
		propertyStatsUpdate := client.Do(&db.Op{
			Resource: db.PropertyStats,
			Type:     db.Update,
			Where: db.WhereMap{
				"item.name":     propertyStats.Name,
				"item.spanType": propertyStats.SpanType,
			},
			Item:   propertyStats,
			Upsert: true,
		})

		if propertyStatsUpdate.Error != nil {
			return errors.New(propertyStatsUpdate.Error, nil)
		}

		return nil
	})
	if err != nil {
		fmt.Println(errors.NewTrace(err).Error())
	}

	db.SummaryCache.Range(func(key, value interface{}) bool {
		interval := key.(string)
		summary := value.(*types.Summary)

		summaryUpdate := client.Do(&db.Op{
			Resource: db.Summaries,
			Type:     db.Update,
			Where: db.WhereMap{
				"item.interval": interval,
			},
			Item:   summary,
			Upsert: true,
		})

		if summaryUpdate.Error != nil {
			fmt.Println(errors.New(summaryUpdate.Error, nil))
		}

		return true
	})
}
