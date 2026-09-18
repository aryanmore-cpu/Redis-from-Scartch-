package main

import (
	"strconv"
	"sync"
	"time"
)

type Item struct {
	value     string
	expiresAt time.Time
}

type Store struct {
	mu sync.RWMutex

	kvStore map[string]Item
}

func NewStore() *Store {

	return &Store{

		kvStore: make(map[string]Item),
	}
}

func (s *Store) Set(key, value string, expiresAt time.Time) {

	s.mu.Lock()
	defer s.mu.Unlock()
	s.kvStore[key] = Item{value: value, expiresAt: expiresAt}

}

func (s *Store) Get(key string) (string, bool) {

	s.mu.Lock()
	defer s.mu.Unlock()
	item, exists := s.kvStore[key]

	if !exists {
		return "", false

	}

	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		delete(s.kvStore, key)
		return "", false
	}
	return item.value, true

}

func (s *Store) Delete(key string) bool {

	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.kvStore[key]
	if !exists {
		return false
	}

	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		delete(s.kvStore, key)
		return false
	}

	delete(s.kvStore, key)
	return true
}

func (s *Store) Exists(key string) bool {

	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.kvStore[key]
	if !exists {
		return false
	}

	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		delete(s.kvStore, key)
		return false
	}
	return true
}

func (s *Store) Incr(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.kvStore[key]
	if exists && !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		delete(s.kvStore, key)
		exists = false
	}

	currentVal := 0
	var existingExpiry time.Time
	if exists {
		valInt, err := strconv.Atoi(item.value)
		if err != nil {
			return 0, err
		}
		currentVal = valInt
		existingExpiry = item.expiresAt

	}

	currentVal++
	s.kvStore[key] = Item{
		value:     strconv.Itoa(currentVal),
		expiresAt: existingExpiry,
	}
	return currentVal, nil

}
