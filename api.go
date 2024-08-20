package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIserver struct {
	listenAdr string
	storeg    Srorage
}

func NewAPIserver(listenAdr string, storeg Srorage) *APIserver {
	return &APIserver{
		listenAdr: listenAdr,
		storeg:    storeg,
	}
}

func (s *APIserver) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandlerFunc(s.handleAcc))
	router.HandleFunc("/account/{id}", makeHTTPHandlerFunc(s.handleGetAccByID))

	log.Println("Server running on port: ", s.listenAdr)

	http.ListenAndServe(s.listenAdr, router)
}

func (s *APIserver) handleAcc(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAcc(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAcc(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAcc(w, r)
	}

	return fmt.Errorf("method not support")
}

func (s *APIserver) handleGetAcc(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIserver) handleGetAccByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]

	//db

	fmt.Println(id)

	return WriteJSON(w, http.StatusOK, &Acc{})
}

func (s *APIserver) handleCreateAcc(w http.ResponseWriter, r *http.Request) error {
	createAccReq := new(CreateAccReq)
	if err := json.NewDecoder(r.Body).Decode(createAccReq); err != nil {
		return err
	}

	acc := NewAcc(createAccReq.Name)
	if err := s.storeg.CreateAcc(acc); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, acc)
}

func (s *APIserver) handleDeleteAcc(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIserver) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func WriteJSON(w http.ResponseWriter, sts int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(sts)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
} //Декоратор apiFunc -> handlerFunc
