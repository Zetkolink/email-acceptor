package models

type Success struct {
	Id string `json:"id"`
}

type Failed struct {
	Error string `json:"error"`
}
