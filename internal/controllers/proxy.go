package controllers

import (
	"fmt"
	"net/url"
	"os"

	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"

	"github.com/gin-gonic/gin"
)

func GetProxyController(c *gin.Context, envKey string) gin.HandlerFunc {
	env := os.Getenv(envKey)
	targetURL, _ := url.Parse(env)
	logger.Info().Msgf("--- RECEIVED REQUEST --- %v", env)
	return ResolveController(c, targetURL)
}

func AllVendorsHIDDENRequestsHandler(c *gin.Context) {
	controller := GetProxyController(c, "PRICE_VP_SERVICE_HOST")
	if controller != nil {
		controller(c)
	} else {
		_, err := fmt.Fprintf(c.Writer, "Method not found")
		if err != nil {
			logger.Error().Err(err).Msg("Method not found")
			return
		}
	}
}
func AllClientsHIDDENRequestsHandler(c *gin.Context) {
	controller := GetProxyController(c, "PRICE_CP_SERVICE_HOST")
	if controller != nil {
		controller(c)
	} else {
		_, err := fmt.Fprintf(c.Writer, "Method not found")
		if err != nil {
			logger.Error().Err(err).Msg("Method not found")
			return
		}
	}
}

func AllHIDDENRequestsHandler(c *gin.Context) {
	env := os.Getenv("HIDDEN_SERVICE_HOST")
	targetURL, _ := url.Parse(env)
	controller := ResolveController(c, targetURL)
	logger.Info().Msgf("--- RECEIVED REQUEST --- %v", env)
	if controller != nil {
		controller(c)
	} else {
		_, err := fmt.Fprintf(c.Writer, "Method not found")
		if err != nil {
			logger.Error().Err(err).Msg("Method not found")
			return
		}
	}
}

func AllPartneHIDDENRequestHandler(c *gin.Context) {
	env := os.Getenv("PARTNEHIDDEN_SERVICE_HOST")
	targetURL, _ := url.Parse(env)
	controller := ResolveController(c, targetURL)
	logger.Info().Msgf("--- RECEIVED REQUEST --- %v", env)
	if controller != nil {
		controller(c)
	} else {
		_, err := fmt.Fprintf(c.Writer, "Method not found")
		if err != nil {
			logger.Error().Err(err).Msg("Method not found")
			return
		}
	}
}

func AllUserRequestsHandler(c *gin.Context) {
	env := os.Getenv("USER_SERVICE_URL")
	targetURL, _ := url.Parse(env)
	controller := ResolveController(c, targetURL)
	logger.Info().Msgf("--- RECEIVED REQUEST --- %v", env)
	if controller != nil {
		controller(c)
	} else {
		_, err := fmt.Fprintf(c.Writer, "Method not found")
		if err != nil {
			logger.Error().Err(err).Msg("Method not found")
			return
		}
	}
}
