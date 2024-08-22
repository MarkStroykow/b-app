package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginReq struct {
	Num  int64  `json:"num"`
	Pass string `json:"pass"`
}

type TransferReq struct {
	ToAcc int `json:"toAcc"`
	Sum   int `json:"sum"`
}

type CreateAccReq struct {
	Name     string `json:"name"`
	Password string `json:"pass"`
}

type Acc struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Number    int64     `json:"num"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	Bal       int       `json:"bal"`
}

func NewAcc(name, password string) (*Acc, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Acc{
		Name:      name,
		Number:    int64(rand.Intn(99999999)),
		Password:  string(encpw),
		CreatedAt: time.Now().UTC(),
	}, nil
}
