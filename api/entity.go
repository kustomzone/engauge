package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/EngaugeAI/engauge/db"
	"github.com/EngaugeAI/engauge/types"

	"github.com/labstack/echo/v4"
)

// EntityList --
func EntityList(c echo.Context) error {
	var limit, offset *int64
	l := c.QueryParam("limit")
	if l != "" {
		i, err := strconv.Atoi(l)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		li := int64(i)
		limit = &li
	}
	o := c.QueryParam("offset")
	if o != "" {
		i, err := strconv.Atoi(o)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		oi := int64(i)
		offset = &oi
	}

	entityDocs := client.Do(&db.Op{
		Resource: db.Entities,
		Type:     db.List,
		Limit:    limit,
		Offset:   offset,
	})
	if entityDocs.Error != nil {
		c.Logger().Error(entityDocs.Error)
		return c.NoContent(http.StatusInternalServerError)
	}

	list := make([]*types.Entity, 0)
	for _, e := range entityDocs.Item.([]*types.Entity) {
		list = append(list, e)
	}

	c.Response().Header().Add("x-total-count", strconv.Itoa(db.EntitiesCache.Len()))
	return c.JSON(http.StatusOK, list)
}

// EntityGet --
func EntityGet(c echo.Context) error {
	id := c.Param("id")
	entityResult := client.Do(&db.Op{
		Resource: db.Entities,
		Type:     db.Read,
		Where: db.WhereMap{
			"item.id": types.UUIDFromString(id),
		},
	})
	if entityResult.Error != nil {
		c.Logger().Error(entityResult.Error)
		return c.NoContent(http.StatusInternalServerError)
	}

	entity := entityResult.Item.(*types.Entity)

	response := &types.EntityResponse{
		ID:         entity.ID,
		EntityType: &entity.EntityType,
		EntityID:   &entity.EntityID,
	}
	for _, interval := range types.Spans {
		switch interval {
		case types.Hourly:
			if db.GlobalSettings.StatsToggles.Hourly {
				stats, err := db.EntityStatsCache.Get(id, interval)
				if err != nil {
					c.Logger().Error(err)
					return c.NoContent(http.StatusInternalServerError)
				}

				if time.Now().UTC().After(stats.End) {
					stats = &types.EntityStats{}
				}

				response.HourlyStats = stats
			}
		case types.Daily:
			if db.GlobalSettings.StatsToggles.Daily {
				stats, err := db.EntityStatsCache.Get(id, interval)
				if err != nil {
					c.Logger().Error(err)
					return c.NoContent(http.StatusInternalServerError)
				}

				if time.Now().UTC().After(stats.End) {
					stats = &types.EntityStats{}
				}

				response.DailyStats = stats
			}
		case types.Weekly:
			if db.GlobalSettings.StatsToggles.Weekly {
				stats, err := db.EntityStatsCache.Get(id, interval)
				if err != nil {
					c.Logger().Error(err)
					return c.NoContent(http.StatusInternalServerError)
				}

				if time.Now().UTC().After(stats.End) {
					stats = &types.EntityStats{}
				}

				response.WeeklyStats = stats
			}
		case types.Monthly:
			if db.GlobalSettings.StatsToggles.Monthly {
				stats, err := db.EntityStatsCache.Get(id, interval)
				if err != nil {
					c.Logger().Error(err)
					return c.NoContent(http.StatusInternalServerError)
				}

				if time.Now().UTC().After(stats.End) {
					stats = &types.EntityStats{}
				}

				response.MonthlyStats = stats
			}
		case types.AllTime:
			stats, err := db.EntityStatsCache.Get(id, interval)
			if err != nil {
				c.Logger().Error(err)
				return c.NoContent(http.StatusInternalServerError)
			}
			response.AlltimeStats = stats
		}
	}

	return c.JSON(http.StatusOK, response)
}
