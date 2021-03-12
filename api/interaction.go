package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/EngaugeAI/engauge/db"
	"github.com/EngaugeAI/engauge/ingest"
	"github.com/EngaugeAI/engauge/types"

	"github.com/labstack/echo/v4"
)

var (
	formats = []string{
		"2006",
		"2006-1",
		"2006-1-2",
		"2006-1-2 15",
		"2006-1-2 15:4",
		"2006-1-2 15:4:5",
		"1-2",
		"15:4:5",
		"15:4",
		"15",
		"2006-01-02'T'15:04:05",
		"2006-01-02T15:04:05-0700",
		"2006 Jan 02 15:04:05",
		"2006-01-02T15:04:05-07:00",
		"15:4:5 Jan 2, 2006 MST",
		"2006-01-02 15:04:05.999999999 -0700 MST",
		"Jan 02 15:04:05 -0700 2006",
		"02/Jan/2006:15:04:05 -0700",
		"Jan 02 2006 15:04:05",
		"Jan 02 15:04:05 2006",
		"2006.1.2",
		"2006.1.2 15:04:05",
		"2006.01.02",
		"2006.01.02 15:04:05",
		"1/2/2006",
		"1/2/2006 15:4:5",
		"2006/01/02",
		"2006/01/02 15:04:05",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}
)

func interactionPost(c echo.Context) error {
	var interaction *types.Interaction
	err := json.NewDecoder(c.Request().Body).Decode(&interaction)
	if err != nil {
		return echo.ErrBadRequest
	}

	// stamp
	t := time.Now().In(timezone)
	interaction.ReceivedAt = &t

	// validate
	err = interaction.Validate()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	go func(i *types.Interaction) {
		if i.Timestamp != nil {
			timestamp, err := parseTimestamp(*i.Timestamp)
			if err != nil {
				log.Println(err)
			}
			i.CreatedAt = &timestamp
		} else {
			i.CreatedAt = i.ReceivedAt
		}

		expiresAt := i.CreatedAt.Add(3 * time.Second)
		tte := expiresAt.UTC().Sub(time.Now().UTC())
		ingest.InteractionsCache.Add(i.String(), i, tte)
	}(interaction)

	return c.NoContent(http.StatusOK)
}

func parseTimestamp(ts string) (timestamp time.Time, err error) {
	for _, frmt := range db.TimestampFormats {
		timestamp, err = time.ParseInLocation(frmt, ts, timezone)
		if err == nil {
			return
		}
	}

	for _, frmt := range formats {
		timestamp, err = time.ParseInLocation(frmt, ts, timezone)
		if err == nil {
			db.TimestampFormats = append(db.TimestampFormats, frmt)
			return
		}
	}

	err = types.ErrTimestamp
	return
}
