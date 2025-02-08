package util

import (
	"sync"
	"time"
)

// CNameMemory stores CNAMEs with timestamps and ensures thread safety.
type CNameMemory struct {
	mu     sync.Mutex
	cnames map[string]time.Time
}

// NewCNameMemory creates a new CNameMemory instance.
func NewCNameMemory() *CNameMemory {
	return &CNameMemory{
		cnames: make(map[string]time.Time),
	}
}

// Refresh updates the timestamp of a CNAME to the current time.
func (c *CNameMemory) Refresh(cname string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cnames[cname] = time.Now()
}

// GetAllActive returns all CNAMEs with timestamps not older than 5 minutes
// and removes stale entries.
func (c *CNameMemory) GetAllActive() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	fiveMinutesAgo := now.Add(-5 * time.Minute)
	activeCnames := []string{}

	for cname, timestamp := range c.cnames {
		if timestamp.After(fiveMinutesAgo) {
			activeCnames = append(activeCnames, cname)
		} else {
			delete(c.cnames, cname) // Remove stale entries
		}
	}

	return activeCnames
}

var (
	globalCNameMemory *CNameMemory
	once              sync.Once
)

// GetCNameMemory returns the singleton instance of CNameMemory.
func GetCNameMemory() *CNameMemory {
	once.Do(func() {
		globalCNameMemory = NewCNameMemory()
	})
	return globalCNameMemory
}
