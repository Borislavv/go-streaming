package entity

type User struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Birthday string `json:"birthday" bson:"birthday,omitempty"`
	Email    string `json:"email" bson:"email"` // unique key
}
