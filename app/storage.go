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

func (s *Storage) WipeEntry(key string, previous Entry) bool {

	if len(key) == 0 {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	currentEntry, found := s.data[key]
	if !found {
		return found
	}

	gottaUpdate := previous.Value == currentEntry.Value &&
		previous.TimeStamp.Equal(currentEntry.TimeStamp)
	if gottaUpdate {
		delete(s.data, key)
	}
	return gottaUpdate

}
