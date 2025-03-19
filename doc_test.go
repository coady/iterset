package iterset

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"testing"
)

type item struct{ id string }

func (it item) Id() string { return it.id }

// Intersect a map with a slice of keys, retaining original order.
func Example_intersect() {
	data := map[string]int{"a": 0, "b": 1, "c": 2}
	keys := []string{"d", "c", "b"}
	// With no sets.
	for _, key := range keys {
		value, ok := data[key]
		if ok {
			fmt.Println(key, value)
		}
	}
	// A typical `mapset` would copy `data`, and have no impact on readability.
	s := Set(slices.Collect(maps.Keys(data))...)
	for _, key := range keys {
		if s.Contains(key) {
			fmt.Println(key, data[key])
		}
	}
	// Using an intersect method would also copy `keys`, and lose ordering.
	// Whereas `iterset` methods have in-lined logic with zero copying and lazy iteration.
	for key, value := range Cast(data).Intersect(slices.Values(keys)) {
		fmt.Println(key, value)
	}
	// Output:
	// c 2
	// b 1
	// c 2
	// b 1
	// c 2
	// b 1
}

// Is one slice a superset of another?
func Example_superset() {
	left, right := []string{"a", "b", "c"}, []string{"b", "c", "d"}
	// With no sets.
	m := map[string]bool{}
	for _, c := range left {
		m[c] = true
	}
	isSuperset := true
	for _, c := range right {
		if !m[c] {
			isSuperset = false
			break
		}
	}
	fmt.Println(isSuperset)
	// Or in functional style.
	fmt.Println(!slices.ContainsFunc(right, func(c string) bool { return !m[c] }))
	// A typical `mapset` would copy both slices, which makes early exits irrelevant.
	// Or it only solves half the problem.
	s := Set(left...)
	fmt.Println(!slices.ContainsFunc(right, func(c string) bool { return !s.Contains(c) }))
	// Whereas `iterset` methods have in-lined logic with minimal copying and early exits.
	fmt.Println(Set(left...).IsSuperset(slices.Values(right)))
	// Output:
	// false
	// false
	// false
	// false
}

// Remove duplicates, retaining original order.
func Example_unique() {
	values := []string{"a", "b", "a"}
	// With no sets.
	m := map[string]bool{}
	keys := []string{}
	for _, c := range values {
		if !m[c] {
			keys = append(keys, c)
		}
		m[c] = true
	}
	fmt.Println(keys)
	// A typical `mapset` would either have no impact on readability, or lose ordering.
	// Whereas `iterset` has this built-in with lazy iteration.
	for c := range Unique(slices.Values(values)) {
		fmt.Println(c)
	}
	// What if a `set` is still needed, in addition to ordering.
	idx := Index(slices.Values(values))
	fmt.Println(idx)
	fmt.Println(Sorted(idx)) // keys sorted by value
	// Output:
	// [a b]
	// a
	// b
	// map[a:0 b:1]
	// [a b]
}

func ExampleMapSet_Contains() {
	s := Set("b", "a", "b")
	fmt.Println(s.Contains("a"), s.Contains("b", "c"))
	// Output: true false
}

func ExampleMapSet_Equal() {
	k := slices.Values([]string{"b", "a", "b"})
	fmt.Println(Set("a", "b").Equal(k), Set("a").Equal(k))
	// Output: true false
}

func ExampleEqual() {
	k := slices.Values([]string{"b", "a", "b"})
	s := slices.Values([]string{"a"})
	fmt.Println(Equal(k, slices.Values([]string{"a", "b"})), Equal(k, s), Equal(k, s))
	// Output: true false false
}

func ExampleEqualCounts() {
	k := slices.Values([]string{"b", "a", "b"})
	s := slices.Values([]string{"a", "b"})
	fmt.Println(EqualCounts(k, k), EqualCounts(k, s), EqualCounts(s, k))
	// Output: true false false
}

func ExampleMapSet_IsSubset() {
	k := slices.Values([]string{"b", "a", "b"})
	fmt.Println(Set("a").IsSubset(k), Set("a", "c").IsSubset(k))
	// Output: true false
}

func ExampleIsSubset() {
	s1 := slices.Values([]string{"a"})
	s2 := slices.Values([]string{"a", "b"})
	fmt.Println(IsSubset(s1, s2), IsSubset(s2, s1))
	// Output: true false
}

func ExampleMapSet_IsSuperset() {
	k := slices.Values([]string{"b", "a", "b"})
	fmt.Println(Set("a", "b").IsSuperset(k), Set("a").IsSuperset(k))
	// Output: true false
}

func ExampleMapSet_IsDisjoint() {
	k := slices.Values([]string{"b", "a", "b"})
	fmt.Println(Set("c").IsDisjoint(k), Set("a").IsDisjoint(k))
	// Output: true false
}

func ExampleMapSet_Add() {
	s := Set("a", "b")
	s.Add("b", "c")
	fmt.Println(len(s))
	// Output: 3
}

func ExampleMapSet_Insert() {
	m := MapSet[string, bool]{}
	m.Insert(slices.Values([]string{"b", "a", "b"}), true)
	fmt.Println(m)
	// Output: map[a:true b:true]
}

func ExampleMapSet_Delete() {
	s := Set("a", "b")
	s.Delete("b", "c")
	fmt.Println(len(s))
	// Output: 1
}

func ExampleMapSet_Remove() {
	s := Set("a", "b")
	s.Remove(slices.Values([]string{"b", "c"}))
	fmt.Println(s)
	// Output: map[a:{}]
}

func ExampleMapSet_Toggle() {
	s := Set("a", "b")
	s.Toggle(maps.Keys(Set("b", "c")), struct{}{})
	fmt.Println(s)
	// Output: map[a:{} c:{}]
}

func ExampleMapSet_Union() {
	m := map[string]int{"a": 0, "b": 1}
	n := map[string]int{"b": 2, "c": 3}
	fmt.Println(Cast(m).Union(maps.All(n)))
	// Output: map[a:0 b:2 c:3]
}

func ExampleMapSet_Intersect() {
	m := MapSet[string, int]{"a": 0, "b": 1}
	s := slices.Values([]string{"b", "c"})
	for key, value := range m.Intersect(s) {
		fmt.Println(key, value)
	}
	// Output: b 1
}

func ExampleIntersect() {
	s1 := slices.Values([]string{"a", "b"})
	s2 := slices.Values([]string{"b", "c"})
	fmt.Println(slices.Collect(Intersect(s1, s2)))
	// Output: [b]
}

func ExampleMapSet_Difference() {
	k := slices.Values([]string{"b", "c"})
	fmt.Println(maps.Collect(Set("a", "b").Difference(k)))
	// Output: map[a:{}]
}

func ExampleDifference() {
	s1 := slices.Values([]string{"a", "b"})
	s2 := slices.Values([]string{"b", "c"})
	fmt.Println(slices.Collect(Difference(s1, s2)))
	// Output: [a]
}

func ExampleMapSet_ReverseDifference() {
	k := slices.Values([]string{"b", "c"})
	fmt.Println(slices.Collect(Set("a", "b").ReverseDifference(k)))
	// Output: [c]
}

func ExampleMapSet_SymmetricDifference() {
	k := slices.Values([]string{"b", "c"})
	fmt.Println(slices.Collect(Set("a", "b").SymmetricDifference(k)))
	// Output: [c a]
}

func ExampleCast() {
	m := map[string]bool{}
	Cast(m).Add("a")
	fmt.Println(m)
	// Output: map[a:false]
}

func ExampleUnique() {
	k := slices.Values([]string{"b", "a", "b"})
	fmt.Println(slices.Collect(Unique(k)))
	// Output: [b a]
}

func ExampleUniqueBy() {
	items := []item{{id: "b"}, {id: "a"}, {id: "b"}}
	for key, value := range UniqueBy(slices.Values(items), item.Id) {
		fmt.Println(key, value)
	}
	// Output:
	// b {b}
	// a {a}
}

func ExampleCompact() {
	k := slices.Values([]string{"b", "b", "a", "a", "b"})
	for key, count := range Compact(k) {
		fmt.Println(key, count)
	}
	// Output:
	// b 2
	// a 2
	// b 1
}

func ExampleCompactBy() {
	items := []item{{id: "b"}, {id: "b"}, {id: "a"}, {id: "a"}, {id: "b"}}
	for key, values := range CompactBy(slices.Values(items), item.Id) {
		fmt.Println(key, values)
	}
	// Output:
	// b [{b} {b}]
	// a [{a} {a}]
	// b [{b}]
}

func ExampleCollect() {
	k := slices.Values([]string{"b", "a", "b"})
	fmt.Println(Collect(k, true))
	// Output: map[a:true b:true]
}

func ExampleSet() {
	fmt.Println(Set("b", "a", "b"))
	// Output: map[a:{} b:{}]
}

func ExampleIndex() {
	fmt.Println(Index(slices.Values([]string{"b", "a", "b"})))
	// Output: map[a:1 b:0]
}

func ExampleCount() {
	fmt.Println(Count(slices.Values([]string{"b", "a", "b"})))
	// Output: map[a:1 b:2]
}

func ExampleIndexBy() {
	items := []item{{id: "b"}, {id: "a"}, {id: "b"}}
	fmt.Println(IndexBy(slices.Values(items), item.Id))
	// Output: map[a:{a} b:{b}]
}

func ExampleGroupBy() {
	items := []item{{id: "b"}, {id: "a"}, {id: "b"}}
	fmt.Println(GroupBy(slices.Values(items), item.Id))
	// Output: map[a:[{a}] b:[{b} {b}]]
}

func ExampleMemoize() {
	fmt.Println(Memoize(slices.Values([]string{"b", "a", "b"}), strings.ToUpper))
	// Output: map[a:A b:B]
}

func ExampleSorted() {
	fmt.Println(Sorted(Index(slices.Values([]string{"b", "a", "b"}))))
	// Output: [b a]
}

func ExampleMin() {
	s := Min(map[string]int{"a": 2, "b": 1, "c": 1})
	slices.Sort(s)
	fmt.Println(s)
	// Output: [b c]
}

func ExampleMax() {
	s := Max(map[string]int{"a": 2, "b": 2, "c": 1})
	slices.Sort(s)
	fmt.Println(s)
	// Output: [a b]
}

func ExampleSize() {
	s := slices.Values([]int{0, 0, 0})
	fmt.Println(Size(s))
	// Output: 3
}

func ExampleKeys() {
	s := slices.All([]string{"a", "b", "c"})
	fmt.Println(slices.Collect(Keys(s)))
	// Output: [0 1 2]
}

func TestIter(t *testing.T) {
	k := slices.Values([]string{"a", "A"})
	for range UniqueBy(k, strings.TrimSpace) {
		break
	}
	for range Compact(k) {
		break
	}
	for range CompactBy(k, strings.TrimSpace) {
		break
	}
	for range Set("a").Intersect(k) {
		break
	}
	for c := range Index(k).Difference(k) {
		t.Errorf("should be empty: %s", c)
	}
	for range Set("b").Difference(k) {
		break
	}
	for range Set("b").ReverseDifference(k) {
		break
	}
	for c := range Index(k).SymmetricDifference(k) {
		t.Errorf("should be empty: %s", c)
	}
	for range Set("b").SymmetricDifference(k) {
		break
	}
	for c := range Set("b").SymmetricDifference(k) {
		if c == "b" {
			break
		}
	}
	for range Keys(slices.All([]string{""})) {
		break
	}
}

func TestEmpty(t *testing.T) {
	var m MapSet[string, struct{}]
	if m.Union() == nil {
		t.Error("should not be nil")
	}
	for range m.Intersect(nil) {
		t.Error("should be empty")
	}
	for range Intersect(nil, maps.Keys(m)) {
		t.Error("should be empty")
	}
	for range m.ReverseDifference(maps.Keys(m)) {
		t.Error("should be empty")
	}
	for range m.SymmetricDifference(maps.Keys(m)) {
		t.Error("should be empty")
	}
	m.Delete("")
	m.Remove(slices.Values([]string{""}))
}
