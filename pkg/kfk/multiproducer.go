package kafka

import (
	"errors"
	"github.com/IBM/sarama"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
)

type KafkaMultiProducer struct {
	producers  map[string]*KafkaProducer
	brokerList []string
}

func NewKafkaMultiProducer(requestTopics map[string]string, brokerList []string) *KafkaMultiProducer {
	producers := make(map[string]*KafkaProducer)
	for name, topic := range requestTopics {
		logger.Info().Msgf("NAME + TOPIC, %v, %v", name, topic)
		producers[name] = NewKafkaProducer(brokerList, topic, 2)
	}
	return &KafkaMultiProducer{producers: producers, brokerList: brokerList}
}

func (k *KafkaMultiProducer) Produce(name string, payload []byte, key string, headers []sarama.RecordHeader) error {
	p := k.producers[name]
	if p == nil {
		err := errors.New("invalid producer id")
		logger.Error().Err(err).Msgf("invalid producer id %v", name)
		return err
	}
	return p.Produce(payload, key, headers)
}

func (k *KafkaMultiProducer) AddProducerSub(name string) {
	k.producers[name] = NewKafkaProducer(k.brokerList, RequestTopics[name], 2)
}

func (k *KafkaMultiProducer) Close() error {
	for _, producer := range k.producers {
		err := producer.Close()
		if err != nil {
			logger.Error().Err(err).Msg("close kafka producer failed")
			return err
		}
	}
	return nil
}
