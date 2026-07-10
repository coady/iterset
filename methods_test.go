//go:build go1.27

package iterset

import (
	"fmt"
	"slices"
)

func ExampleMapSet_IsSubset() {
	k := []string{"b", "a", "b"}
	fmt.Println(Set("a").IsSubset(slices.Values(k)), Set("a", "c").IsSubset(k))
	fmt.Println(Set("a").IsSubset(Set(k...)), Set("a", "c").IsSubset(map[string]struct{}{}))
	// Output:
	// true false
	// true false
}
