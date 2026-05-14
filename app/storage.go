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
	Value     string
	TimeStamp time.Time
}

func (s *Storage) Get(key string) (bool, Entry) {

	if len(key) == 0 {
		return false, Entry{}
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	data, exists := s.data[key]
	return exists, data

}

func (s *Storage) Set(key string, data Entry) bool {

	if len(key) == 0 {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = data
	return true

}
