package controllers

import (
	kafka "bff/pkg/kfk"
	"github.com/gin-gonic/gin"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/constants"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/customErrors"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/HIDDEN"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/utils"
	"time"
)

func (rc *RestController) CalculateMargins(c *gin.Context) {
	logger.Debug().Msg("CalculateMargins")
	var payload HIDDEN.CalculateMarginsRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Error().Err(err).Msg("ShouldBindJSON")
		utils.SendErrorResponse(c, customErrors.ErrRequestBadRequestError)
		return
	}

	if len(payload.MccMnc) == 0 {
		logger.Error().Msg("No mccMnc was supplied")
		utils.SendErrorResponse(c, customErrors.ErrRequestBadRequestError)
		return
	}

	logger.Debug().Msg("_____CalculateMargins")

	headers := utils.FormKafkaHeaders(c, constants.Margins, constants.ActionCalculate, kafka.Settings.ReplyTopic)

	//fmt.Println("headers", headers)
	kfkMsg, err := rc.avro.Encode(payload)

	//fmt.Println("kfkMsg", kfkMsg)
	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		utils.SendErrorResponse(c, customErrors.ErrRequestBadRequestError)
		return
	}

	logger.Debug().Msg("PRODUCING MSG")
	//HIDDEN.request.margin
	chann, err := rc.adapter.ProduceMessage("Request.Margin", kfkMsg, headers)

	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		utils.SendErrorResponse(c, err)
		return
	}

	d, h, err := utils.ListenReply[[][]HIDDEN.MarginPayloadByMccMnc](chann, rc.avro, time.Second*30)

	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkSliceResponse(c, h, d)
}
