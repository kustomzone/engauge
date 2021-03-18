package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/EngaugeAI/engauge/api"
	"github.com/EngaugeAI/engauge/db"
	"github.com/EngaugeAI/engauge/db/local"
	"github.com/EngaugeAI/engauge/ingest"
	"github.com/EngaugeAI/engauge/types"

	rice "github.com/GeertJohan/go.rice"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme/autocert"
)

type EnvVars struct {
	Env          string
	Https        bool
	Basepath     string
	Timezone     string
	Sessiondelay int
	User         string
	Password     string
	Jwt          string
	Apikey       string
}

func main() {
	// init
	var env EnvVars
	envconfig.Process("engauge", &env)
	var dev bool
	if env.Env == "dev" {
		dev = true
	}

	client, err := local.NewClient(env.Basepath)
	if err != nil {
		log.Fatal(err)
	}
	ingest.Init(client)
	api.Init(client, env.Timezone)

	if env.Timezone != "" {
		types.DefaultTimeZone = env.Timezone
	}

	if env.Sessiondelay != 0 {
		types.SessionExpiryDuration = time.Duration(time.Duration(int64(env.Sessiondelay)) * time.Minute)
	}

	db.GlobalSettings.User = env.User
	db.GlobalSettings.Password = env.Password
	db.GlobalSettings.APIKey = env.Apikey
	db.GlobalSettings.JWTSecret = env.Jwt

	e := echo.New()
	e.HideBanner = true

	// middleware
	if env.Https {
		e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
		e.Pre(middleware.HTTPSRedirect())
		e.Pre(middleware.HTTPSNonWWWRedirect())

	}
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(middleware.Recover())

	// handlers
	staticHandler := http.FileServer(rice.MustFindBox("dashboard/build").HTTPBox())
	e.GET("/", echo.WrapHandler(staticHandler))
	e.GET("/static/*", echo.WrapHandler(staticHandler))
	api.AttachHandlers(e, dev)
	e.GET("/health", healthCheck)
	if dev {
		api.Pprof(e)
	}

	// start server
	go func() {
		if env.Https {
			e.Logger.Fatal(e.StartAutoTLS(":443"))
		} else {
			e.Logger.Fatal(e.Start("localhost:8080"))
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "alive")
}
