package main

import (
	"log"
)

func main() {
	storge, err := NewStorageDB()
	if err != nil {
		log.Fatal(err)
	}

	if err := storge.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIserver(":8000", storge)
	server.Run()
}
