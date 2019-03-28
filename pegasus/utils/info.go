package utils

// define types that can be represented as entries of database
// like user, item ...

type User struct {
	UserID string `json:"user_id"` // why Username should be uppercase ?
	Password string `json:"password"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
}

type Item struct {

}