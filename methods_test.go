//go:build go1.27

package iterset

import (
	"fmt"
	"slices"
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
