package api

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/imrishuroy/legal-referral/util"
	"github.com/rs/zerolog/log"
)

func createKafkaProducer(config util.Config) (*kafka.Producer, error) {

	conf := kafka.ConfigMap{
		// User-specific properties that you must set
		"bootstrap.servers": config.BootStrapServers,
		"sasl.username":     config.SASLUsername,
		"sasl.password":     config.SASLPassword,

		// Fixed properties
		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "PLAIN",
		"acks":              "all"}

	p, err := kafka.NewProducer(&conf)
	if err != nil {
		log.Error().Msgf("Failed to create producer: %s", err)
		return nil, err
	}
	return p, nil
}

func (srv *Server) publishToKafka(topic string, key string, value string) {
	// create a new producer instance
	p, err := createKafkaProducer(srv.Config)
	if err != nil {
		log.Error().Msgf("Failed to create producer: %s", err)
		return
	}
	defer p.Close()

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

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          []byte(value),
	}, nil)
	if err != nil {
		log.Error().Msgf("Failed to produce message: %v", err)
		return
	}

	// send any outstanding or buffered messages to the Kafka broker and close the connection
	p.Flush(15 * 1000)
	p.Close()
}
