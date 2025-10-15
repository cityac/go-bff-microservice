package kafka

import (
	"github.com/IBM/sarama"
	"github.com/oklog/ulid/v2"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/utils"
	"log"
	_ "strconv"
	"strings"
)

type KafkaService struct {
	producers *KafkaMultiProducer
	consumer  *KafkaConsumer
	waitList  map[string]chan *utils.RespMsg
}

func NewKafkaService(kafkaUrl string, reply string) (controller *KafkaService) {
	// If multiple services need to be used than add topic to Producer
	servers := strings.Split(kafkaUrl, ",")
	ret := &KafkaService{
		producers: NewKafkaMultiProducer(RequestTopics, servers),
		consumer:  NewKafkaConsumer(servers, reply, 5),
		waitList:  make(map[string]chan *utils.RespMsg),
	}
	go ret.KafkaConsumeReplies()
	return ret
}

//func (c *KafkaService) DoPost(payload []byte) (*Response, error) {
//	id := uuid.NewString()
//	comm := make(chan *types.ReplyMessagePayload)
//	c.waitList[id] = comm
//
//	defer func() {
//		delete(c.waitList, id)
//		close(comm)
//	}()
//
//	err := c.producer.Produce(payload, id, nil)
//	if err != nil {
//		return nil, err
//	}
//	select {
//	case reply := <-comm:
//		return &Response{StatusCode: reply.status, Body: reply.payload}, nil
//	case <-time.After(5 * time.Second):
//	}
//	return &Response{StatusCode: 408, Body: []byte("timeout")}, nil
//}

//func (c *KafkaService) kafkaConsumeReplies() {
//	comm, err := c.consumer.Consume()
//	if err != nil {
//		panic(err)
//	}
//	for {
//		reply := <-comm
//		replyChan, exists := c.waitList[string(reply.Key)]
//		if exists {
//			statusStr := string(reply.Headers[0].Value)
//			status, err := strconv.ParseUint(statusStr, 10, 32)
//			if err != nil {
//				status = 404
//			}
//			replyChan <- &RestMessage{id: string(reply.Key), payload: reply.Value, status: int(status)}
//			close(replyChan)
//			delete(c.waitList, string(reply.Key))
//		}
//	}
//}

type RespMsg struct {
	Value   []byte
	Headers []*sarama.RecordHeader
}

func (c *KafkaService) ProduceMessage(name string, message []byte, headers []sarama.RecordHeader) (chan *utils.RespMsg, error) {
	chann := make(chan *utils.RespMsg)
	requestId := ulid.Make().String()
	logger.Info().Msgf("REQUEST Id %v", requestId)
	c.waitList[requestId] = chann
	headers = append(headers,
		sarama.RecordHeader{Key: []byte("MessageId"), Value: []byte(requestId)},
	)

	err := c.producers.Produce(name, message, requestId, headers)

	return chann, err
}

func (c *KafkaService) KafkaConsumeReplies() {
	comm, err := c.consumer.Consume()
	if err != nil {
		log.Println("ERROR", err)
	}
	go func() {
		for {
			reply, ok := <-comm
			if !ok {
				logger.Info().Msg("Channel closed")
				break // Exit the loop if the channel is closed
			}
			respMessage := utils.RespMsg{
				Value:   reply.Value,
				Headers: reply.Headers,
			}
			logger.Info().Msgf("MESSAGE TO CHANNEL %v", string(reply.Key))
			chann := c.waitList[string(reply.Key)]
			if chann != nil {
				chann <- &respMessage
				delete(c.waitList, string(reply.Key))
				close(chann)
			}
		}
	}()
}
