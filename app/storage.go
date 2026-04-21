package main

import (
	"sync"
	"time"
)

type Storage struct {
	Data map[string]map[string]time.Time
	Lock sync.Mutex
}
