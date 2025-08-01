package util


import (
	"os"
	"fmt"
	"sync"
	"time"
)

var (
	genesisOnce sync.Once
	genesisTime time.Time
)

func getGenesis() time.Time {
	genesisOnce.Do(func() {
		const defaultGenesisRFC3339 = "2020-12-01T12:00:23Z"

		raw := os.Getenv("GENESIS_TIME")
		if raw == "" {
			raw = defaultGenesisRFC3339
		}

		// Try RFC3339 first.
		if t, err := time.Parse(time.RFC3339, raw); err == nil {
			genesisTime = t.UTC()
			fmt.Printf("beacon: parse GENESIS_TIME=%q as %q\n", raw, genesisTime)
			return
		}

		// If both parsers fail, use the baked-in default.
		fmt.Printf("beacon: cannot parse GENESIS_TIME=%q, using default %s\n", raw, defaultGenesisRFC3339)
		genesisTime, _ = time.Parse(time.RFC3339, defaultGenesisRFC3339)
	})
	return genesisTime
}

// ----------------------------------------------------------------------------
// Slot timing helper that uses getGenesis()
// ----------------------------------------------------------------------------

func slotTime(slot uint64) time.Time {
	const slotDuration = 12 * time.Second
	return getGenesis().Add(time.Duration(slot) * slotDuration)
}
