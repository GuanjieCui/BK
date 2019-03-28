package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pegasus/mySQL"
	"pegasus/utils"
)

func Test(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one test request")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var user utils.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "cannot decode user data from client", http.StatusBadRequest)
		fmt.Println("cannot decode user data from client")
		return
	}

	db, err := mySQL.Connect()
	defer db.Close()

	if err != nil {
		http.Error(w, "DB Connection failed", http.StatusInternalServerError)
		fmt.Println("DB Connection failed")
	}

	res := db.Get(user.UserID)

	w.Write([]byte(res))
}