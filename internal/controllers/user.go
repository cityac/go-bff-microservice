package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v3/jwt"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/constants"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/customErrors"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/user"
	response "gitlab.HIDDEN.com/HIDDEN/sms-core/utils"
	"io"
	"net/http"
	"os"
	"strings"
)

func (rc *RestController) GetUser(c *gin.Context) {
	//id := c.Param("id")
	// Get the Authorization header from the request
	authHeader := c.GetHeader("Authorization")

	// Check if the Authorization header exists and starts with 'Bearer '
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		response.SendErrorResponse(c, errors.New("no token found in the Authorization header"), customErrors.RequestUnauthorized)
		return
	}

	// Extract the JWT token from the Authorization header
	token := strings.TrimPrefix(authHeader, "Bearer ")

	decode, err := jwt.ParseSigned(token)

	if err != nil {
		logger.Error().Err(err).Str("token", token).Msg("failed to parse token")
		response.SendErrorResponse(c, customErrors.ErrRequestUnauthorizedError)
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
			logger.Error().Err(err).Str("token", token).Msg("decode.UnsafeClaimsWithoutVerification")
			claims.Name = "Jane Doe"
			claims.Email = "janedoe@HIDDEN.com"
			claims.Subject = "4dfd9722-c6d4-4947-80c2-7fb9b8a8c5f6"
		}
	}

	logger.Warn().Msgf("CLAIMS %v", claims)

	req, err := http.NewRequest("GET", os.Getenv("USER_SERVICE_URL")+"/api/user/"+claims.UserId, nil)

	if err != nil {
		logger.Error().Err(err).Str("token", token).Msg("failed to create request")
		response.SendErrorResponse(c, err, customErrors.BadRequest)
		return
	}

	// Add a header to the request
	req.Header.Set(constants.OrgIdHeader, claims.OrgId)

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		logger.Error().Err(err).Str("token", token).Msg("failed to send request")
		response.SendErrorResponse(c, errors.New("cannot fetch user service url"), customErrors.BadRequest)
		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		response.SendErrorResponse(c, err)
		return
	}

	logger.Info().Msgf("RESPONSE %v", body)

	var user *sms_user.User

	err = json.Unmarshal(body, &user)

	if err != nil {
		response.SendErrorResponse(c, err)
		return
	}

	response.SendOkResponse(c, nil, user)
}
