package ratelimiter

import (
	"sync"
	"time"
)

type Store interface {
	increment(key string, window time.Duration) (prev_current int, curr_current int, window_start time.Time, err error)
}

type Window struct {
	prev_current int
	curr_current int
	windowStart  time.Time
}

type MemoryStore struct {
	ms     sync.RWMutex
	data   map[string]*Window
	stopCH chan struct{}
}

func NewMemory(clear_up_interval time.Duration) *MemoryStore {

	ms := &MemoryStore{
		ms:     sync.RWMutex{},
		data:   make(map[string]*Window),
		stopCH: make(chan struct{}),
	}

	go ms.cleanupMemoryRamUssage(clear_up_interval)
	return ms

}

func (ms *MemoryStore) increment(key string, window time.Duration) (prev_current int, curr_current int, windowStart time.Time, err error) {

	ms.ms.Lock()
	defer ms.ms.Unlock()

	now := time.Now()
	entry, exists := ms.data[key]

	if !exists {
		ws := now.Truncate(window)
		entry := &Window{
			prev_current: 0,
			curr_current: 1,
			windowStart:  now,
		}
		ms.data[key] = entry
		return 0, 1, ws, nil
	}

	if entry.windowStart.After(entry.windowStart.Add(window)) {
		entry.prev_current = 0
		entry.curr_current = 1
		entry.windowStart = now
	} else if entry.windowStart.After(entry.windowStart) {
		entry.prev_current = entry.curr_current
		entry.curr_current = 0
		entry.windowStart = now
	} else {
		entry.curr_current++
	}

	return entry.prev_current, entry.curr_current, entry.windowStart, nil

}

func (ms *MemoryStore) cleanupMemoryRamUssage(interval_duration time.Duration) {

	ticker := time.NewTicker(interval_duration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ms.revictDeleteData(interval_duration)
		case <-ms.stopCH:
			return
		}
	}

}

func (ms *MemoryStore) revictDeleteData(interval_duration time.Duration) {

	ms.ms.Lock()
	defer ms.ms.Unlock()

	cutoff := time.Now().Add(-2 * interval_duration)

	for key, result := range ms.data {
		if result.windowStart.Before(cutoff) {
			delete(ms.data, key)
		}
	}

}

func (ms *MemoryStore) Stop() {
	close(ms.stopCH)
}
