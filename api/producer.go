package api

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog/log"
)

func (server *Server) createProducer(userID string, postID string) {

	// creates a new producer instance

	conf := kafka.ConfigMap{
		// User-specific properties that you must set
		"bootstrap.servers": server.config.BootStrapServers,
		"sasl.username":     server.config.SASLUsername,
		"sasl.password":     server.config.SASLPassword,

		// Fixed properties
		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "PLAIN",
		"acks":              "all"}

	p, _ := kafka.NewProducer(&conf)
	topic := "publish-feed"

	// go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()

	// TODO: see if we can create the initializer function in the main.go file and
	// create a separate function for the producer with the topic, data as parameters

	// produces a sample message to the user-created topic
	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(userID),
		Value:          []byte(postID),
	}, nil)
	if err != nil {
		return
	}

	// send any outstanding or buffered messages to the Kafka broker and close the connection
	p.Flush(15 * 1000)
	p.Close()

}

func (server *Server) publishToKafka(userID string, postID string) {
	topic := "publish-feed"

	// go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go func() {
		for e := range server.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()

	// create a separate function for the producer with the topic, data as parameters

	// produces a sample message to the user-created topic
	err := server.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(userID),
		Value:          []byte(postID),
	}, nil)
	if err != nil {
		log.Error().Msgf("Failed to produce message: %v", err)

		// reconnect to the producer
		//server.producer.Close()
		//server.createProducer(userID, postID)

		return
	}

	// send any outstanding or buffered messages to the Kafka broker and close the connection
	server.producer.Flush(15 * 1000)
	server.producer.Close()
}
