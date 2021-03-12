package ingest

import (
	"time"

	"github.com/EngaugeAI/engauge/types"

	"github.com/JKhawaja/cache"
)

var (
	// InteractionsCache holds the interaction initially and sets an expiration according to it's created timestamp
	// or its' received timestamp (if no created timestamp). This is used to attempt to preserve ordering of interactions
	// based on the actual time of the interaction on the client-side.
	InteractionsCache *cache.Cache
)

func initCache() {
	onExpires := func(item interface{}) {
		i := item.(*types.Interaction)
		buffDrop(i)
	}
	config := &cache.CacheConfig{
		OnExpires:     onExpires,
		CleanDuration: 1 * time.Second,
	}
	InteractionsCache = cache.NewCache(config)
}

func buffDrop(interaction *types.Interaction) {
	bufferChan <- interaction
}
