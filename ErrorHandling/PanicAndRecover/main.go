package main

import "fmt"

func handlePanic() {
	a := recover()
	if a != nil {
		fmt.Printf("Recover: %v\n", a)
	}
}

func runPanic() {
	defer handlePanic()
	panic("This is a panic")
}

func main() {
	fmt.Println("Start")
	runPanic()
	fmt.Println("End")
}
