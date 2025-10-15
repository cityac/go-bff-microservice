package controllers

import (
	"bff/pkg/domain/model"
	kafka "bff/pkg/kfk"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/partner"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/constants"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/utils"
)

type MigrateMRsCountPayload struct {
	ToDaily bool `json:"toDaily"`
}

type MigrateMRsCountByIdsPayload struct {
	RuleIDs []int `json:"ruleIds"`
}

type MigrateRatePlanPayload struct {
	RatePlanType string `json:"ratePlanType"`
	Tenant       string `json:"tenant"`
	OrgId        string `json:"orgId"`
}

type MigrateIzonetPartnersPayload struct {
	Tenant string `json:"tenant"`
	OrgId  string `json:"orgId"`
}

func (rc *RestController) MigrateMRsCount(c *gin.Context) {
	logger.Info().Msg("MigrateMRsCount")

	var request MigrateMRsCountPayload
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error().Err(err).Msg("Failed to bind request body")
		utils.SendErrorResponse(c, err)
		return
	}
	logger.Info().Interface("request", request).Msg("MigrateMRsCount successfully parsed request body")

	headers := utils.FormKafkaHeaders(c, constants.Rules, constants.ActionMigrate, kafka.Settings.ReplyTopic)

	kfkMsg, err := json.Marshal(request)
	if err != nil {
		logger.Error().Err(err).Msg("MigrateMRsCount json.Marshal")
		utils.SendErrorResponse(c, err)
		return
	}

	chann, err := rc.adapter.ProduceMessage(constants.Request, kfkMsg, headers)

	if err != nil {
		logger.Error().Err(err).Msg("MigrateMRsCount rc.adapter.ProduceMessage")
		utils.SendErrorResponse(c, err)
		return
	}

	d, _, err := utils.ListenReply[[]string](chann, rc.avro, time.Minute*10)

	if err != nil {
		logger.Error().Err(err).Msg("MigrateMRsCount utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkResponse(c, nil, d)
}

func (rc *RestController) MigrateMRsCountByIds(c *gin.Context) {
	logger.Info().Msg("MigrateMRsCountByIds")

	var request MigrateMRsCountByIdsPayload
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error().Err(err).Msg("Failed to bind request body")
		utils.SendErrorResponse(c, err)
		return
	}
	logger.Info().Interface("request", request).Msg("MigrateMRsCountByIds successfully parsed request body")

	headers := utils.FormKafkaHeaders(c, constants.Rules, constants.ActionMigrateById, kafka.Settings.ReplyTopic)

	kfkMsg, err := json.Marshal(request)
	if err != nil {
		logger.Error().Err(err).Msg("MigrateMRsCountByIds json.Marshal")
		utils.SendErrorResponse(c, err)
		return
	}

	chann, err := rc.adapter.ProduceMessage(constants.Request, kfkMsg, headers)

	if err != nil {
		logger.Error().Err(err).Msg("MigrateMRsCountByIds rc.adapter.ProduceMessage")
		utils.SendErrorResponse(c, err)
		return
	}

	d, _, err := utils.ListenReply[[]string](chann, rc.avro, time.Minute*10)

	if err != nil {
		logger.Error().Err(err).Msg("MigrateMRsCountByIds utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkResponse(c, nil, d)
}

func (rc *RestController) MigrateIzonetRatePlansFromCGRates(c *gin.Context) {
	var payload MigrateRatePlanPayload

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v\n", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Error().Err(err).Msg("Failed to bind request body")
		utils.SendErrorResponse(c, err)
		return
	}

	requestURL := fmt.Sprintf("%s/api/migration/ratePlans/%s", rc.config.HIDDENProviderHost, payload.RatePlanType)
	res, headers, err := utils.HttpPost(c, requestURL, payload)

	if err != nil {
		logger.Error().Err(err).Msg("MigrateIzonetRatePlansFromCGRates utils.HttpGet")
		utils.SendErrorResponse(c, err)
		return
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		logger.Error().Err(err).Msg("MigrateIzonetRatePlansFromCGRates utils.HttpGet")
		utils.SendErrorResponse(c, err)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error().Err(err).Msg("Body.Close")
		}
	}(res.Body)

	// uncomment for migrations to price-cp-service
	// c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	//priceCPCOntroller := GetProxyController(c, "PRICE_CP_SERVICE_HOST")
	//priceCPCOntroller(c)

	responseBody, err := io.ReadAll(res.Body)
	var response []model.MigrationResult

	err = json.Unmarshal(responseBody, &response)

	utils.SendOkResponse(c, headers, response)
}

func (rc *RestController) MigrateGlobalProductsLinkageFromCGRates(c *gin.Context) {
	var payload MigrateRatePlanPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Error().Err(err).Msg("Failed to bind request body")
		utils.SendErrorResponse(c, err)
		return
	}

	requestURL := fmt.Sprintf("%s/api/migration/ratePlans/%s/linkage", rc.config.HIDDENProviderHost, payload.RatePlanType)
	res, headers, err := utils.HttpPost(c, requestURL, payload)

	if err != nil {
		logger.Error().Err(err).Msg("MigrateGlobalProductsLinkageFromCGRates utils.HttpGet")
		utils.SendErrorResponse(c, err)
		return
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		logger.Error().Err(err).Msg("MigrateGlobalProductsLinkageFromCGRates utils.HttpGet")
		utils.SendErrorResponse(c, err)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error().Err(err).Msg("Body.Close")
		}
	}(res.Body)

	responseBody, err := io.ReadAll(res.Body)
	var response []model.MigrationResult

	err = json.Unmarshal(responseBody, &response)

	utils.SendOkResponse(c, headers, response)
}

func (rc *RestController) MigrateIzonetPartners(c *gin.Context) {
	var payload MigrateIzonetPartnersPayload

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v\n", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Error().Err(err).Msg("Failed to bind request body")
		utils.SendErrorResponse(c, err)
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	requestURL := fmt.Sprintf("%s/api/migration/partners", rc.config.HIDDENProviderHost)
	res, headers, err := utils.HttpPost(c, requestURL, payload)

	if err != nil {
		logger.Error().Err(err).Msg("MigrateIzonetPartners utils.HttpPost")
		utils.SendErrorResponse(c, err)
		return
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		logger.Error().Err(err).Msg("MigrateIzonetPartners utils.HttpPost")
		utils.SendErrorResponse(c, err)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error().Err(err).Msg("Body.Close")
		}
	}(res.Body)

	responseBody, err := io.ReadAll(res.Body)

	var partners []partner.Partner

	err = json.Unmarshal(responseBody, &partners)

	type scriptMigrationReply struct {
		Partners []partner.Partner
		Total    int
	}

	utils.SendOkResponse(c, headers, scriptMigrationReply{
		Partners: partners,
		Total:    len(partners),
	})
}
