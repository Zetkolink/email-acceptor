package helpers

import (
	"github.com/Zetkolink/go-amqp-reconnect/rabbitmq"
	"github.com/streadway/amqp"
	"log"
)

type Rabbit struct {
	*rabbitmq.Connection
	amqpUrl string
}

func InitRabbit(cfg Config) Rabbit {
	rb := Rabbit{}
	rb.amqpUrl = cfg.Rb.Amqp
	rb.connect()

	return rb
}

func (r *Rabbit) connect() {
	conn, err := rabbitmq.Dial(r.amqpUrl)
	if err != nil {
		log.Fatal(err)
	}

	r.Connection = conn
}

func (r *Rabbit) InitChannel(cfg Config) *rabbitmq.Channel {
	exchange := &rabbitmq.Exchange{
		ExchangeName: cfg.Rb.Exchange,
		Kind:         amqp.ExchangeDirect,
		Durable:      true,
		AutoDelete:   false,
		Internal:     false,
		NoWait:       false,
		Args:         nil,
	}

	queue := &rabbitmq.Queue{
		QueueName:  cfg.Rb.Queue,
		Durable:    false,
		AutoDelete: false,
		Exclusive:  false,
		Args:       nil,
	}

	bind := &rabbitmq.QueueBind{
		Key:  "",
		Args: nil,
	}

	sendCh, err := r.Channel(exchange, queue, bind)
	if err != nil {
		log.Panic(err)
	}

	return sendCh
}
