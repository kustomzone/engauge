package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/EngaugeAI/engauge/db"
	"github.com/EngaugeAI/engauge/types"

	"github.com/labstack/echo/v4"
)

// SummaryList --
func SummaryList(c echo.Context) error {
	summaryList := make([]*types.SummaryListView, 0)

	db.SummaryCache.Range(func(key, value interface{}) bool {
		summary := value.(*types.Summary)
		summaryList = append(summaryList, summary.ListView())
		return true
	})

	c.Response().Header().Add("x-total-count", strconv.Itoa(4))
	return c.JSON(http.StatusOK, summaryList)
}

// SummaryGet --
func SummaryGet(c echo.Context) error {
	interval := c.Param("id")

	// check toggle
	switch interval {
	case types.Hourly:
		if !db.GlobalSettings.StatsToggles.Hourly {
			return c.NoContent(http.StatusBadRequest)
		}
	case types.Daily:
		if !db.GlobalSettings.StatsToggles.Daily {
			return c.NoContent(http.StatusBadRequest)
		}
	case types.Weekly:
		if !db.GlobalSettings.StatsToggles.Weekly {
			return c.NoContent(http.StatusBadRequest)
		}
	case types.Monthly:
		if !db.GlobalSettings.StatsToggles.Monthly {
			return c.NoContent(http.StatusBadRequest)
		}
	case types.Quarterly:
		if !db.GlobalSettings.StatsToggles.Quarterly {
			return c.NoContent(http.StatusBadRequest)
		}
	case types.Yearly:
		if !db.GlobalSettings.StatsToggles.Yearly {
			return c.NoContent(http.StatusBadRequest)
		}
	default:
		return c.NoContent(http.StatusBadRequest)
	}

	item, ok := db.SummaryCache.Load(interval)
	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}
	summary := item.(*types.Summary)

	if time.Now().UTC().After(summary.End) {
		response := &types.SummaryResponse{
			ID: summary.SpanType,
		}
		return c.JSON(http.StatusOK, response)
	}

	return c.JSON(http.StatusOK, summary.Response())
}
