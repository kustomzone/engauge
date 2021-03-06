package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/EngaugeAI/engauge/db"
	"github.com/EngaugeAI/engauge/types"

	"github.com/labstack/echo/v4"
)

// EndpointList ...
func EndpointList(c echo.Context) error {
	c.Logger().Debug(c.Request().Header.Get(echo.HeaderAuthorization))
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

	endpointDocs := client.Do(&db.Op{
		Resource: db.Endpoints,
		Type:     db.List,
		Limit:    limit,
		Offset:   offset,
	})

	listView := make([]*types.EndpointListView, 0)
	for _, endpoint := range endpointDocs.Item.([]*types.Endpoint) {
		listView = append(listView, endpoint.ListView())
	}

	c.Response().Header().Add("x-total-count", strconv.Itoa(db.EndpointsCache.Len()))
	return c.JSON(http.StatusOK, listView)
}

// EndpointGet ...
func EndpointGet(c echo.Context) error {
	id := c.Param("id")

	endpointResult := client.Do(&db.Op{
		Resource: db.Endpoints,
		Type:     db.Read,
		Where: db.WhereMap{
			"item.id": types.UUIDFromString(id),
		},
	})
	if endpointResult.Error != nil {
		c.Logger().Error(endpointResult.Error)
		return c.NoContent(http.StatusInternalServerError)
	}

	endpoint := endpointResult.Item.(*types.Endpoint)

	response := &types.EndpointResponse{
		ID:         endpoint.ID,
		Action:     endpoint.Action,
		EntityType: endpoint.EntityType,
		EntityID:   endpoint.EntityID,
		OriginType: endpoint.OriginType,
		OriginID:   endpoint.OriginID,
	}

	stats, err := db.EndpointsStatsCache.AllIntervalStats(endpoint.ID.String())
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}
	response.Stats = stats

	return c.JSON(http.StatusOK, response)
}

// EndpointPost ...
func EndpointPost(c echo.Context) error {
	var endpoint types.Endpoint
	err := json.NewDecoder(c.Request().Body).Decode(&endpoint)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if existingID := db.EndpointsCache.ID(&endpoint); existingID != types.UUIDNil.UUID {
		return c.NoContent(http.StatusBadRequest)
	}

	endpoint.ID = types.NewUUID()

	createResult := client.Do(&db.Op{
		Resource: db.Endpoints,
		Type:     db.Create,
		Item:     endpoint,
	})
	if createResult.Error != nil {
		c.Logger().Error(createResult.Error)
		return c.NoContent(http.StatusInternalServerError)
	}

	db.EndpointsCache.Set(endpoint.ID, &endpoint)

	return c.JSON(http.StatusOK, endpoint)
}
