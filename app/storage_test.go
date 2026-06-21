package main

import (
	"sync"
	"testing"
	"time"
)

const ttl = 2 * time.Second

func TestStorageRace(t *testing.T) {

	storage := Storage{
		data: map[string]*Entry{},
		mu:   sync.RWMutex{},
	}
	key := "random"
	starter := "generic"
	follower := "default"

	var wg sync.WaitGroup
	wg.Add(2)

	wg.Go(func() {
		defer wg.Done()
		storage.Set(key, &Entry{Value: starter}, ttl)
	})
	wg.Go(func() {
		defer wg.Done()
		storage.Set(key, &Entry{Value: follower}, ttl)
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
		data: map[string]*Entry{},
		mu:   sync.RWMutex{},
	}
	key := "random"

	entry := &Entry{Value: "arbitrary"}
	added := storage.Set(key, entry, ttl)
	time.Sleep(7 * time.Second)
	if !added {
		t.Errorf("failed to find the expected value by key %s", key)
	}
	added = storage.Set(key, entry, ttl)

	time.Sleep(11 * time.Second)

	if _, found := storage.Get(key); found {
		t.Errorf("expected to wipe %s, yet found data", key)
	}

}

func TestWipeEntry_Negative(t *testing.T) {

	storage := Storage{
		data: map[string]*Entry{},
		mu:   sync.RWMutex{},
	}
	key := "random"

	entry := &Entry{Value: "arbitrary"}

	storage.Set(key, entry, ttl)
	time.Sleep(1 * time.Second)
	if _, found := storage.Get(key); !found {
		t.Errorf("expected to keep the data %s, yet wiped it", key)
	}

}
