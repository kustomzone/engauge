package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/EngaugeAI/engauge/db"
	"github.com/EngaugeAI/engauge/types"

	"github.com/labstack/echo/v4"
)

// SettingsList --
func SettingsList(c echo.Context) error {
	db.GlobalSettings.ID = strconv.Itoa(1)
	settingsList := []*types.Settings{db.GlobalSettings}
	c.Response().Header().Add("x-total-count", strconv.Itoa(1))
	return c.JSON(http.StatusOK, settingsList)
}

// SettingsGet --
func SettingsGet(c echo.Context) error {
	db.GlobalSettings.ID = strconv.Itoa(1)
	return c.JSON(http.StatusOK, db.GlobalSettings)
}

// SettingsPut --
func SettingsPut(c echo.Context) error {
	var request types.Settings
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return echo.ErrBadRequest
	}

	db.GlobalSettings.StatsToggles = request.StatsToggles
	db.GlobalSettings.InteractionsStorage = request.InteractionsStorage

	// update in db
	campaignUpdate := client.Do(&db.Op{
		Resource: db.Settings,
		Type:     db.Update,
		Item:     db.GlobalSettings,
	})
	if campaignUpdate.Error != nil {
		c.Logger().Error(campaignUpdate.Error)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, db.GlobalSettings)
}
