package sorted

import (
	"encoding/json"
	"slices"
	"testing"
)

func assertEqual[V comparable](t *testing.T, expected, actual V) {
	t.Helper()

	if expected != actual {
		t.Errorf("expected `%v` got `%v`", expected, actual)
	}
}

func assertTrue(t *testing.T, val bool) {
	t.Helper()

	if val == false {
		t.Errorf("expected true, got false")
	}
}

func TestSortedMapBasics(t *testing.T) {
	m := NewMap[int, string]()

	assertEqual(t, 0, m.Len())

	_, found := m.Get(4711)
	assertTrue(t, !found)

	m.Put(4711, "foo")
	assertEqual(t, 1, m.Len())
	assertTrue(t, m.Contains(4711))
	v, found := m.Get(4711)
	assertTrue(t, found)
	assertEqual(t, "foo", v)

	m.Put(4711, "hej")
	v, found = m.Get(4711)
	assertTrue(t, found)
	assertEqual(t, "hej", v)

	m.Put(4712, "bar")
	assertEqual(t, 2, m.Len())
	assertTrue(t, m.Contains(4712))

	assertTrue(t, slices.Equal([]int{4711, 4712}, m.Keys()))

	var keys []int
	var vals []*string

	for k, v := range m.Items() {
		keys = append(keys, k)
		vals = append(vals, v)
	}

	found = m.Del(4711)
	assertTrue(t, found)
	found = m.Del(4711)
	assertTrue(t, !found)

	assertTrue(t, slices.Equal([]int{4712}, m.Keys()))
}

func TestSortedMapJSON(t *testing.T) {
	m := NewMap[int, string]()
	m.Put(4711, "foo")
	m.Put(4713, "hej")
	m.Put(4712, "bar")

	bytes, err := json.Marshal(&m)
	if err != nil {
		t.Errorf("expected err to nil, got %v", err)
	}

	assertEqual(t, `{"4711":"foo","4712":"bar","4713":"hej"}`, string(bytes))

	m2 := NewMap[string, int]()
	m2.Put("foo", 4711)
	m2.Put("hej", 4713)
	m2.Put("bar", 4712)

	bytes2, err := json.Marshal(&m2)
	if err != nil {
		t.Errorf("expected err to nil, got %v", err)
	}

	assertEqual(t, `{"bar":4712,"foo":4711,"hej":4713}`, string(bytes2))
}
