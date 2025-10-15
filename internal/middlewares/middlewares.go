package middlewares

import (
	"bff/env"
	"errors"
	"github.com/go-jose/go-jose/v3/jwt"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/constants"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/customErrors"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS(config *env.EnvConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		origin := c.GetHeader("Origin")
		allowedOrigins := []string{
			config.DashboardHost,
			config.PartnerClientHost,
			config.PriceClientHost,
			config.HIDDENClientHost,
			config.HIDDENOauthClientHost,
			config.HIDDENAdminHost,
			config.OrgAdminHost,
		}

		if !config.PRODUCTION {
			allowedOriginsForDev := []string{
				"http://localhost:8000",
				"http://localhost:8001",
				"http://localhost:8002",
				"http://localhost:8003",
				"http://localhost:8004",
				"http://localhost:8005",
			}
			allowedOrigins = append(allowedOrigins, allowedOriginsForDev...)
		}

		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin != "" && allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Ids")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// Check if the Authorization header exists and starts with 'Bearer '
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.SendErrorResponse(c, errors.New("no token found in the Authorization header"), customErrors.RequestUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		decode, err := jwt.ParseSigned(token)

		if err != nil {
			utils.SendErrorResponse(c, customErrors.ErrRequestUnauthorizedError)
			return
		}
		var claims struct {
			Email    string `json:"email"`
			Name     string `json:"name"`
			UserId   string `json:"user_id"`
			OrgId    string `json:"org_id"`
			UserRole string `json:"user_role"`
			// You can include standard claims here too
			jwt.Claims
		}

		// Extract claims without validating the signature
		if decode != nil {
			err := decode.UnsafeClaimsWithoutVerification(&claims)
			if err != nil {
				logger.Error().Err(err).Msg("decode.UnsafeClaimsWithoutVerification")
				claims.Name = "Jane Doe"
				claims.Email = "janedoe@HIDDEN.com"
				claims.Subject = "4dfd9722-c6d4-4947-80c2-7fb9b8a8c5f6"
			}
		}
		c.Set(constants.OrgId, claims.OrgId)
		c.Set(constants.UserId, claims.UserId)
		c.Request.Header.Set(constants.OrgIdHeader, claims.OrgId)
		c.Request.Header.Set(constants.UserIdHeader, claims.UserId)
		c.Request.Header.Set("x-user-email", claims.Email)
		c.Next()
	}
}
