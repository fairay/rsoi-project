package models

import (
	"log"
	"net/http"
	"os"

	"github.com/Shopify/sarama"
)

type KafkaSettings struct {
	KafkaTopic	string
	Producer	sarama.SyncProducer
}

type Models struct {
	Client     *http.Client
	Flights    *FlightsM
	Privileges *PrivilegesM
	Tickets    *TicketsM

	Kafka	*KafkaSettings
}

func InitKafka() *KafkaSettings {
	kafkaBrokers := []string{"kafka:29092"}
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	config := sarama.NewConfig()
	config.Net.TLS.Enable = false
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(kafkaBrokers, config)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}

	return &KafkaSettings{
		KafkaTopic: "quickstart2",
		Producer: producer,
	}
}

func InitModels() *Models {
	models := new(Models)
	client := &http.Client{}

	models.Client = client
	models.Flights = NewFlightsM(client)
	models.Privileges = NewPrivilegesM(client)
	models.Tickets = NewTicketsM(client, models.Flights)
	models.Kafka = InitKafka()

	return models
}
