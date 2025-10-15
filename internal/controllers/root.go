package controllers

import (
	"bff/env"
	kafka2 "bff/pkg/kfk"
	"bff/pkg/redis"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/auth"

	"gitlab.HIDDEN.com/HIDDEN/sms-core/avro"
)

var Controller *RestController

type RestController struct {
	config  *env.EnvConfig
	adapter *kafka2.KafkaService
	avro    *avro.AvroRepo
	fga     *auth.OFGARepo
	redis   *redis.RedisProvider
}

func NewRestController(registry *avro.AvroRepo, kafka *kafka2.KafkaService, fga *auth.OFGARepo, redis *redis.RedisProvider, config *env.EnvConfig) {
	ret := &RestController{
		config:  config,
		adapter: kafka,
		avro:    registry,
		fga:     fga,
		redis:   redis,
	}
	Controller = ret
}
