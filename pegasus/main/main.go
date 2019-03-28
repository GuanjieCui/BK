package main

import (
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"pegasus/handler"
)

var mySigningKey = []byte("a")

func main() {
	//mySQL.NewTable()
	//db, _ := mySQL.Connect()
	//defer db.Close()

	fmt.Println("fire-up-engine")
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	r := mux.NewRouter()
	r.Handle("/signup", http.HandlerFunc(handler.Signup)).Methods("POST", "OPTIONS")
	r.Handle("/login", http.HandlerFunc(handler.Login)).Methods("POST", "OPTIONS")
	// test handle
	r.Handle("/testget", jwtMiddleware.Handler(http.HandlerFunc(handler.Test))).Methods("GET", "OPTIONS")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

