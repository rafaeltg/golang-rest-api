package models

type Customer struct {
	Id     string  `json:"id" bson:"id"`
	Name   string  `json:"name"`
	Phones []Phone `json:"phones"`
}

type Phone struct {
	Type   string `json:"type"`
	Number string `json:"number"`
}
