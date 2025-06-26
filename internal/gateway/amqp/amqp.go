package amqp

import (
	"context"
	"encoding/json"
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

func (a *AMQP) Start(ctx context.Context, connection string) error {

	// 1. Подключение к RabbitMQ
	conn, err := amqp.Dial(connection)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// 2. Создание канала
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to create channel: %v", err)
	}
	defer ch.Close()

	// 3. Объявление очереди
	_, err = ch.QueueDeclare(
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
	msgs, err := ch.Consume(
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
		a.processMessage(ctx, ch, msg)
	}

	return err
}

func (a *AMQP) processMessage(ctx context.Context, ch *amqp.Channel, msg amqp.Delivery) {
	// 1. Парсим запрос
	reportID := string(msg.Body)
	log.Printf("📥 Received request for report: %s", reportID)

	// 2. Формируем данные (ваша бизнес-логика)
	data, err := a.Service.ReportExecute(ctx, reportID)
	// data.ID = reportID
	if err != nil {
		log.Print(data)
		return
	}

	// 3. Сериализуем в JSON
	response, err := json.Marshal(data)
	if err != nil {
		log.Printf("❌ Failed to marshal data: %v", err)
		return
	}

	// 4. Отправляем ответ
	err = ch.Publish(
		"",          // exchange
		msg.ReplyTo, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          response,
			CorrelationId: msg.CorrelationId,
		},
	)
	if err != nil {
		log.Printf("❌ Failed to send response: %v", err)
	} else {
		log.Printf("📤 Sent response for report %s", reportID)
	}
}
