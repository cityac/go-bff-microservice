//go:build dev

package debug

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func LockTimeout(ctx *gin.Context) {
	value, err := ctx.Cookie("debug-lock-timeout")
	if err != nil {
		log.Debug().Err(err).Msg("debug-lock-timeout cookie not found")
		return
	}

	timeout, err := strconv.Atoi(value)
	if err != nil {
		log.Warn().Err(err).Msg("Invalid debug-lock-timeout cookie value")
		return
	}

	if timeout > 0 {
		log.Warn().Msgf("Lock timeout set to %d milliseconds", timeout)
		time.Sleep(time.Duration(timeout) * time.Millisecond)
	}
}
