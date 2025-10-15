package controllers

import (
	"github.com/go-jose/go-jose/v3/jwt"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func ResolveController(c *gin.Context, url *url.URL) gin.HandlerFunc {
	host := c.Request.Host
	path := c.Request.URL.Path
	logger.Info().Msgf("host: %v, path: %v", host, path)

	return ProxyHandler(url)
}

func ProxyHandler(target *url.URL) gin.HandlerFunc {
	proxy := httputil.NewSingleHostReverseProxy(target)
	return func(c *gin.Context) {
		var claims struct {
			Email string `json:"email"`
			Name  string `json:"name"`
			jwt.Claims
		}

		claims.Name = "Jane Doe"
		userEmail := c.Request.Header.Get("X-User-Email")

		c.Request.URL.Host = target.Host
		c.Request.URL.Scheme = target.Scheme
		c.Request.Header.Set("X-Forwarded-Host", c.Request.Header.Get("Host"))

		c.Request.Header.Set("X-User-Name", claims.Name)
		c.Request.Header.Set("X-User-Email", userEmail)
		c.Request.Host = target.Host

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
