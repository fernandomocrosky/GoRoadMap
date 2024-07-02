package main

import (
	"errors"
	"log"
)

func throwErrror() error {
	return errors.New("new error detected")
}

func main() {
	err := throwErrror()
	log.Println(err)
}
