package main

import (
	"sync"
	"time"
)

type Storage struct {
	data map[string]Entry
	mu   sync.RWMutex
}

type Entry struct {
	Value string
}

func (s *Storage) Get(key string) (Entry, bool) {

	if len(key) == 0 {
		return Entry{}, false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	data, exists := s.data[key]
	return data, exists

}

func (s *Storage) Set(key string, data Entry, lifespan time.Duration) bool {

	if len(key) == 0 {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = data

	time.AfterFunc(lifespan, func() {

		s.mu.Lock()
		defer s.mu.Unlock()
		if s.data[key] == data {
			delete(s.data, key)
		}
	})

	return true

}
