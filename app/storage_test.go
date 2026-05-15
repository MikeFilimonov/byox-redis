package main

import (
	"sync"
	"testing"
	"time"
)

func TestStorageRace(t *testing.T) {

	storage := Storage{
		data: map[string]Entry{},
		mu:   sync.RWMutex{},
	}
	key := "random"
	starter := "generic"
	follower := "default"

	var wg sync.WaitGroup
	wg.Add(2)

	wg.Go(func() {
		defer wg.Done()
		storage.Set(key, Entry{Value: starter, TimeStamp: time.Now()})
	})
	wg.Go(func() {
		defer wg.Done()
		storage.Set(key, Entry{Value: follower, TimeStamp: time.Now()})
	})
	wg.Wait()

	result, found := storage.Get(key)
	if !found {
		t.Errorf("failed to find the expected value by key %s", key)
	}
	if !(result.Value == starter || result.Value == follower) {
		t.Errorf("expected %s(%s), got %s", starter, follower, result.Value)
	}

}
