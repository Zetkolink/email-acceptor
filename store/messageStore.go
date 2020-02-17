package store

import (
	"context"
	"email-acceptor/helpers"
	"email-acceptor/models"
	"encoding/json"
	"github.com/Zetkolink/go-amqp-reconnect/rabbitmq"
	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/lithammer/shortuuid"
	"github.com/streadway/amqp"
	"strings"
)

type MessageStore struct {
	db helpers.DbConnection
	ex string
	ch *rabbitmq.Channel
}

func NewMessageStore(db helpers.DbConnection, ch *rabbitmq.Channel, ex string) *MessageStore {
	return &MessageStore{
		db: db,
		ch: ch,
		ex: ex,
	}
}

func (ms MessageStore) Migrate() {
	msg := models.Messages{}
	msg.Migrate(ms.db.DB)
}

func (ms MessageStore) Send(msg models.MessageRequest) (string, error) {
	id := shortuuid.New()
	msg.UniqueId = &id
	request, err := json.Marshal(msg)
	if err != nil {
		return id, err
	}
	err = ms.ch.Publish(ms.ex, "", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        request,
	})
	if err != nil {
		return id, err
	}

	resp := models.Success{Id: id}

	respJson, _ := json.Marshal(resp)

	return string(respJson), nil
}

func (ms MessageStore) GetMessageRequest(_ context.Context, uniqueId string) (string, error) {
	messages := models.Messages{}
	err := ms.db.Where("id = ?", uniqueId).Find(&messages).Error
	if err != nil {
		return "", err
	}

	return messages.Obj, nil
}

func (ms MessageStore) GetMessageRequestAll(_ context.Context, page int, perPage int) (string, int, int, int, int) {
	messages := make([]models.Messages, 0)

	pagin := pagination.Paging(&pagination.Param{
		DB:      ms.db.Where("id != ?", ""),
		Page:    page,
		Limit:   perPage,
		OrderBy: []string{"id desc"},
		ShowSQL: false,
	}, &messages)

	res := make([]string, 0, len(messages))

	for _, v := range messages {
		res = append(res, v.Obj)
	}

	return "[" + strings.Join(res, ",") + "]", pagin.TotalRecord, pagin.TotalPage, pagin.NextPage, pagin.PrevPage
}
