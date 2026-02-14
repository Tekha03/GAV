package memory

import "sync"

type ReactionRepository struct {
	mu sync.RWMutex
	
}