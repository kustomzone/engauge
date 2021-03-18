package api

import (
	"net/http"
	"strconv"

	"github.com/EngaugeAI/engauge/db"
	"github.com/EngaugeAI/engauge/types"

	"github.com/labstack/echo/v4"
)

// OriginsList --
func OriginsList(c echo.Context) error {
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

	originDocs := client.Do(&db.Op{
		Resource: db.Origins,
		Type:     db.List,
		Limit:    limit,
		Offset:   offset,
	})
	if originDocs.Error != nil {
		c.Logger().Error(originDocs.Error)
		return c.NoContent(http.StatusInternalServerError)
	}

	list := originDocs.Item.([]*types.Origin)
	if len(list) == 0 {
		list = make([]*types.Origin, 0)
	}

	c.Response().Header().Add("x-total-count", strconv.Itoa(db.OriginsCache.Len()))
	return c.JSON(http.StatusOK, list)
}

// OriginGet --
func OriginGet(c echo.Context) error {
	id := c.Param("id")
	originResult := client.Do(&db.Op{
		Resource: db.Origins,
		Type:     db.Read,
		Where: db.WhereMap{
			"item.id": types.UUIDFromString(id),
		},
	})
	if originResult.Error != nil {
		c.Logger().Error(originResult.Error)
		return c.NoContent(http.StatusInternalServerError)
	}

	origin := originResult.Item.(*types.Origin)

	response := &types.OriginResponse{
		ID:         origin.ID,
		OriginType: origin.OriginType,
		OriginID:   origin.OriginID,
	}

	stats, err := db.OriginsStatsCache.AllIntervalStats(origin.ID.String())
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}
	response.Stats = stats

	return c.JSON(http.StatusOK, response)
}
