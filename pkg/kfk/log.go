package kafka

import (
	"github.com/IBM/sarama"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/actionlog"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/avro"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/constants"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/utils"
)

type LogProducer struct {
	producer sarama.SyncProducer
	avro     *avro.AvroRepo
}

func NewLogProducer(brokerList []string, maxRetry int, avro *avro.AvroRepo) *LogProducer {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = maxRetry
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		logger.Panic().Err(err).Msg("sarama.NewSyncProducer")
	}
	return &LogProducer{producer: producer, avro: avro}
}

func (k *LogProducer) Produce(payload []byte, key string, headers []sarama.RecordHeader) error {
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
		Topic:   "HIDDEN.us",
		Value:   sarama.ByteEncoder(payload),
		Key:     sarama.StringEncoder(key),
		Headers: ([]sarama.RecordHeader)(saramaHeaders),
	}
	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		logger.Error().Err(err).Msg(" k.producer.SendMessage")
		return err
	}
	logger.Info().Msgf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", "HIDDEN.us", partition, offset)
	return nil
}

func (k *LogProducer) LogEvent(data any, user actionlog.User, action string, entity string, orgId string) {

	bytes, err := utils.SerializeData(data)

	if err != nil {
		return
	}

	logEvent := actionlog.LogEvent{
		User:   user,
		Data:   bytes,
		Action: action,
	}

	encode, err := k.avro.Encode(logEvent)
	if err != nil {
		return
	}

	rHeaders := []sarama.RecordHeader{
		{Key: []byte("Entity"), Value: []byte(entity)},
		{Key: []byte("OrganizationId"), Value: []byte(constants.OrganizationID)},
	}

	err = k.Produce(encode, "", rHeaders)
	if err != nil {
		logger.Error().Err(err).Msg(" k.Produce")
	}
}
