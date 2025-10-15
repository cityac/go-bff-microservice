package controllers

import (
	kafka "bff/pkg/kfk"
	"errors"
	"github.com/gin-gonic/gin"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/constants"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/customErrors"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/fp"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	partnercore "gitlab.HIDDEN.com/HIDDEN/sms-core/partner"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/types"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/utils"
	"reflect"
	"time"
)

func (rc *RestController) CreateContact(c *gin.Context) {
	var clientPayload partnercore.ContactClientPayload
	partnerId := c.Param("partnerId")

	if err := c.ShouldBindJSON(&clientPayload); err != nil {
		logger.Error().Err(err).Msg("ShouldBindJSON")
		utils.SendErrorResponse(c, err, customErrors.JSONMarshal)
		return
	}

	/*Validate clientPayload*/
	err := clientPayload.Validate()
	if err != nil {
		logger.Error().Err(err).Msg("clientPayload.Validate")
		utils.SendErrorResponse(c, err, customErrors.UnprocessableEntity)
		return
	}

	isExists, err := rc.checkIsContactAlreadyExists(c, partnerId, partnercore.CheckContactPayload{
		Name:         clientPayload.Name,
		Email:        clientPayload.Email,
		ContactTypes: clientPayload.ContactTypes,
	})

	if err != nil {
		logger.Error().Err(err).Msg("rc.checkIsContactAlreadyExists")
		utils.SendErrorResponse(c, err, customErrors.BadRequest)
		return
	}

	if isExists {
		err = errors.New("make sure the contact has unique data")
		logger.Error().Err(err).Msg("An error occurred")
		utils.SendErrorResponse(c, errors.New("make sure the contact has unique data"), customErrors.UnprocessableEntity)
		return
	}

	/*Create kafka payload*/
	kafkaPayload := partnercore.NewCreateContactPayload(clientPayload, partnerId)

	headers := utils.FormKafkaHeaders(c, constants.Contacts, constants.ActionCreate, kafka.Settings.ReplyTopic)

	kfkMsg, err := rc.avro.Encode(kafkaPayload)
	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		utils.SendErrorResponse(c, err)
		return
	}

	channel, err := rc.adapter.ProduceMessage(constants.Insert, kfkMsg, headers)
	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		utils.SendErrorResponse(c, err)
		return
	}

	d, _, err := utils.ListenReply[partnercore.Contact](channel, rc.avro, time.Second*30)
	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendCreatedResponse(c, d)
}

func (rc *RestController) UpdateContact(c *gin.Context) {
	var clientPayload partnercore.ContactClientPayload
	partnerId := c.Param("partnerId")
	contactId := c.Param("contactId")

	if err := c.ShouldBindJSON(&clientPayload); err != nil {
		logger.Error().Err(err).Msg("ShouldBindJSON")
		utils.SendErrorResponse(c, err, customErrors.JSONMarshal)
		return
	}

	/*Validate clientPayload*/
	err := clientPayload.Validate()
	if err != nil {
		logger.Error().Err(err).Msg("clientPayload.Validate")
		utils.SendErrorResponse(c, err, customErrors.UnprocessableEntity)
		return
	}

	isExists, err := rc.checkIsContactAlreadyExists(c, partnerId, partnercore.CheckContactPayload{
		ContactID:    contactId,
		Name:         clientPayload.Name,
		Email:        clientPayload.Email,
		ContactTypes: clientPayload.ContactTypes,
	})

	if err != nil {
		logger.Error().Err(err).Msg("rc.checkIsContactAlreadyExists")
		utils.SendErrorResponse(c, err, customErrors.BadRequest)
		return
	}

	if isExists {
		err = errors.New("make sure the contact has unique data")
		logger.Warn().Err(err).Msg("rc.checkIsContactAlreadyExists")
		utils.SendErrorResponse(c, err, customErrors.UnprocessableEntity)
		return
	}

	/*Create kafka payload*/
	kafkaPayload := partnercore.NewUpdateContactPayload(clientPayload, partnerId, contactId)

	headers := utils.FormKafkaHeaders(c, constants.Contacts, constants.ActionUpdate, kafka.Settings.ReplyTopic)

	kfkMsg, err := rc.avro.Encode(kafkaPayload)
	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		utils.SendErrorResponse(c, err)
		return
	}

	channel, err := rc.adapter.ProduceMessage(constants.Insert, kfkMsg, headers)
	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		utils.SendErrorResponse(c, err)
		return
	}

	d, _, err := utils.ListenReply[partnercore.Contact](channel, rc.avro, time.Second*30)
	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkResponse(c, nil, d)
}

func (rc *RestController) GetContactsByPartnerId(c *gin.Context) {
	value := partnercore.GetContactsByPartnerIdParams{
		Skip:      c.Query("skip"),
		Limit:     c.Query("limit"),
		Search:    c.Query("search"),
		PartnerId: c.Param("partnerId"),
	}

	headers := utils.FormKafkaHeaders(c, constants.Contacts, constants.ActionGetByPartnerId, kafka.Settings.ReplyTopic)

	kfkMsg, err := rc.avro.Encode(value)
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

	d, h, err := utils.ListenReply[[]partnercore.Contact](chann, rc.avro, time.Second*30)

	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkSliceResponse(c, h, d)
}

func (rc *RestController) DeleteContact(c *gin.Context) {
	partnerId := c.Param("partnerId")
	contactId := c.Param("contactId")

	/*Create kafka payload*/
	kafkaPayload := partnercore.NewDeleteContactPayload(partnerId, contactId)

	headers := utils.FormKafkaHeaders(c, constants.Contacts, constants.ActionDeleteOneById, kafka.Settings.ReplyTopic)

	kfkMsg, err := rc.avro.Encode(kafkaPayload)
	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		utils.SendErrorResponse(c, err)
		return
	}

	channel, err := rc.adapter.ProduceMessage(constants.Insert, kfkMsg, headers)
	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		utils.SendErrorResponse(c, err)
		return
	}

	_, _, err = utils.ListenReply[types.IdResponse](channel, rc.avro, time.Second*30)
	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendNoContentResponse(c)
}

func (rc *RestController) checkIsContactAlreadyExists(c *gin.Context, partnerID string, p partnercore.CheckContactPayload) (bool, error) {
	value := partnercore.GetContactsByPartnerIdParams{
		Skip:      c.Query("skip"),
		Limit:     c.Query("limit"),
		PartnerId: c.Param("partnerId"),
	}

	headers := utils.FormKafkaHeaders(c, constants.Contacts, constants.ActionGetByPartnerId, kafka.Settings.ReplyTopic)

	kfkMsg, err := rc.avro.Encode(value)
	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		return false, err
	}

	chann, err := rc.adapter.ProduceMessage(constants.Request, kfkMsg, headers)

	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		return false, err
	}

	d, _, err := utils.ListenReply[[]partnercore.Contact](chann, rc.avro, time.Second*30)

	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		return false, err
	}

	isContactAlreadyExists := fp.Some(*d, func(c partnercore.Contact) bool {
		return c.ID != p.ContactID && c.Name == p.Name && c.Email == p.Email && reflect.DeepEqual(c.ContactTypes, p.ContactTypes)
	})

	return isContactAlreadyExists, nil
}
