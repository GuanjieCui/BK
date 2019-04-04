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
	RobotId uint16 `json:"RobotId"`
    Username uint16 `json:"Username"`
    Size string `json:"Size"`
    ArrivalTime string `json:"Arrival"`
    Weight float64 `json:"Weight"`
    PickupLoc Location `json:"PickupLoc"`
    DropoffLoc Location `json:"DropoffLoc"`
    
}
type resp struct{
	OrderId uint16
}
type arrival struct{
	ArrivalTime string
}
func main() {
	fmt.Println("started-service")
	jwtMiddleware:=jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter:func(token *jwt.Token) (interface{},error) {
			return []byte(mySigningKey),nil
		},
		SigningMethod:jwt.SigningMethodHS256,
	})
	r:=mux.NewRouter()
	r.Handle("/order",jwtMiddleware.Handler(http.HandlerFunc(handlerOrder))).Methods("POST","OPTIONS")
	r.Handle("/track",jwtMiddleware.Handler(http.HandlerFunc(handlerTrack))).Methods("GET","OPTIONS")
	http.Handle("/",r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlerOrder(w http.ResponseWriter,r *http.Request){
    fmt.Println("Receive order")
    w.Header().Set("Content-Type","application/json")
    w.Header().Set("Access-Control-Allow-Origin","*")
    w.Header().Set("Access-Control-Allow-Headers","Content-Type,Authorization")

    user:=r.Context().Value("user")
    claims:=user.(*jwt.Token).Claims
    username:=claims.(jwt.MapClaims)["Username"]

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
    fmt.Println(p.ArrivalTime)
	_,err=q.Exec(sz+1,username,p.RobotId,p.Size,p.Weight,p.PickupLoc.Lat,p.PickupLoc.Lon,p.DropoffLoc.Lat,p.DropoffLoc.Lon,p.ArrivalTime)
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

func handlerTrack(w http.ResponseWriter,r *http.Request){
	fmt.Println("Received one request for track")
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type,Authorization")

	orderId,_:strconv.ParseInt(r.URL.Query().Get("orderid"),10,32)
	db, err := sql.Open("mysql", USERNAME + ":" + PASSWORD + "@tcp(" +
		HOSTNAME + ":" + PORT_NUMBER + ")/pegasus")
    if err != nil {
		panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	// close the db
	defer db.Close()
	query,err:=db.Prepare("select ArrivalTime from orders where OrderId=?")
	if err != nil {
		panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	defer query.Close()
	q,_:=query.Exec(orderId)
    
	var atime string
	if q.Next(){
		if err := q.Scan(&atime); err != nil {
			fmt.Println("err", err)
	    }
	}
	res:=arrival{
		ArrivalTime:atime,
	}
	b,err:=json.Marshal(res)
	if err!=nil{
		fmt.Println("error:",err)
	}
	w.Write(b)

}