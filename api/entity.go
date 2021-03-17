package api

import (
	"net/http"
	"strconv"

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
		Stats:      &types.AllIntervalStats{},
	}

	stats, err := db.EntityStatsCache.AllIntervalStats(entity.ID.String())
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}
	response.Stats = stats

	return c.JSON(http.StatusOK, response)
}
