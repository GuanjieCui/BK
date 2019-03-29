package main

import (
	"encoding/json"
	"fmt"
	"log"
	
	"net/http"
	
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)
const (
	USERNAME = "root"
	PASSWORD = "root"
	PORT_NUMBER = "3306"
	HOSTNAME = "localhost"
)
type Location struct{
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
type Order struct{
    UserId uint16 `json:"UserId"`
    Size string `json:"Size"`
    arrivalTime string `json:"arrival"`
    Weight float64 `json:"Weight"`
    PickupLoc Location `json:"PickupLoc"`
    DropoffLoc Location `json:"DropoffLoc"`
    
}
type resp struct{
	OrderId uint16
}
func main() {
	fmt.Println("started-service")
	http.HandleFunc("/order", handlerOrder)
	//http.HandleFunc("/track", handlerTrack)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlerOrder(w http.ResponseWriter,r *http.Request){
    fmt.Println("Receive order")
    decoder:=json.NewDecoder(r.Body)
    var p Order
    if err:=decoder.Decode(&p);err!=nil{
        panic(err)
    }
    //fmt.Fprintf(w,"Order received: %d,%s,%f",p.UserId,p.Size,p.Weight)
    db, err := sql.Open("mysql", USERNAME + ":" + PASSWORD + "@tcp(" +
		HOSTNAME + ":" + PORT_NUMBER + ")/pegasus")
    fmt.Println("Receive order2")
    if err != nil {
		panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	// close the db
	defer db.Close()
	querySize:="select count(*) from orders"
	qsz,_:=db.Query(querySize)
    fmt.Println("Receive order3")
	var sz uint16
	if qsz.Next(){
		if err := qsz.Scan(&sz); err != nil {
			fmt.Println("err", err)
	    }
	}
	
	q,err:=db.Prepare("insert into orders values(?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	defer q.Close()
    fmt.Println(p.arrivalTime)
	_,err=q.Exec(sz+1,p.UserId,p.Size,p.Weight,p.PickupLoc.Lat,p.PickupLoc.Lon,p.DropoffLoc.Lat,p.DropoffLoc.Lon,p.arrivalTime)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	res:=resp{
		OrderId: sz+1,
	}
	b,err:=json.Marshal(res)
	if err!=nil{
		fmt.Println("error:",err)
	}
	w.Write(b)
}