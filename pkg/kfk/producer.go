package kafka

import (
	"github.com/IBM/sarama"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaProducer(brokerList []string, topic string, maxRetry int) *KafkaProducer {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = maxRetry
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		logger.Panic().Err(err).Msg("sarama.NewSyncProducer")
	}
	return &KafkaProducer{producer: producer, topic: topic}
}

func (k *KafkaProducer) Produce(payload []byte, key string, headers []sarama.RecordHeader) error {
	var saramaHeaders []sarama.RecordHeader = nil
	if headers != nil && len(headers) > 0 {
		for _, header := range headers {
			saramaHeaders = append(
				saramaHeaders,
				sarama.RecordHeader{
					Key:   header.Key,
					Value: header.Value})
		}
	}
	msg := &sarama.ProducerMessage{
		Topic:   k.topic,
		Value:   sarama.ByteEncoder(payload),
		Key:     sarama.StringEncoder(key),
		Headers: ([]sarama.RecordHeader)(saramaHeaders),
	}
	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		logger.Error().Err(err).Msgf("k.producer.SendMessage")
		return err
	}
	logger.Info().Msgf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", k.topic, partition, offset)
	return nil
}

func (k *KafkaProducer) Close() error {
	return k.producer.Close()
}
