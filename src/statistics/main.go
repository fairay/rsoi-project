package main

import (
	"context"
	"log"
	"os"
	"statistics/controllers"
	"statistics/utils"

	"fmt"
	"math/rand"
	"time"

	"github.com/Shopify/sarama"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type KafkaSettings struct {
	Consumer sarama.ConsumerGroup
}

func InitKafka() *KafkaSettings {
	kafkaBrokers := utils.Config.Kafka.Endpoints
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	config := sarama.NewConfig()
	config.Net.TLS.Enable = false
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	consumer, err := sarama.NewConsumerGroup(kafkaBrokers, "1", config)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}

	return &KafkaSettings{
		Consumer: consumer,
	}
}

var messages = make(chan string)

type TaskHandler struct{}

func (*TaskHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (*TaskHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (*TaskHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}

func main() {
	// createTopic()
	rand.Seed(time.Now().UnixNano())

	utils.InitConfig()
	utils.InitLogger()
	defer utils.CloseLogger()

	kafka := InitKafka()
	go func() {
		ctx := context.Background()
		handler := &TaskHandler{}
		for {
			err := kafka.Consumer.Consume(ctx, utils.Config.Kafka.Topics, handler)
			if err != nil {
				fmt.Println(err.Error())
				panic(err)
			}

			if ctx.Err() != nil {
				fmt.Println(ctx.Err().Error())
				panic(err)
			}
			fmt.Println("Next loop")
		}
	}()

	r := controllers.InitRouter()
	utils.Logger.Print("Server started")
	fmt.Printf("Server is running on http://localhost:%d\n", utils.Config.Port)
	code := controllers.RunRouter(r, utils.Config.Port)

	utils.Logger.Printf("Server ended with code %s", code)
}
