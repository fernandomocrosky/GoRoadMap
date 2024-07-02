package main

import (
	"log"
)

func executeSomehting() {
	defer log.Println("This will execute after the function completes")

	log.Println("This will execute before the defer statement is executed")
}

func main() {
	executeSomehting()
}
