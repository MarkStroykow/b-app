package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Srorage interface {
	CreateAcc(*Acc) error
	DeleteAcc(int) error
	UpdateAcc(*Acc) error
	GetAccs() ([]*Acc, error)
	GerAccID(int) (*Acc, error)
}

type StorageDB struct {
	db *sql.DB
}

func NewStorageDB() (*StorageDB, error) {
	connStr := "user=postgres dbname=postgres password=q sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &StorageDB{
		db: db,
	}, nil
}

func (s *StorageDB) Init() error {
	return s.createAccTable()
}

func (s *StorageDB) createAccTable() error {
	req := `create table if not exists acc (
	id serial primary key,
	name varchar(255),
	num serial,
	created_at timestamp
	)`
	_, err := s.db.Exec(req)
	return err
}

func (s *StorageDB) CreateAcc(acc *Acc) error {
	q := `insert into acc 
	(name, num, created_at)
	values ($1,$2,$3)`
	resp, err := s.db.Query(
		q,
		acc.Name,
		acc.Number,
		acc.CreatedAt)

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)

	return nil
}

func (s *StorageDB) UpdateAcc(*Acc) error {
	return nil
}

func (s *StorageDB) DeleteAcc(in int) error {
	return nil
}

func (s *StorageDB) GetAccs() ([]*Acc, error) {
	rows, err := s.db.Query("select * from acc")
	if err != nil {
		return nil, err
	}

	accs := []*Acc{}
	for rows.Next() {
		acc := new(Acc)
		if err := rows.Scan(
			&acc.ID,
			&acc.Name,
			&acc.Number,
			&acc.CreatedAt); err != nil {
			return nil, err
		}
		accs = append(accs, acc)
	}

	return accs, nil
}

func (s *StorageDB) GerAccID(id int) (*Acc, error) {
	return nil, nil
}
