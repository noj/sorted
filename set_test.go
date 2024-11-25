package sorted

import (
	"encoding/json"
	"slices"
	"testing"
)

func TestSortedSet(t *testing.T) {
	s := NewSet[string]()

	assertEqual(t, 0, s.Len())
	assertTrue(t, !s.Contains("foo"))

	ok := s.Put("foo")
	assertTrue(t, ok)
	assertTrue(t, s.Contains("foo"))

	ok = s.Put("foo")
	assertTrue(t, !ok)
	assertTrue(t, s.Contains("foo"))

	ok = s.Put("bar")
	assertTrue(t, ok)
	assertTrue(t, s.Contains("bar"))

	assertTrue(t, slices.Equal([]string{"bar", "foo"}, s.Values()))

	var res []string
	for x := range s.Items() {
		res = append(res, x)
	}
	assertTrue(t, slices.Equal([]string{"bar", "foo"}, res))

	ok = s.Del("bar")
	assertTrue(t, ok)
	assertTrue(t, !s.Contains("bar"))

	ok = s.Del("bar")
	assertTrue(t, !ok)
	assertTrue(t, !s.Contains("bar"))

	assertTrue(t, slices.Equal([]string{"foo"}, s.Values()))

	ok = s.Del("foo")
	assertTrue(t, ok)
	assertTrue(t, !s.Contains("foo"))
	assertTrue(t, slices.Equal([]string{}, s.Values()))

	// Var args ctor
	s2 := NewSet(4, 3, 4, 4)
	assertTrue(t, slices.Equal([]int{3, 4}, s2.Values()))
}

func TestSortedSetJSON(t *testing.T) {
	s := NewSet[int]()
	s.Put(5)
	s.Put(3)
	s.Put(1)

	bytes, err := json.Marshal(s)
	if err != nil {
		t.Errorf("expected err to nil, got %v", err)
	}

	assertEqual(t, "[1,3,5]", string(bytes))

	in := NewSet[int]()
	err = json.Unmarshal([]byte("[7,8,1,1,1,1]"), &in)
	if err != nil {
		t.Errorf("expected err to nil, got %v", err)
	}

	assertTrue(t, slices.Equal([]int{1, 7, 8}, in.Values()))
}
