package utils

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
)

// Данные статистики запросов
type RequestStat struct {
	HandlerName string
	RequestTime time.Duration
}

// Функция для отправки статистики запросов в Kafka
func sendRequestStatToKafka(stat RequestStat, topic string, producer sarama.SyncProducer) {
	// Преобразуем RequestStat в массив байтов, чтобы отправить в Kafka
	statBytes := []byte(fmt.Sprintf("%s:%d", stat.HandlerName, stat.RequestTime.Milliseconds()))

	// Создаем сообщение Kafka
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(statBytes),
	}

	// Отправляем сообщение в Kafka
	_, _, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("Error sending request stat to Kafka: %v", err)
		return
	}

	log.Printf("Request stat sent to Kafka: %s", string(statBytes))
}

// Обертка для обработчиков HTTP, чтобы сохранять статистику запросов
func RequestStatMiddleware(next http.Handler, handlerName string, topic string, producer sarama.SyncProducer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		requestTime := time.Since(start)

		stat := RequestStat{HandlerName: handlerName, RequestTime: requestTime}
		go sendRequestStatToKafka(stat, topic, producer)
	})
}
