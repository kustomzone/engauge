package api

import (
	"time"

	"github.com/EngaugeAI/engauge/db"
	"github.com/EngaugeAI/engauge/types"
)

var (
	client   db.Client
	timezone *time.Location

	/* actions */

	// Start is a query param that will set an entities status to Started
	Start = "start"
	// Stop is a query param that will set an entities status to Finished
	Stop = "stop"
)

// Init intialized the global db client
func Init(c db.Client, tz string) {
	location, err := time.LoadLocation(types.DefaultTimeZone)
	if err != nil {
		panic(err)
	}
	timezone = location

	client = c
}
