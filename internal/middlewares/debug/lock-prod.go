//go:build !dev

package debug

import (
	"github.com/gin-gonic/gin"
)

func LockTimeout(_ *gin.Context) {
	// In production, this middleware does nothing.
	// It is used to prevent accidental lock timeouts in production environments.
	// The actual logic is handled in the dev build.
	// This is a placeholder to ensure the middleware exists without any functionality.
}
