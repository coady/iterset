//go:build go1.27

package iterset

import (
	"fmt"
	"maps"
	"slices"
	"testing"
)

func ExampleMapSet_Equal() {
	k := []string{"b", "a", "b"}
	fmt.Println(Set("a", "b").Equal(slices.Values(k)), Set("a").Equal(k))
	fmt.Println(Set("a").Equal([]string{}), Set("a").Equal(Set("a")))
	// Output:
	// true false
	// false true
}

func ExampleMapSet_SubsetOf() {
	k := []string{"b", "a", "b"}
	fmt.Println(Set("a").SubsetOf(slices.Values(k)), Set("a", "c").SubsetOf(k))
	fmt.Println(Set("a").SubsetOf(Set(k...)), Set("a").SubsetOf([]string{}))
	// Output:
	// true false
	// true false
}

func ExampleMapSet_Intersects() {
	k := []string{"b", "a", "b"}
	fmt.Println(Set("c").Intersects(slices.Values(k)), Set("a").Intersects(Set(k...)))
	// Output: false true
}

func ExampleMapSet_IntersectionWith() {
	s := Set("a", "b", "c")
	s.IntersectionWith(slices.Values([]string{"b", "c", "d"}))
	fmt.Println(s)
	s.IntersectionWith(Set("a"))
	fmt.Println(s)
	// Output:
	// map[b:{} c:{}]
	// map[]
}

func ExampleMapSet_Intersection() {
	m := MapSet[string, int]{"a": 0, "b": 1}
	s := slices.Values([]string{"b", "c", "d"})
	for key, value := range m.Intersection(s) {
		fmt.Println(key, value)
	}
	fmt.Println(maps.Collect(m.Intersection(Count(s))))
	fmt.Println(maps.Collect(Count(s).Intersection(m)))
	// Output:
	// b 1
	// map[b:1]
	// map[b:1]
}

func ExampleMapSet_IntersectionCount() {
	s, k := Set("a", "b"), []string{"b", "c", "d"}
	fmt.Println(s.IntersectionCount(slices.Values(k)), s.IntersectionCount(Set(k...)))
	// Output: 1 1
}

func ExampleMapSet_Difference() {
	s := Set("a", "b")
	k := []string{"b", "c"}
	fmt.Println(maps.Collect(s.Difference(slices.Values(k))))
	fmt.Println(maps.Collect(s.Difference(Set(k...))), maps.Collect(s.Difference(Set[string]())))
	// Output:
	// map[a:{}]
	// map[a:{}] map[a:{} b:{}]
}

func ExampleMapSet_DifferenceWith() {
	s := Set("a", "b", "c")
	s.DifferenceWith(slices.Values([]string{"c", "d"}))
	fmt.Println(s)
	s.DifferenceWith(Set("a", "x"))
	fmt.Println(s)
	// Output:
	// map[a:{} b:{}]
	// map[b:{}]
}

func TestDifferenceWith(t *testing.T) {
	Set("a").DifferenceWith(Set("a", "b"))
}

func ExampleMapSet_SymmetricDifference() {
	k := []string{"b", "c"}
	fmt.Println(slices.Collect(Set("a", "b").SymmetricDifference(slices.Values(k))))
	fmt.Println(Set("a").Equal(Set("b").SymmetricDifference(Set("a", "b"))))
	// Output:
	// [c a]
	// true
}

func ExampleMapSet_Union_map() {
	m := MapSet[string, int]{"a": 0, "b": 1}
	fmt.Println(m.Union(MapSet[string, int]{"b": 2, "c": 3}))
	// Output: map[a:0 b:2 c:3]
}

func ExampleMapSet_UnionWith_map() {
	m := MapSet[string, int]{"a": 0, "b": 1}
	m.UnionWith(MapSet[string, int]{"b": 2, "c": 3})
	fmt.Println(m)
	// Output: map[a:0 b:2 c:3]
}

func TestSymmetricDifference(t *testing.T) {
	for range Set("a", "b").SymmetricDifference(Set("c")) {
		break
	}
	for range Set("a").SymmetricDifference(Set("a", "b")) {
		break
	}
}

func ExampleMapSet_Overlap() {
	s, k := Set("a", "b", "c"), []string{"b", "c", "d"}
	fmt.Println(s.Overlap(slices.Values(k)))
	fmt.Println(s.Overlap(Set(k...)))
	// Output:
	// 1 2 1
	// 1 2 1
}
