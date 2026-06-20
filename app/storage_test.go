package main

import (
	"sync"
	"testing"
	"time"
)

const secondsInADay = 86400 * time.Second

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
		storage.Set(key, Entry{Value: starter}, secondsInADay)
	})
	wg.Go(func() {
		defer wg.Done()
		storage.Set(key, Entry{Value: follower}, secondsInADay)
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
func TestWipeEntry_Positive(t *testing.T) {

	storage := Storage{
		data: map[string]Entry{},
		mu:   sync.RWMutex{},
	}
	key := "random"

	entry := Entry{Value: "arbitrary"}
	added := storage.Set(key, entry, secondsInADay)
	if !added {
		t.Errorf("failed to find the expected value by key %s", key)
	}

	if _, found := storage.Get(key); found {
		t.Errorf("expected to wipe %s, yet found data", key)
	}

}

func TestWipeEntry_Negative(t *testing.T) {

	storage := Storage{
		data: map[string]Entry{},
		mu:   sync.RWMutex{},
	}
	key := "random"

	entry := Entry{Value: "arbitrary"}

	added := storage.Set(key, entry, secondsInADay)
	if !added {
		t.Errorf("failed to find the expected value by key %s", key)
	}

	if _, found := storage.Get(key); !found {
		t.Errorf("expected to keep the data %s, yet wiped it", key)
	}

}
