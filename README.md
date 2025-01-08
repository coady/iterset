[![build](https://github.com/coady/iterset/actions/workflows/build.yml/badge.svg)](https://github.com/coady/iterset/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/coady/iterset/branch/main/graph/badge.svg)](https://codecov.io/gh/coady/iterset/)
[![go report](https://goreportcard.com/badge/github.com/coady/iterset)](https://goreportcard.com/report/github.com/coady/iterset)

# iterset
A [set](https://en.wikipedia.org/wiki/Set_(abstract_data_type)) library based on maps and [iterators](https://pkg.go.dev/iter). A set type is not necessary to have set operations.

There are many `mapset` implementations available, but they restrict the values to `struct{}` or `bool`. In practice this has downsides.
* Maps must be copied even though they already support iteration and O(1) lookup.
* Map values are lost.
* Slices must be copied even if they would have only been iterated.
* Slice ordering is lost.
* Copying effectively means no early exits, e.g., in a subset check.

Since sets are not built-in, they realistically will always be a secondary type. Even in languages with built-in sets, it is common to call set operations on keys while still keeping data in a map, and common to want to retain ordering.

So `iterset` is built around generic maps with `any` value type. Inspired by [Python sets](https://docs.python.org/3/library/stdtypes.html#set-types-set-frozenset), its methods support iterators. This integrates well with functions in [maps](https://pkg.go.dev/maps) and [slices](https://pkg.go.dev/slices), and addresses the typical `mapset` issues.
* Maps can be casted instead of copied.
* Map values are kept without affecting set operations.
* Slices can be iterated using `slices.Values` without copying.
* Slice iterators retain ordering.
* Iterators are lazily evaluated, inherently supporting early exits.

## Usage
There are constructors for all common use cases.
* `Cast` a map
* `Unique{By}` iterates keys in order
* `Compact{By}` iterates consecutive grouped keys
* `Collect` with default value
* `Set` from variadic args
* `Index` retains original position
* `Count` stores key counts
* `IndexBy` stores values by key function
* `GroupBy` stores slices grouped by key function
* `Memoize` caches function call

Methods support iterators, compatible with `slices.Values` and `maps.Keys`. Implementations are asymptotically optimal, and exit early where relevant.
* `Equal`
* `IsSubset`
* `IsSuperset`
* `IsDisjoint`
* `Union`
* `Intersect`
* `Difference`
* `ReverseDifference`
* `SymmetricDifference`

Scalar operations can be passed as bound methods for functional programming. 
* `Get`
* `Contains`
* `Missing`
* `Add`
* `Delete`

Some operations are better expressed as functions, to avoid making unnecessary maps.
* `Sorted`
* `IsSubset`
* `Intersect`
* `Difference`

## Installation
No dependencies. [Go >=1.23](https://go.dev/doc/go1.23) required.

```console
go get github.com/coady/iterset
```

## Tests
100% code coverage.

```console
go test -cover
```
