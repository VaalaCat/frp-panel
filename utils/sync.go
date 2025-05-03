package utils

// from https://gist.github.com/tarampampam/f96538257ff125ab71785710d48b3118

import "sync"

// SyncMap is like a Go sync.Map but type-safe using generics.
//
// The zero SyncMap is empty and ready for use. A SyncMap must not be copied after first use.
type SyncMap[K comparable, V any] struct {
	mu sync.Mutex
	m  map[K]V
}

// Grow grows the map to the given size. It can be called before the first write operation used.
func (s *SyncMap[K, V]) Grow(size int) {
	s.mu.Lock()
	s.grow(size)
	s.mu.Unlock()
}

func (s *SyncMap[K, V]) grow(size ...int) {
	if s.m == nil {
		if len(size) == 0 {
			s.m = make(map[K]V) // let runtime decide the needed map size
		} else {
			s.m = make(map[K]V, size[0])
		}
	}
}

// Clone returns a copy (clone) of current SyncMap.
func (s *SyncMap[K, V]) Clone() SyncMap[K, V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	var clone = make(map[K]V, len(s.m))

	for k, v := range s.m {
		clone[k] = v
	}

	return SyncMap[K, V]{m: clone}
}

// Load returns the value stored in the map for a key, or nil if no value is present.
// The ok result indicates whether value was found in the map.
func (s *SyncMap[K, V]) Load(key K) (value V, loaded bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.m == nil { // fast operation terminator
		return
	}

	value, loaded = s.m[key]

	return
}

// Store sets the value for a key.
func (s *SyncMap[K, V]) Store(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.grow()

	s.m[key] = value
}

// LoadOrStore returns the existing value for the key if present. Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (s *SyncMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if actual, loaded = s.m[key]; !loaded {
		s.grow()

		s.m[key], actual = value, value
	}

	return
}

// LoadAndDelete deletes the value for a key, returning the previous value if any. The loaded result reports whether
// the key was present.
func (s *SyncMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.m == nil { // fast operation terminator
		return
	}

	s.grow()

	if value, loaded = s.m[key]; loaded {
		delete(s.m, key)
	}

	return
}

// Delete deletes the value for a key.
func (s *SyncMap[K, V]) Delete(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.m == nil { // fast operation terminator
		return
	}

	s.grow()

	delete(s.m, key)
}

// Range calls f sequentially for each key and value present in the map. If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's contents: no key will be visited more
// than once. Range does not block other methods on the receiver; even f itself may call any method on m.
func (s *SyncMap[K, V]) Range(f func(key K, value V) (shouldContinue bool)) {
	s.mu.Lock()

	if s.m == nil { // fast operation terminator
		s.mu.Unlock()

		return
	}

	s.grow()

	for k, v := range s.m {
		s.mu.Unlock()

		if !f(k, v) {
			return
		}

		s.mu.Lock()
	}

	s.mu.Unlock()
}

// Len returns the count of values in the map.
func (s *SyncMap[K, V]) Len() (l int) {
	s.mu.Lock()
	l = len(s.m)
	s.mu.Unlock()

	return
}

// Keys return slice with all map keys.
func (s *SyncMap[K, V]) Keys() []K {
	s.mu.Lock()
	defer s.mu.Unlock()

	var keys, i = make([]K, len(s.m)), 0

	for k := range s.m {
		keys[i], i = k, i+1
	}

	return keys
}

// Values return slice with all map values.
func (s *SyncMap[K, V]) Values() []V {
	s.mu.Lock()
	defer s.mu.Unlock()

	var values, i = make([]V, len(s.m)), 0

	for _, v := range s.m {
		values[i], i = v, i+1
	}

	return values
}

func (s *SyncMap[K, V]) ToMap() map[K]V {
	s.mu.Lock()
	defer s.mu.Unlock()

	var m = make(map[K]V, len(s.m))

	for k, v := range s.m {
		m[k] = v
	}

	return m
}
