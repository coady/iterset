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

func ExampleMapSet_IsSubset() {
	k := []string{"b", "a", "b"}
	fmt.Println(Set("a").IsSubset(slices.Values(k)), Set("a", "c").IsSubset(k))
	fmt.Println(Set("a").IsSubset(Set(k...)), Set("a").IsSubset([]string{}))
	// Output:
	// true false
	// true false
}

func ExampleMapSet_IsDisjoint() {
	k := []string{"b", "a", "b"}
	fmt.Println(Set("c").IsDisjoint(slices.Values(k)), Set("a").IsDisjoint(Set(k...)))
	// Output: true false
}

func ExampleMapSet_Keep() {
	s := Set("a", "b", "c")
	s.Keep(slices.Values([]string{"b", "c", "d"}))
	fmt.Println(s)
	s.Keep(Set("a"))
	fmt.Println(s)
	// Output:
	// map[b:{} c:{}]
	// map[]
}

func ExampleMapSet_Intersect() {
	m := MapSet[string, int]{"a": 0, "b": 1}
	s := slices.Values([]string{"b", "c", "d"})
	for key, value := range m.Intersect(s) {
		fmt.Println(key, value)
	}
	fmt.Println(maps.Collect(m.Intersect(Count(s))))
	fmt.Println(maps.Collect(Count(s).Intersect(m)))
	// Output:
	// b 1
	// map[b:1]
	// map[b:1]
}

func ExampleMapSet_IntersectCount() {
	s, k := Set("a", "b"), []string{"b", "c", "d"}
	fmt.Println(s.IntersectCount(slices.Values(k)), s.IntersectCount(Set(k...)))
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

func ExampleMapSet_Remove() {
	s := Set("a", "b", "c")
	s.Remove(slices.Values([]string{"c", "d"}))
	fmt.Println(s)
	s.Remove(Set("a", "x"))
	fmt.Println(s)
	// Output:
	// map[a:{} b:{}]
	// map[b:{}]
}

func TestRemovet(t *testing.T) {
	Set("a").Remove(Set("a", "b"))
}

func ExampleMapSet_SymmetricDifference() {
	k := []string{"b", "c"}
	fmt.Println(slices.Collect(Set("a", "b").SymmetricDifference(slices.Values(k))))
	fmt.Println(Set("a").Equal(Set("b").SymmetricDifference(Set("a", "b"))))
	// Output:
	// [c a]
	// true
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
