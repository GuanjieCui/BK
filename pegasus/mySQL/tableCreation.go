package mySQL

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// delete elder pegasus mySQL and create a new one
func NewTable() {
	fmt.Println("start to create a new pegasus datasbase ......")
	db, err := sql.Open("mysql", USERNAME + ":" + PASSWORD + "@tcp(" +
		HOSTNAME + ":" + PORT_NUMBER + ")/")
	// close the db
	defer db.Close()
	checkErr(err)
	// drop pre pegasus database
	q := "drop database if exists pegasus"
	db.Query(q)

	// create users table
	q = "create table pegasus.users (" +
		"user_id varchar(255) not null," +
		"password varchar(255) not null," +
		"first_name varchar(255) not null," +
		"last_name varchar(255) not null," +
		"primary key (user_id)" +
		")"
	db.Query(q)
	// create items table
	q = "create table pegasus.items (" +
		"item_id varchar(255) not null," +
		"description varchar(255) not null," +
		"order_datetime datetime not null," +
		"price float not null," +
		"primary key (item_id)" +
		")"
	db.Query(q)
	// create packages table to record delivery info of items
	q = "create table pegasus.packages (" +
		"item_id varchar(255) not null," +
		"from_lat float not null," +
		"from_lon float not null," +
		"to_lat float not null," +
		"to_lon float not null," +
		"drop_time datetime not null," +
		"delivery_time datetime," +
		"primary key (item_id)," +
		"foreign key (item_id) REFERENCES items(item_id)" +
		")"
	db.Query(q)
	// create station table to record all post robots info
	q = "create table pegasus.station (" +
		"robot_id varchar(255) not null," +
		"robot_type varchar(255) not null," +
		"robot_status varchar(255) not null," +
		"primary key (robot_id)" +
		")"
	db.Query(q)

	// insert a fake user
	// "INSERT INTO users VALUES('1111', '3229c1097c00d497a0fd282d586be050', 'John', 'Smith')";
	q = "insert into pegasus.users values ('1111', '3229c1097c00d497a0fd282d586be050', 'thanks', 'plane')"
	db.Query(q)
	fmt.Println("pegasus express datatbase created!")
}

func checkErr(err error) {
	if (err != nil) {
		panic(err)
	}
}
