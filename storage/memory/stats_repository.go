package memory

import (
	"gav/internal/stats"
	"sync"
)

type StatsReposytory struct {
	mu 		sync.RWMutex
	stats	map[uint]*stats.UserStats	
}