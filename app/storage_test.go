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
func TestWipeEntry_Positive(t *testing.T) {

	storage := Storage{
		data: map[string]Entry{},
		mu:   sync.RWMutex{},
	}
	key := "random"

	entry := Entry{Value: "arbitrary", TimeStamp: time.Now()}
	added := storage.Set(key, entry)
	if !added {
		t.Errorf("failed to find the expected value by key %s", key)
	}

	wiped := storage.WipeEntry(key, entry)

	if !wiped {
		t.Errorf("failed to wipe the data %s", key)
	}

	if _, found := storage.Get(key); found {
		t.Errorf("expected to wipe %s, yet found data", key)
	}

	added = storage.Set(key, entry)
	if !added {
		t.Errorf("failed to find the expected value by key %s", key)
	}

	missing := Entry{Value: "missing", TimeStamp: time.Now()}
	wiped = storage.WipeEntry(key, missing)

	if wiped {
		t.Errorf("wipe the data %s by mistake", key)
	}

	if _, found := storage.Get(key); !found {
		t.Errorf("expected to keep the data %s, yet wiped it", key)
	}

}

func TestWipeEntry_Negative(t *testing.T) {

	key := "random"

	entry := Entry{Value: "arbitrary", TimeStamp: time.Now()}

	added := storage.Set(key, entry)
	if !added {
		t.Errorf("failed to find the expected value by key %s", key)
	}

	missing := Entry{Value: "missing", TimeStamp: time.Now()}
	wiped := storage.WipeEntry(key, missing)

	if wiped {
		t.Errorf("wipe the data %s by mistake", key)
	}

	if _, found := storage.Get(key); !found {
		t.Errorf("expected to keep the data %s, yet wiped it", key)
	}

}
