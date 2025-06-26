package amqp

import (
	"context"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

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
