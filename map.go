package sorted

import (
	"bytes"
	"cmp"
	"encoding/json"
	"iter"
	"slices"
)

// A sorted map implementation. Keeps keys and values in separate arrays,
// to improve memory locality when searching through the key space.
type Map[K cmp.Ordered, V any] struct {
	keys []K
	vals []V
}

// Creates a new, empty sorted map
func NewMap[K cmp.Ordered, V any]() *Map[K, V] {
	return &Map[K, V]{
		keys: []K{},
		vals: []V{},
	}
}

func (m Map[K, V]) keySearch(key K) (int, bool) {
	return slices.BinarySearch(m.keys, key)
}

// Checks if the key is present in the map
func (m Map[K, V]) Contains(key K) bool {
	_, found := m.keySearch(key)
	return found
}

// Returns the value in the map, if found and an indicator whether found or not.
// If not found the value will be the 0-value
func (m Map[K, V]) Get(key K) (V, bool) {
	if idx, found := m.keySearch(key); found {
		return m.vals[idx], true
	}

	var empty V
	return empty, false
}

// Sets a key-value pair at a specific key slot, returns true if the key was not
// already in the key-set.
func (m *Map[K, V]) Put(key K, value V) bool {
	idx, found := m.keySearch(key)
	if found {
		m.vals[idx] = value
		return false
	}

	// Append place holder, resize key array + put new key
	m.keys = append(m.keys, key)
	copy(m.keys[idx+1:], m.keys[idx:])
	m.keys[idx] = key

	// Append place holder, resize value array + put new value
	m.vals = append(m.vals, value)
	copy(m.vals[idx+1:], m.vals[idx:])
	m.vals[idx] = value

	return true
}

// Iterates through the key-value set, in key order
func (m Map[K, V]) Items() iter.Seq2[K, *V] {
	return func(yield func(k K, v *V) bool) {
		for idx, k := range m.keys {
			if !yield(k, &m.vals[idx]) {
				return
			}
		}
	}
}

// Returns list of all keys
func (m Map[K, V]) Keys() []K {
	return m.keys
}

// Returns list of all values
func (m Map[K, V]) Vals() []V {
	return m.vals
}

// Returns the number of element in the map
func (m Map[K, V]) Len() int {
	return len(m.keys)
}

// Deletes the given key, returns true if the key was found, false otherwise
func (m *Map[K, V]) Del(key K) bool {
	idx, found := slices.BinarySearch(m.keys, key)
	if !found {
		return false
	}

	m.keys = slices.Delete(m.keys, idx, idx+1)
	m.vals = slices.Delete(m.vals, idx, idx+1)

	return true
}

// Custom JSON marshalling, to make its use transparent with Go's json package.
func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	i := 0
	count := m.Len()

	if err := buf.WriteByte(byte('{')); err != nil {
		return nil, err
	}

	for k, v := range m.Items() {
		keyBytes, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}

		if keyBytes[0] != byte('"') {
			if err := buf.WriteByte(byte('"')); err != nil {
				return nil, err
			}
		}

		_, err = buf.Write(keyBytes)
		if err != nil {
			return nil, err
		}

		if keyBytes[0] != byte('"') {
			if err := buf.WriteByte(byte('"')); err != nil {
				return nil, err
			}
		}

		if err := buf.WriteByte(byte(':')); err != nil {
			return nil, err
		}

		valBytes, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		_, err = buf.Write(valBytes)
		if err != nil {
			return nil, err
		}

		i += 1
		if i < count {
			if err := buf.WriteByte(byte(',')); err != nil {
				return nil, err
			}
		}
	}

	if err := buf.WriteByte(byte('}')); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
