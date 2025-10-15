package kafka

import (
	"github.com/IBM/sarama"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	"log"
)

type KafkaConsumer struct {
	consumer sarama.Consumer
	topic    string
}

func NewKafkaConsumer(brokerList []string, topic string, maxRetry int) *KafkaConsumer {
	consumer, err := sarama.NewConsumer(brokerList, nil)
	if err != nil {
		logger.Panic().Err(err).Msg("Error creating new kafka consumer")
	}
	return &KafkaConsumer{consumer: consumer, topic: topic}
}

func (k *KafkaConsumer) Consume() (chan *sarama.ConsumerMessage, error) {
	partitionList, err := k.consumer.Partitions(k.topic) //get all partitions
	if err != nil {
		logger.Error().Err(err).Msg("k.consumer.Partitions")
		return nil, err
	}
	messages := make(chan *sarama.ConsumerMessage, 256)
	initialOffset := sarama.OffsetNewest //offset to start reading message from
	for _, partition := range partitionList {
		pc, err := k.consumer.ConsumePartition(k.topic, partition, initialOffset)
		if err != nil {
			logger.Error().Err(err).Msgf("Failed to start consumer for partition %d: %v", partition, err)
			continue
		}
		go func(pc sarama.PartitionConsumer) {
			defer pc.Close()
			for message := range pc.Messages() {
				log.Println("NEW MESSAGE", string(message.Key), message.Partition)
				logger.Info().Msgf("NEW MESSAGE %v, %v", string(message.Key), message.Partition)
				select {
				case messages <- message:
					// Successfully sent the message
				default:
					// Channel is full, handle this case (e.g., log, drop message, etc.)
					logger.Warn().Msgf("WARN: message channel is full, dropping message %v, %v, %v", message.Headers, message.Topic, message.Offset)
				}
			}
		}(pc)
	}
	return messages, nil
}

func (k *KafkaConsumer) PauseAll() {
	k.consumer.PauseAll()
}

func (k *KafkaConsumer) ResumeAll() {
	k.consumer.ResumeAll()
}

func (k *KafkaConsumer) Close() error {
	return k.consumer.Close()
}
