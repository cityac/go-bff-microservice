package kafka

import "gitlab.HIDDEN.com/HIDDEN/sms-core/constants"

//	type Topic struct {
//		Request string
//	}
type kafkatopic struct {
	Partners string
	Price    string
	Reply    string
}
type settings struct {
	KafkaServer   string
	GroupId       string
	KafkaTopic    kafkatopic
	ReplyTopic    string
	RequestTopics kafkatopic
}

var Settings = settings{
	KafkaServer:   "localhost:9092",
	GroupId:       "bff",
	KafkaTopic:    kafkatopic{Partners: "HIDDEN.request.partners", Reply: "HIDDEN.reply.bff"},
	ReplyTopic:    "HIDDEN.reply.bff",
	RequestTopics: kafkatopic{Partners: "HIDDEN.request.partners"},
}

var RequestTopics = map[string]string{
	constants.Request: "HIDDEN.request",
	constants.Insert:  "HIDDEN.insert",
	constants.Update:  "HIDDEN.update",
	constants.Delete:  "HIDDEN.delete",
	constants.RCreate: "HIDDEN.relay.create",
	constants.RUpdate: "HIDDEN.relay.update",
	constants.RDelete: "HIDDEN.relay.delete",
	"Request.Margin":  "HIDDEN.margin.request",
	"Requests.Email":  "HIDDEN.email.requests",
}
