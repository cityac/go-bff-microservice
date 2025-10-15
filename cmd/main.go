package main

import (
	"bff/env"
	"bff/internal/controllers"
	"bff/internal/middlewares"
	"bff/internal/routers"
	"bff/pkg/events"
	kafka "bff/pkg/kfk"
	"bff/pkg/redis"
	ws "bff/pkg/ws"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/auth"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/avro"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/constants"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/types"
)

const portNumber = ":8080"

func main() {
	config := env.NewEnvConfig()
	ofgaConfig := env.NewOFGAClientConfig()

	logger.InitLogger(config.LogLevel, os.Stdout)

	registry := avro.NewAvroRepo(config.SchemaRegistryUrl)
	fga := auth.NewOFGARepo(ofgaConfig)
	redisProvider := redis.NewRedisProvider()
	k := kafka.NewKafkaService(config.KafkaUrl, kafka.Settings.KafkaTopic.Reply)
	l := kafka.NewLogProducer(strings.Split(config.KafkaUrl, ","), 2, registry)
	controllers.NewRestController(registry, k, fga, &redisProvider, config)

	if config.PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	}
	router := gin.Default()
	pprof.Register(router)

	router.Use(middlewares.CORS(config))

	router = routers.NewRouter(router, l, fga, &redisProvider)

	server := http.Server{
		Addr:    portNumber,
		Handler: router,
	}

	socket := ws.Start()

	router.POST("/test/notify", func(c *gin.Context) {
		var event struct {
			EventName string `json:"event"`
			Message   string `json:"message"`
		}

		organizationID := c.Query("org_id")

		if organizationID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Please specify org_id in query"})
			return
		}

		err := c.ShouldBindJSON(&event)
		if err != nil {
			logger.Error().Msgf("Failed get the message: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		notification := types.Notification{
			Message:          event.Message,
			OrganizationId:   constants.OrganizationID,
			NotificationType: "balance",
		}

		payload, err := json.Marshal(notification)
		if err != nil {
			logger.Error().Msgf("Failed to encode the message: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		socket.EmitMessage(organizationID, event.EventName, payload)
		c.JSON(http.StatusOK, gin.H{"status": "message sent"})
	})

	events.StartEventService(config.KafkaUrl, socket, registry)
	// Create a channel to receive OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	// Start the server in a separate goroutine
	go func() {
		logger.Info().Msg("Starting server on port 8080...")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Msgf("Failed to start server: %v", err)
		}
	}()

	<-stop
	logger.Info().Msg("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal().Msgf("Server shutdown error: %v", err)
	}

	logger.Info().Msg("Server gracefully stopped")
}
