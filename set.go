package sorted

import (
	"cmp"
	"encoding/json"
	"iter"
	"slices"
)

// A sorted set implementation. Keeps keys in a sorted vector to improve memory
// locality.
type Set[T cmp.Ordered] struct {
	vals []T
}

func toSliceSet[S ~[]E, E cmp.Ordered](s S) S {
	slices.Sort(s)
	return slices.Compact(s)
}

// Creates a new empty sorted set, or populated with the given arguments
func NewSet[T cmp.Ordered](items ...T) *Set[T] {
	return &Set[T]{
		vals: toSliceSet(items),
	}
}

func (s Set[T]) keySearch(key T) (int, bool) {
	return slices.BinarySearch(s.vals, key)
}

func (s Set[T]) Contains(val T) bool {
	_, found := s.keySearch(val)
	return found
}

func (s *Set[T]) Put(val T) bool {
	idx, found := s.keySearch(val)
	if found {
		return false
	}

	// Append place holder, resize value array + put new value
	s.vals = append(s.vals, val)
	copy(s.vals[idx+1:], s.vals[idx:])
	s.vals[idx] = val

	return true
}

func (s Set[T]) Values() []T {
	return s.vals
}

func (s Set[T]) Items() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, val := range s.vals {
			if !yield(val) {
				return
			}
		}
	}
}

func (s Set[T]) Len() int {
	return len(s.vals)
}

func (s *Set[T]) Del(val T) bool {
	idx, found := slices.BinarySearch(s.vals, val)
	if !found {
		return false
	}

	s.vals = slices.Delete(s.vals, idx, idx+1)

	return true
}

// Custom JSON marshalling, to make its use transparent with Go's json package.
func (s *Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.vals)
}

// Custom JSON marshalling, to make its use transparent with Go's json package.
func (s *Set[T]) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &s.vals); err != nil {
		return err
	}

	// Maintain ordered set invariant:
	s.vals = toSliceSet(s.vals)
	return nil
}
