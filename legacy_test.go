//go:build !go1.27

package iterset

import (
	"fmt"
	"maps"
	"slices"
)

func ExampleMapSet_Equal() {
	k := slices.Values([]string{"b", "a", "b"})
	fmt.Println(Set("a", "b").Equal(k), Set("a").Equal(k))
	// Output: true false
}

func ExampleMapSet_IsSubset() {
	k := slices.Values([]string{"b", "a", "b"})
	fmt.Println(Set("a").IsSubset(k), Set("a", "c").IsSubset(k))
	// Output: true false
}

func ExampleMapSet_IsDisjoint() {
	k := slices.Values([]string{"b", "a", "b"})
	fmt.Println(Set("c").IsDisjoint(k), Set("a").IsDisjoint(k))
	// Output: true false
}

func ExampleMapSet_Keep() {
	s := Set("a", "b", "c")
	s.Keep(slices.Values([]string{"b", "c", "d"}))
	fmt.Println(s)
	// Output: map[b:{} c:{}]
}

func ExampleMapSet_Intersect() {
	m := MapSet[string, int]{"a": 0, "b": 1}
	s := slices.Values([]string{"b", "c"})
	for key, value := range m.Intersect(s) {
		fmt.Println(key, value)
	}
	// Output: b 1
}

func ExampleMapSet_IntersectCount() {
	s, k := Set("a", "b"), Set("b", "c", "d")
	fmt.Println(s.IntersectCount(maps.Keys(k)))
	// Output: 1
}

func ExampleMapSet_Difference() {
	k := slices.Values([]string{"b", "c"})
	fmt.Println(maps.Collect(Set("a", "b").Difference(k)))
	// Output: map[a:{}]
}
