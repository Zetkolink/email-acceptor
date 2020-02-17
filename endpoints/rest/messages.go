package rest

import (
	"email-acceptor/models"
	"email-acceptor/pkg/logger"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func addMessagesApi(router *mux.Router, ms messages, lg logger.Logger) {
	pc := &messagesController{}
	pc.pc = ms
	pc.Logger = lg

	router.HandleFunc("/notifs/", pc.getAll).Methods(http.MethodGet)
	router.HandleFunc("/notifs/", pc.post).Methods(http.MethodPost)
	router.HandleFunc("/notifs/{id}", pc.getById).Methods(http.MethodGet)
}

type messagesController struct {
	logger.Logger

	pc messages
}

func (cc messagesController) post(wr http.ResponseWriter, req *http.Request) {
	messageR := models.MessageRequest{}
	if err := readRequest(req, &messageR); err != nil {
		cc.Warnf("failed to read user request: %s", err)
		respond(wr, http.StatusBadRequest, "")
		return
	}

	err := messageR.Validate()
	if err != nil {
		respondErr(wr, err)
		return
	}
	id, err := cc.pc.Send(messageR)
	if err != nil {
		respond(wr, http.StatusInternalServerError, "")
		return
	}

	wr.WriteHeader(http.StatusOK)
	wr.Write([]byte(id))
}

func (cc messagesController) getAll(wr http.ResponseWriter, req *http.Request) {
	page := 1
	pageSize := 20
	p, ok := req.URL.Query()["page"]
	if ok {
		pq, err := strconv.Atoi(p[0])
		if err != nil {
			respond(wr, http.StatusBadRequest, "")
			return
		}
		page = pq
	}
	ps, ok := req.URL.Query()["per_page"]
	if ok {
		psq, err := strconv.Atoi(ps[0])
		if err != nil {
			respond(wr, http.StatusBadRequest, "")
			return
		}
		pageSize = psq
	}

	response, total, totalPage, next, prev := cc.pc.GetMessageRequestAll(req.Context(), page, pageSize)
	wr.Header().Set("X-Total", strconv.Itoa(total))
	wr.Header().Set("X-Total-Pages", strconv.Itoa(totalPage))
	wr.Header().Set("X-Per-Page", strconv.Itoa(pageSize))
	wr.Header().Set("X-Page", strconv.Itoa(page))
	wr.Header().Set("X-Next-Page", strconv.Itoa(next))
	wr.Header().Set("X-Prev-Page", strconv.Itoa(prev))
	wr.WriteHeader(http.StatusOK)
	wr.Write([]byte(response))
}

func (cc messagesController) getById(wr http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		respond(wr, http.StatusBadRequest, "")
		return
	}

	response, err := cc.pc.GetMessageRequest(req.Context(), id)
	if err != nil {
		respond(wr, http.StatusInternalServerError, "")
		return
	}
	wr.WriteHeader(http.StatusOK)
	wr.Write([]byte(response))
}
