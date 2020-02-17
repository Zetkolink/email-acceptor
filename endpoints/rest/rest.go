package rest

import (
	"context"
	"email-acceptor/models"
	"email-acceptor/pkg/errors"
	"email-acceptor/pkg/logger"
	"email-acceptor/pkg/render"
	"email-acceptor/store"
	"github.com/gorilla/mux"
	"net/http"
)

type Rest struct {
	logger.Logger
}

// New initializes the server with routes exposing the given usecases.
func New(logger logger.Logger, ms *store.MessageStore) http.Handler {
	// setup router with default handlers
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedHandler)

	sh := http.StripPrefix("/docs/", http.FileServer(http.Dir("./dist/")))
	router.PathPrefix("/docs/").Handler(sh)

	// setup api endpoints
	addMessagesApi(router, ms, logger)

	return router
}

func notFoundHandler(wr http.ResponseWriter, req *http.Request) {
	_ = render.JSON(wr, http.StatusNotFound, errors.ResourceNotFound("path", req.URL.Path))
}

func methodNotAllowedHandler(wr http.ResponseWriter, req *http.Request) {
	_ = render.JSON(wr, http.StatusMethodNotAllowed, errors.ResourceNotFound("path", req.URL.Path))
}

type messages interface {
	GetMessageRequestAll(ctx context.Context, page int, perPage int) (string, int, int, int, int)
	GetMessageRequest(ctx context.Context, uniqueId string) (string, error)
	Send(msg models.MessageRequest) (string, error)
}
