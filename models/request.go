package models

import (
	"errors"
)

type MessageRequest struct {
	UniqueId *string  `json:"unique_id"`
	Sender   string   `json:"sender"`
	To       []string `json:"to"`
	Subject  *string  `json:"subject"`
	Message  string   `json:"message"`
}

func (m MessageRequest) Validate() error {
	if !ValidateEmail(m.Sender) {
		return errors.New("validation 'sender' email error")
	}
	for _, v := range m.To {
		if !ValidateEmail(v) {
			return errors.New("validation 'to' email error")
		}
	}
	if m.Message == "" {
		return errors.New("validation 'message' error")
	}

	return nil
}
