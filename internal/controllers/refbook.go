package controllers

import (
	kafka "bff/pkg/kfk"
	"encoding/json"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/utils"
	"time"

	"gitlab.HIDDEN.com/HIDDEN/sms-core/constants"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/HIDDEN"

	"github.com/gin-gonic/gin"
)

type CheckMccMncByPayload struct {
	MccMnc []string `json:"mccMnc"`
}

func (rc *RestController) Refbook(c *gin.Context) {
	refbookQuery := HIDDEN.GetRefBookByPayload{
		CountryName: c.Query("countryName"),
		MccMnc:      c.Query("mccMnc"),
		Network:     c.Query("network"),
		Criteria:    c.Query("criteria"),
		Skip:        c.Query("skip"),
		Limit:       c.Query("limit"),
	}

	headers := utils.FormKafkaHeaders(c, constants.Refbook, constants.ActionGet, kafka.Settings.ReplyTopic)

	kfkMsg, err := rc.avro.Encode(refbookQuery)

	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		utils.SendErrorResponse(c, err)
		return
	}

	chann, err := rc.adapter.ProduceMessage(constants.Request, kfkMsg, headers)

	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		utils.SendErrorResponse(c, err)
		return
	}

	d, h, err := utils.ListenReply[[]HIDDEN.RefBook](chann, rc.avro, time.Second*30)

	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkSliceResponse(c, h, d)
}

func (rc *RestController) Countries(c *gin.Context) {

	headers := utils.FormKafkaHeaders(c, constants.Refbook, constants.ActionGetGroupByCountries, kafka.Settings.ReplyTopic)

	chann, err := rc.adapter.ProduceMessage(constants.Request, []byte{}, headers)

	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		utils.SendErrorResponse(c, err)
		return
	}

	d, _, err := utils.ListenReply[[]HIDDEN.CountriesAndNets](chann, rc.avro, time.Second*30)

	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkSliceResponse(c, nil, d)
}

func (rc *RestController) MccMnc(c *gin.Context) {

	var mccMncQuery HIDDEN.CheckMccMncByPayload
	if err := c.ShouldBindJSON(&mccMncQuery); err != nil {
		logger.Error().Err(err).Msg("ShouldBindJSON")
		utils.SendErrorResponse(c, err)
		return
	}

	headers := utils.FormKafkaHeaders(c, constants.Refbook, constants.ActionCheckIfExists, kafka.Settings.ReplyTopic)

	// TODO: ADD HEADERS AND AVRO IF NEEDED
	kfkMsg, err := json.Marshal(mccMncQuery.MccMnc)

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

	d, _, err := utils.ListenReply[[]HIDDEN.ExistingMccMnc](chann, rc.avro, time.Second*30)

	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkSliceResponse(c, nil, d)
}
