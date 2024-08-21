package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v4"
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
	router.HandleFunc("/account/{id}", withAuth(makeHTTPHandlerFunc(s.handleGetAccByID), s.storeg))
	router.HandleFunc("/transfer", makeHTTPHandlerFunc(s.handleTransfer))

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

	return fmt.Errorf("method not support")
}

func (s *APIserver) handleGetAcc(w http.ResponseWriter, r *http.Request) error {
	accs, err := s.storeg.GetAccs()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accs)
}

func (s *APIserver) handleGetAccByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id, err := getID(r)
		if err != nil {
			return err
		}

		acc, err := s.storeg.GerAccID(id)
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, acc)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAcc(w, r)
	}
	return fmt.Errorf("wrong method %s", r.Method)
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

	tokenS, err := createJWT(acc)
	if err != nil {
		return err
	}

	fmt.Println(tokenS)

	return WriteJSON(w, http.StatusOK, acc)
}

func (s *APIserver) handleDeleteAcc(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.storeg.DeleteAcc(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIserver) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transfer := new(TransferReq)
	if err := json.NewDecoder(r.Body).Decode(transfer); err != nil {
		return err
	}
	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, transfer)
}

func WriteJSON(w http.ResponseWriter, sts int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(sts)

	return json.NewEncoder(w).Encode(v)
}

func createJWT(acc *Acc) (string, error) {
	claim := &jwt.MapClaims{
		"expiresAt": 15000,
		"accNum":    acc.Number,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString([]byte(secret))
}

//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NOdW0iOjE0OTY1NjYwLCJleHBpcmVzQXQiOjE1MDAwfQ.P-0iwqSlVH_GNOrd3arrYBY_CxzqbYZYkayDAdJv1Y8
//13

func withAuth(handlerFunc http.HandlerFunc, s Srorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling auth")

		tokenS := r.Header.Get("JWTtoken")

		token, err := validJWT(tokenS)
		if err != nil {
			WriteJSON(w, http.StatusOK, ApiError{Error: "invalid token"})
			return
		}

		if !token.Valid {
			WriteJSON(w, http.StatusOK, ApiError{Error: "invalid token(2)"})
			return
		}

		userID, err := getID(r)
		if err != nil {
			WriteJSON(w, http.StatusOK, ApiError{Error: "invalid token(3)"})
			return
		}
		acc, err := s.GerAccID(userID)
		if err != nil {
			WriteJSON(w, http.StatusOK, ApiError{Error: "invalid token(4)"})
			return
		}

		claim := token.Claims.(jwt.MapClaims)

		if acc.Number != int64(claim["accNum"].(float64)) { //Чут-чут костыль, нужен свой claim((
			WriteJSON(w, http.StatusOK, ApiError{Error: "invalid token(5)"})
			return
		}

		handlerFunc(w, r)
	}
}

//const jwrSecret = "qwerty1"

func validJWT(tokerS string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokerS, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
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

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}
