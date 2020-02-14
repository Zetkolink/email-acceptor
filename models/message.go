package models

import "github.com/jinzhu/gorm"

type Message struct {
	Sender  string  `json:"sender"`
	To      string  `json:"to"`
	Subject *string `json:"subject"`
	Message string  `json:"message"`
	gorm.Model
}
