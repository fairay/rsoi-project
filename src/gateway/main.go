package main

import (
	"crypto/tls"
	"gateway/controllers"
	"gateway/utils"
	"log"

	"fmt"
	"math/rand"
	"time"

	"github.com/Shopify/sarama"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func createTopic() {
	// Настройка Kafka
	kafkaBrokers := []string{"kafka:29092"}
   
	// Настройки Kafka админ клиента
	config := sarama.NewConfig()
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = &tls.Config{
		InsecureSkipVerify: true,
	}
	admin, err := sarama.NewClusterAdmin(kafkaBrokers, config)
	if err != nil {
		log.Fatalf("Failed to create a Kafka admin: %v", err)
	}
	defer admin.Close()
   
	// Создание топика
	topicName := "test-topic"
	numPartitions := 3
	replicationFactor := int16(1)
   
	err = admin.CreateTopic(topicName, &sarama.TopicDetail{
	 NumPartitions:     3,
	 ReplicationFactor: replicationFactor,
	}, false)
   
	if err != nil {
	 log.Fatalf("Failed to create Kafka topic: %v", err)
	}
   
	fmt.Printf("Kafka topic '%s' with %d partitions created successfully\n", topicName, numPartitions)
}

func main() {
	// createTopic()
	rand.Seed(time.Now().UnixNano())

	utils.InitConfig()
	utils.InitLogger()
	defer utils.CloseLogger()

	r := controllers.InitRouter()
	utils.Logger.Print("Server started")
	fmt.Printf("Server is running on http://localhost:%d\n", utils.Config.Port)
	code := controllers.RunRouter(r, utils.Config.Port)

	utils.Logger.Printf("Server ended with code %s", code)
}
