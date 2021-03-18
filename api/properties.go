package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/EngaugeAI/engauge/db"
	"github.com/EngaugeAI/engauge/types"

	"github.com/labstack/echo/v4"
)

// PropertiesList --
func PropertiesList(c echo.Context) error {
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

	propertiesResult := client.Do(&db.Op{
		Resource: db.Properties,
		Type:     db.List,
		Limit:    limit,
		Offset:   offset,
	})
	if propertiesResult.Error != nil {
		c.Logger().Error(propertiesResult.Error)
		return c.NoContent(http.StatusInternalServerError)
	}

	items := propertiesResult.Item.([]*types.Property)

	listView := make(types.PropertyListViews, 0)
	for _, property := range items {
		listView = append(listView, property.ListView())
	}

	c.Response().Header().Add("x-total-count", strconv.Itoa(len(listView)))
	return c.JSON(http.StatusOK, listView)
}

// PropertiesGet --
func PropertiesGet(c echo.Context) error {
	id := c.Param("id")

	p := db.PropertiesCache.Get(id)
	if p == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	prop := p.(*types.Property)

	response := prop.Response()
	for _, interval := range types.Intervals {
		switch interval {
		case types.Hourly:
			if db.GlobalSettings.StatsToggles.Hourly {
				stats, err := db.PropertyStatsCache.Get(id, interval)
				if err != nil {
					c.Logger().Error(err)
					return c.NoContent(http.StatusInternalServerError)
				}

				if time.Now().UTC().After(stats.End) {
					stats = &types.PropertyStats{}
				}

				response.HourlyStats = stats
			}
		case types.Daily:
			if db.GlobalSettings.StatsToggles.Daily {
				stats, err := db.PropertyStatsCache.Get(id, interval)
				if err != nil {
					c.Logger().Error(err)
					return c.NoContent(http.StatusInternalServerError)
				}

				if time.Now().UTC().After(stats.End) {
					stats = &types.PropertyStats{}
				}

				response.DailyStats = stats
			}
		case types.Weekly:
			if db.GlobalSettings.StatsToggles.Weekly {
				stats, err := db.PropertyStatsCache.Get(id, interval)
				if err != nil {
					c.Logger().Error(err)
					return c.NoContent(http.StatusInternalServerError)
				}

				if time.Now().UTC().After(stats.End) {
					stats = &types.PropertyStats{}
				}

				response.WeeklyStats = stats
			}
		case types.Monthly:
			if db.GlobalSettings.StatsToggles.Monthly {
				stats, err := db.PropertyStatsCache.Get(id, interval)
				if err != nil {
					c.Logger().Error(err)
					return c.NoContent(http.StatusInternalServerError)
				}

				if time.Now().UTC().After(stats.End) {
					stats = &types.PropertyStats{}
				}

				response.MonthlyStats = stats
			}
		case types.Quarterly:
			if db.GlobalSettings.StatsToggles.Quarterly {
				stats, err := db.PropertyStatsCache.Get(id, interval)
				if err != nil {
					c.Logger().Error(err)
					return c.NoContent(http.StatusInternalServerError)
				}

				if time.Now().UTC().After(stats.End) {
					stats = &types.PropertyStats{}
				}

				response.QuarterlyStats = stats
			}
		case types.Yearly:
			if db.GlobalSettings.StatsToggles.Yearly {
				stats, err := db.PropertyStatsCache.Get(id, interval)
				if err != nil {
					c.Logger().Error(err)
					return c.NoContent(http.StatusInternalServerError)
				}

				if time.Now().UTC().After(stats.End) {
					stats = &types.PropertyStats{}
				}

				response.YearlyStats = stats
			}
		}
	}

	return c.JSON(http.StatusOK, response)
}
