package rest

import (
	"email-acceptor/models"
	"email-acceptor/pkg/errors"
	"email-acceptor/pkg/render"
	"encoding/json"
	"net/http"
)

func respond(wr http.ResponseWriter, status int, v interface{}) {
	if err := render.JSON(wr, status, v); err != nil {
		if loggable, ok := wr.(errorLogger); ok {
			loggable.Errorf("failed to write data to http ResponseWriter: %s", err)
		}
	}
}

func respondErr(wr http.ResponseWriter, err error) {
	resp := models.Failed{Error: err.Error()}
	respond(wr, http.StatusBadRequest, resp)
}

func readRequest(req *http.Request, v interface{}) error {
	if err := json.NewDecoder(req.Body).Decode(v); err != nil {
		return errors.Validation("Failed to read request body")
	}

	return nil
}

type errorLogger interface {
	Errorf(msg string, args ...interface{})
}
