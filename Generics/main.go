package main

import "fmt"

func SumInts(m map[string]int64) int64 {
	var s int64

	for _, v := range m {
		s += v
	}

	return s
}

func SumFloats(m map[string]float64) float64 {
	var s float64

	for _, v := range m {
		s += v
	}

	return s
}

func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}

	return s
}

type Number interface {
	int64 | float64
}

func SumNumbers[K comparable, V Number](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}

	return s
}

func main() {

	ints := map[string]int64{
		"first":  1,
		"second": 2,
	}

	floats := map[string]float64{
		"first":  2.0,
		"second": 3.0,
	}

	fmt.Printf("Non generics sums %v and %v\n", SumInts(ints), SumFloats(floats))

	// Just to remember the type can be ommited SumIntsOrFloats(ints) for example
	fmt.Println("Generic sum of integers", SumIntsOrFloats[string, int64](ints))
	fmt.Println("Generic sum of floats", SumIntsOrFloats[string, float64](floats))

	fmt.Println("Generic sum of integers with interface", SumNumbers(ints))
	fmt.Println("Generic sum of floats with interface", SumNumbers(floats))
}
