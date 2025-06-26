package amqp

import (
	"context"
	"log"

	"github.com/imirjar/rb-diver/internal/models"
	"github.com/streadway/amqp"
)

type Service interface {
	ReportExecute(context.Context, string) (models.Data, error)
}

type AMQP struct {
	Conn    *amqp.Connection
	Chan    *amqp.Channel
	Service Service
}

func New() *AMQP {
	return &AMQP{}
}

func (a *AMQP) Connect(ctx context.Context, connection string) error {
	conn, err := amqp.Dial(connection)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
		return err
	}
	a.Conn = conn

	ch, err := a.Conn.Channel()
	if err != nil {
		log.Fatalf("Failed to create channel: %v", err)
	}

	a.Chan = ch

	return nil
}

func (a *AMQP) Disconnect() error {
	defer a.Conn.Close()
	return a.Chan.Close()
}

func (a *AMQP) Start(ctx context.Context, connection string) error {

	// 3. Объявление очереди
	_, err := a.Chan.QueueDeclare(
		"data_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// 4. Подписка на сообщения
	msgs, err := a.Chan.Consume(
		"data_queue", // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		log.Fatalf("Failed to consume: %v", err)
	}

	log.Println("🚀 Consumer started. Waiting for messages...")

	// 5. Обработка сообщений
	for msg := range msgs {
		a.processMessage(ctx, a.Chan, msg)
	}

	return err
}
