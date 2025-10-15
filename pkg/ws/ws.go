package ws

import (
	"github.com/Baiguoshuai1/shadiaosocketio"
	"github.com/Baiguoshuai1/shadiaosocketio/websocket"
	"github.com/go-jose/go-jose/v3/jwt"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	"net/http"
)

type Message struct {
	Id      int    `json:"id"`
	Channel string `json:"channel"`
}

type Desc struct {
	Text string `json:"text"`
}

type Server struct {
	server *shadiaosocketio.Server
}

const (
	addr   = ":2233"
	socket = "/socket.io/"
)

func Start() *Server {
	server := shadiaosocketio.NewServer(*websocket.GetDefaultWebsocketTransport())

	err := server.On(shadiaosocketio.OnConnection, func(c *shadiaosocketio.Channel) {})
	if err != nil {
		logger.Error().Err(err).Msg("server.On")
	}

	err = server.On(shadiaosocketio.OnDisconnection, func(c *shadiaosocketio.Channel, reason websocket.CloseError) {
		logger.Info().Msgf("sIO Channel closed for client: %v", c.Id())
	})
	if err != nil {
		logger.Error().Err(err).Msg("server.On")
	}

	err = server.On("login", func(c *shadiaosocketio.Channel, token string) {
		decode, err := jwt.ParseSigned(token)
		if err != nil {
			logger.Error().Err(err).Msg("jwt.ParseSigned")
			c.Close()
		}

		var claims struct {
			Email    string `json:"email"`
			Name     string `json:"name"`
			UserId   string `json:"user_id"`
			OrgId    string `json:"org_id"`
			UserRole string `json:"user_role"`
			jwt.Claims
		}

		// Extract claims without validating the signature
		if decode != nil {
			if err := decode.UnsafeClaimsWithoutVerification(&claims); err != nil {
				logger.Error().Err(err).Msg("decode.UnsafeClaimsWithoutVerification")
				c.Close()
			}
		}

		if claims.OrgId == "" {
			logger.Info().Msgf("No org_id is present in the token, closing the channel %v", c.Id())
			c.Close()
			return
		}

		err = c.Join(claims.OrgId)
		if err != nil {
			logger.Error().Err(err).Msgf("Client failed to connect to organization room: orgID: %s", claims.OrgId)
			c.Close()
			return
		}

		err = c.Join(claims.UserId)
		if err != nil {
			logger.Error().Err(err).Msgf("Client failed to connect to user room: orgID: %s, userId: %s", claims.OrgId, claims.UserId)
			c.Close()
			return
		}

		logger.Info().Msgf("Client successfully authorized, id:'%s', org_id:'%s'", c.Id(), claims.OrgId)
	})
	if err != nil {
		logger.Error().Err(err).Msg("server.On")
	}

	err = server.On("/ackFromClient", func(c *shadiaosocketio.Channel, msg Message, num int) (int, Desc, string) {
		// Your existing ack handling logic
		return 1, Desc{}, ""
	})
	if err != nil {
		logger.Error().Err(err).Msg("server.On")
	}

	serveMux := http.NewServeMux()
	serveMux.Handle(socket, server)

	logger.Info().Msg("[ws] Starting server on :2233...")

	go func() {
		if err := http.ListenAndServe(addr, serveMux); err != nil {
			logger.Error().Err(err).Msg("Failed to start server:")
		}
	}()

	return &Server{
		server: server,
	}
}

func (s *Server) EmitMessage(room, event string, message interface{}) {
	if s.server != nil {
		s.server.BroadcastTo(room, event, message)
	}
}
