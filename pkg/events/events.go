package events

import (
	kafka "bff/pkg/kfk"
	"bff/pkg/ws"
	"github.com/IBM/sarama"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/avro"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/constants"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/types"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/utils"
	"log"
	"os"
	"os/signal"
	_ "strconv"
	"strings"
	"syscall"
)

type EventService struct {
	consumer *kafka.KafkaConsumer
	socket   *ws.Server
	avro     *avro.AvroRepo
}

var (
	topic            = "HIDDEN.notifications"
	defaultEventName = "notifications"
)

func StartEventService(kafkaUrl string, socket *ws.Server, avro *avro.AvroRepo) {
	// If multiple services need to be used than add topic to Producer
	servers := strings.Split(kafkaUrl, ",")
	ret := &EventService{
		consumer: kafka.NewKafkaConsumer(servers, topic, 5),
		socket:   socket,
		avro:     avro,
	}
	go ret.ConsumeEvents()
}

type RespMsg struct {
	Value   []byte
	Headers []*sarama.RecordHeader
}

func handleNotification[T types.EmitableNotification](c *EventService, value []byte, eventName string) {
	var notification T
	err := c.avro.Decode(value, &notification)
	if err != nil {
		logger.Error().Err(err).Msg("c.avro.Decode")
		return
	}
	log.Println("Transmitting notification", notification)
	c.socket.EmitMessage(notification.GetRoomName(), eventName, notification)
}

func (c *EventService) ConsumeEvents() {
	messages, err := c.consumer.Consume()
	if err != nil {
		logger.Error().Err(err).Msg("c.consumer.Consume")
		return
	}
	logger.Info().Msg("Started the consumer for notifications")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	for {
		select {
		case msg, ok := <-messages:
			if !ok {
				logger.Info().Msg("Message channel closed.")
				return
			}

			logger.Info().Msgf("Received message: Topic: %s, Partition: %d, Offset: %d \n",
				msg.Topic, msg.Partition, msg.Offset)

			headers := utils.ParseEventHeaders(msg.Headers)
			HIDDEN headers.EventType {
			case constants.CreditLimit:
				handleNotification[types.CreditLimitNotification](c, msg.Value, constants.CreditLimit)
			case constants.LogoutCommand:
				handleNotification[types.LogoutUserCommand](c, msg.Value, constants.LogoutCommand)
			case constants.PartnerVerified:
				handleNotification[types.PartnerVerifiedNotification](c, msg.Value, constants.PartnerVerified)
			default:
				handleNotification[types.Notification](c, msg.Value, defaultEventName)
			}

		case <-signals:
			log.Println("Interrupt signal received. Exiting...")
			return
		}
	}
}
