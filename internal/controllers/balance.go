package controllers

import (
	kafka "bff/pkg/kfk"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/constants"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/customErrors"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	partnercore "gitlab.HIDDEN.com/HIDDEN/sms-core/partner"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/types"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/utils"
	"time"
)

func (rc *RestController) GetPartnerBalance(c *gin.Context) {
	partnerId := c.Param("partnerId")

	value := partnerId

	headers := utils.FormKafkaHeaders(c, constants.Balance, constants.ActionGetBy, kafka.Settings.ReplyTopic)

	kfkMsg, err := json.Marshal(value)
	if err != nil {
		logger.Error().Err(err).Msg("json.Marshal")
		utils.SendErrorResponse(c, err)
		return
	}

	chann, err := rc.adapter.ProduceMessage(constants.Request, kfkMsg, headers)

	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		utils.SendErrorResponse(c, err)
		return
	}

	d, _, err := utils.ListenReply[partnercore.Balance](chann, rc.avro, time.Second*5)

	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkResponse(c, nil, d)
}

func (rc *RestController) CreatePartnerBalance(c *gin.Context) {
	var balancePayload partnercore.CreateBalancePayload
	var req partnercore.Balance

	if err := c.ShouldBindJSON(&balancePayload); err != nil {
		logger.Error().Err(err).Msg("ShouldBindJSON")
		utils.SendErrorResponse(c, customErrors.ErrRequestBadRequestError)
		return
	}

	// Validate payload
	err := balancePayload.Validate()

	if err != nil {
		logger.Error().Err(err).Msg("balancePayload.Validate")
		utils.SendErrorResponse(c, customErrors.ErrRequestUnprocessableEntityError)
		return
	}

	// Verify if partner with CreateBalancePayload.PartnerId exists
	err = rc.ValidatePartnerID(c, balancePayload.PartnerId)

	if err != nil {
		logger.Error().Err(err).Msg("rc.ValidatePartnerID")
		utils.SendErrorResponse(c, customErrors.ErrRequestUnprocessableEntityError)
		return
	}

	// Prepare kafka payload
	req.PartnerId = balancePayload.PartnerId
	req.CreditLimit = balancePayload.CreditLimit
	req.Comment = balancePayload.Comment
	req.SendRateChangeNotifications = balancePayload.SendRateChangeNotifications
	req.HIDDENPeriods = balancePayload.HIDDENPeriods

	headers := utils.FormKafkaHeaders(c, constants.Balance, constants.ActionCreate, kafka.Settings.ReplyTopic)

	kfkMsg, err := rc.avro.Encode(req)

	logger.Info().Msgf("kfkMsg: %v", kfkMsg)

	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		utils.SendErrorResponse(c, err)
		return
	}

	chann, err := rc.adapter.ProduceMessage(constants.Insert, kfkMsg, headers)

	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		utils.SendErrorResponse(c, err)
		return
	}

	d, _, err := utils.ListenReply[types.IdResponse](chann, rc.avro, time.Second*30)

	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkResponse(c, nil, d)

}

func (rc *RestController) UpdatePartnerBalance(c *gin.Context) {
	var balancePayload partnercore.UpdatePartnerBalancePayload

	if err := c.ShouldBindJSON(&balancePayload); err != nil {
		logger.Error().Err(err).Msg("ShouldBindJSON")
		utils.SendErrorResponse(c, customErrors.ErrRequestBadRequestError)
		return
	}

	// Validate payload
	err := balancePayload.Validate()

	if err != nil {
		logger.Error().Err(err).Msg("balancePayload.Validate")
		utils.SendErrorResponse(c, customErrors.ErrRequestUnprocessableEntityError)
		return
	}

	// Verify if partner with CreateBalancePayload.PartnerId exists
	err = rc.ValidatePartnerID(c, balancePayload.PartnerId)

	if err != nil {
		logger.Error().Err(err).Msg("rc.ValidatePartnerID")
		utils.SendErrorResponse(c, customErrors.ErrRequestUnprocessableEntityError)
		return
	}

	headers := utils.FormKafkaHeaders(c, constants.Balance, constants.ActionUpdate, kafka.Settings.ReplyTopic)

	kfkMsg, err := rc.avro.Encode(balancePayload)

	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		utils.SendErrorResponse(c, err)
		return
	}

	chann, err := rc.adapter.ProduceMessage(constants.RUpdate, kfkMsg, headers)

	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		utils.SendErrorResponse(c, err)
		return
	}

	d, _, err := utils.ListenReply[types.IdResponse](chann, rc.avro, time.Second*30)

	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkResponse(c, nil, d)
}

func (rc *RestController) SendTestBalanceUpdate(c *gin.Context) {
	var update partnercore.BalanceHIDDENUpdate

	if err := c.ShouldBindJSON(&update); err != nil {
		logger.Error().Err(err).Msg("ShouldBindJSON")
		utils.SendErrorResponse(c, customErrors.ErrRequestBadRequestError)
		return
	}

	headers := utils.FormKafkaHeaders(c, constants.Balance, constants.ActionPatch, kafka.Settings.ReplyTopic)

	message, err := rc.avro.Encode(update)
	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		utils.SendErrorResponse(c, err)
		return
	}

	_, err = rc.adapter.ProduceMessage(constants.RUpdate, message, headers)

	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkResponse(c, nil, "ok")
}
