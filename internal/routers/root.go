package routers

import (
	"bff/internal/middlewares"
	kafka "bff/pkg/kfk"
	"bff/pkg/redis"

	"github.com/gin-gonic/gin"

	"gitlab.HIDDEN.com/HIDDEN/sms-core/auth"
)

type BffRouter struct {
	router *gin.Engine
	log    *kafka.LogProducer
	fga    *auth.OFGARepo
	redis  *redis.RedisProvider
}

var Router *BffRouter

func NewRouter(router *gin.Engine, lg *kafka.LogProducer, fga *auth.OFGARepo, redis *redis.RedisProvider) *gin.Engine {
	Router = &BffRouter{router: router, log: lg, fga: fga, redis: redis}
	// router.Use(middlewares.VerifyToken())
	Router.MountHealthzRouter()

	router.Use(middlewares.Auth())

	Router.MountAccountsRouter()
	Router.MountRefBookRoutes()
	Router.MountPartnerRouter()
	Router.DictionariesRouter()
	Router.MountUploadRoutes()
	Router.MountRulesRoutes()
	Router.MountPricesRoutes()
	Router.MountHIDDENRoutes()
	Router.MountRatePlanRouter()
	Router.MountPriceListRouter()
	Router.MountBalanceRouter()
	Router.MountUserRouter()
	Router.MountOrgRouter()
	Router.MountTestRoutes()
	Router.MountEmailRouter()
	Router.MountLocksRouter()
	Router.MountFeatureRouter()

	return router
}
