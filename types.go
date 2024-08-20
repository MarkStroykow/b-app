package main

import (
	"math/rand"
	"time"
)

type CreateAccReq struct {
	Name string `json:"name"`
}

type Acc struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Number    int64     `json:"num"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAcc(name string) *Acc {
	return &Acc{
		//ID:        rand.Intn(9999),
		Name:      name,
		Number:    int64(rand.Intn(99999999)),
		CreatedAt: time.Now().UTC(),
	}
}
