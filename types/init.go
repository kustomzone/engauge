package types

import (
	"encoding/gob"
	"time"

	"github.com/humilityai/sam"
	"github.com/humilityai/temporal"
)

func init() {
	// extras
	gob.Register(NewUUID())
	gob.Register(sam.SliceFloat64{})
	gob.Register(make(temporal.YearSeasons))
	gob.Register(time.Duration(0))

	// stats
	gob.Register(new(SimpleValueStats))
	gob.Register(make([]*SimpleValueStats, 0))
	gob.Register(new(SimpleStats))
	gob.Register(new(ValueStats))
	gob.Register(make([]*ValueStats, 0))
	gob.Register(new(Stats))
	gob.Register(new(NamedStats))
	gob.Register(make([]*NamedStats, 0))
	gob.Register(new(NamedStatsList))
	gob.Register(new(SessionStatsList))
	gob.Register(new(ConversionStatsList))
	gob.Register(new(ConversionStats))
	gob.Register(new(UnitMetrics))

	// summmaries
	gob.Register(new(Summary))

	// origins
	gob.Register(new(Origin))
	gob.Register(new(Origins))
	gob.Register(new(OriginProfile))

	// endpoints
	gob.Register(new(Endpoints))
	gob.Register(new(EndpointProfile))

	// entities
	gob.Register(new(Entity))
	gob.Register(new(EntityProfile))

	// interval stats
	gob.Register(new(Updater))
	gob.Register(new(IntervalStats))
	gob.Register(new(IntervalStatsList))
}
