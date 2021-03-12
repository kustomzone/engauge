package api

import (
	"net/http/pprof"
	"strings"

	"github.com/labstack/echo/v4"
)

// Pprof will add all pprof endpoints to the provided echo server
func Pprof(e *echo.Echo) {
	g := e.Group("")

	pprofPaths := []string{
		"/debug/pprof/",
		"/debug/pprof/heap",
		"/debug/pprof/goroutine",
		"/debug/pprof/block",
		"/debug/pprof/threadcreate",
		"/debug/pprof/profile",
		"/debug/pprof/symbol",
		"/debug/pprof/trace",
		"/debug/pprof/mutex",
	}
	for _, path := range pprofPaths {
		g.GET(path, pprofHandler(path))
		if strings.Contains(path, "/symbol") {
			g.POST(path, pprofHandler(path))
		}
	}
}

func pprofHandler(path string) echo.HandlerFunc {
	switch path {
	case "/debug/pprof/":
		return func(ctx echo.Context) error {
			pprof.Index(ctx.Response().Writer, ctx.Request())
			return nil
		}
	case "/debug/pprof/heap":
		return func(ctx echo.Context) error {
			pprof.Handler("heap").ServeHTTP(ctx.Response(), ctx.Request())
			return nil
		}
	case "/debug/pprof/goroutine":
		return func(ctx echo.Context) error {
			pprof.Handler("goroutine").ServeHTTP(ctx.Response().Writer, ctx.Request())
			return nil
		}
	case "/debug/pprof/block":
		return func(ctx echo.Context) error {
			pprof.Handler("block").ServeHTTP(ctx.Response().Writer, ctx.Request())
			return nil
		}
	case "/debug/pprof/threadcreate":
		return func(ctx echo.Context) error {
			pprof.Handler("threadcreate").ServeHTTP(ctx.Response().Writer, ctx.Request())
			return nil
		}
	case "/debug/pprof/cmdline":
		return func(ctx echo.Context) error {
			pprof.Cmdline(ctx.Response().Writer, ctx.Request())
			return nil
		}
	case "/debug/pprof/profile":
		return func(ctx echo.Context) error {
			pprof.Profile(ctx.Response().Writer, ctx.Request())
			return nil
		}
	case "/debug/pprof/symbol":
		return func(ctx echo.Context) error {
			pprof.Symbol(ctx.Response().Writer, ctx.Request())
			return nil
		}
	case "/debug/pprof/trace":
		return func(ctx echo.Context) error {
			pprof.Trace(ctx.Response().Writer, ctx.Request())
			return nil
		}
	case "/debug/pprof/mutex":
		return func(ctx echo.Context) error {
			pprof.Handler("mutex").ServeHTTP(ctx.Response().Writer, ctx.Request())
			return nil
		}
	}

	return nil
}
