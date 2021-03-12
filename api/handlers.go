package api

import (
	"github.com/EngaugeAI/engauge/db"
	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// AttachHandlers will attach handlers to the provided server object
func AttachHandlers(server *echo.Echo, dev bool) {
	server.POST("/login", Login)
	server.GET("/refresh-token", RefreshToken)
	server.GET("/logout", Logout)

	api := server.Group("/api")
	if !dev {
		api.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
			KeyLookup: "header:api-key",
			Validator: func(key string, c echo.Context) (bool, error) {
				return key == db.GlobalSettings.APIKey, nil
			},
		}))
	}

	dashboard := server.Group("/dashboard")
	dashboard.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    []byte(db.GlobalSettings.JWTSecret),
		SigningMethod: middleware.AlgorithmHS256,
		ContextKey:    "user",
		TokenLookup:   "cookie:token",
		Claims:        jwt.MapClaims{},
	}))

	api.POST("/interaction", interactionPost, middleware.BodyLimit("2K"))

	// summary
	dashboard.GET("/summaries", SummaryList)
	dashboard.GET("/summaries/:id", SummaryGet)

	// properties
	dashboard.GET("/properties/:id", PropertiesGet)
	dashboard.GET("/properties", PropertiesList)

	// endpoints
	dashboard.GET("/endpoint", EndpointList)
	dashboard.POST("/endpoint", EndpointPost)
	dashboard.GET("/endpoint/:id", EndpointGet)

	// origins
	dashboard.GET("/origin", OriginsList)
	dashboard.GET("/origin/:id", OriginGet)

	// entities
	dashboard.GET("/entity", EntityList)
	dashboard.GET("/entity/:id", EntityGet)

	// settings
	dashboard.GET("/settings", SettingsList)
	dashboard.GET("/settings/:id", SettingsGet)
	dashboard.PUT("/settings", SettingsPut)
}
