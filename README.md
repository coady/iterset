[![build](https://github.com/coady/iterset/actions/workflows/build.yml/badge.svg)](https://github.com/coady/iterset/actions/workflows/build.yml)
[![image](https://codecov.io/gh/coady/iterset/branch/main/graph/badge.svg)](https://codecov.io/gh/coady/iterset/)

# iterset
[Golang](https://go.dev) set operations using maps and iterators.

There are many `mapset` implementations available, but they restrict the values to `struct{}` or `bool`. This means slices, maps, and iterators have to be converted to sets. Which besides being inefficient, loses slice ordering and map values. Additionally since sets are not built-in, they realistically will always be secondary types. Even in languages with built-in sets, it is common to call set operations on keys while still retaining data in a map.

So `iterset` is built around generic maps with `any` value type. Maps can be casted without copying, and the methods support set operations which integrate with functions in [maps](https://pkg.go.dev/maps) and [slices](https://pkg.go.dev/slices). Additionally the methods support iterators for further efficiency. Inspired by [Python sets](https://docs.python.org/3/library/stdtypes.html#set-types-set-frozenset), which allow iterables in methods.

## Usage
There are constructors for all common use cases.
* `Cast` a map
* `Unique` iterates keys in order
* `Collect` with default value
* `Set` from variadic args
* `Index` retains original position
* `Count` stores key counts
* `IndexBy` stores values by key function
* `GroupBy` stores slices grouped by key function
* `Memoize` caches function call

Methods support iterators, compatible with `slices.Values` and `maps.Keys`.
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

## Installation
No dependencies. [Go 1.23](https://go.dev/doc/go1.23) required.

```console
go get github.com/coady/iterset
```

## Tests
```console
go test .
```
