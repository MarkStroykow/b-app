package main

import (
	"flag"
	"fmt"
	"log"
)

func seedAcc(storg Srorage, name, pass string) *Acc {
	acc, err := NewAcc(name, pass)
	if err != nil {
		log.Fatal(err)
	}

	if err := storg.CreateAcc(acc); err != nil {
		log.Fatal(err)
	}

	return acc
}

func seedAccs(s Srorage) {
	seedAcc(s, "Mark", "qwerty111")
}

func main() {

	seed := flag.Bool("seed", false, "seed the db")
	flag.Parse()

	storge, err := NewStorageDB()
	if err != nil {
		log.Fatal(err)
	}

	if err := storge.Init(); err != nil {
		log.Fatal(err)
	}
	if *seed {
		fmt.Println("seed from db")
		seedAccs(storge)
	}

	server := NewAPIserver(":8000", storge)
	server.Run()
}
