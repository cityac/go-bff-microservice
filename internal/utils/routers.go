package utils

import (
	"github.com/gin-gonic/gin"
	"os"
)

func FeatureFlag(flagName string, activeHandler, defaultHandler gin.HandlerFunc) gin.HandlerFunc {
	flagValue := os.Getenv(flagName)

	if flagValue == "true" {
		return activeHandler
	} else {
		return defaultHandler
	}
}
