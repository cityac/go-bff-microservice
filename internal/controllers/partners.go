package controllers

import (
	internalUtils "bff/internal/utils"
	kafka "bff/pkg/kfk"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"gitlab.HIDDEN.com/HIDDEN/sms-core/auth"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/customErrors"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/fp"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"

	"github.com/gin-gonic/gin"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/constants"
	partnercore "gitlab.HIDDEN.com/HIDDEN/sms-core/partner"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/types"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/utils"
)

type GetPartnerByPayload struct {
	BusinessUnit   string   `json:"businessUnit"`
	SubUnit        string   `json:"subUnit"`
	AccountManager string   `json:"accountManager"`
	ContractType   string   `json:"contractType"`
	Skip           string   `json:"skip"`
	Limit          string   `json:"limit"`
	Search         string   `json:"search"`
	Active         string   `json:"active"`
	SearchBy       []string `json:"searchBy"`
	Fields         []string `json:"fields"`
}

func (rc *RestController) doGetPartnersProducts(c *gin.Context, value partnercore.GetPartnerProductListPayload) ([]partnercore.PartnerProduct, *types.ReplyMessageHeaders, error) {
	headers := utils.FormKafkaHeaders(c, constants.Products, constants.ActionGetByType, kafka.Settings.ReplyTopic)

	kfkMsg, err := rc.avro.Encode(value)

	if err != nil {
		return nil, nil, err
	}

	chann, err := rc.adapter.ProduceMessage(constants.Request, kfkMsg, headers)

	if err != nil {
		return nil, nil, err
	}

	d, h, err := utils.ListenReply[[]partnercore.PartnerProduct](chann, rc.avro, time.Second*30)

	if err != nil {
		return nil, nil, err
	}
	return *d, h, nil
}

func (rc *RestController) GetPartnersProducts(c *gin.Context) {
	value := partnercore.GetPartnerProductListPayload{
		Type:     c.Query("partnerType"),
		Criteria: c.Query("criteria"),
		Skip:     c.Query("skip"),
		Limit:    c.Query("limit"),
	}
	d, h, err := rc.doGetPartnersProducts(c, value)

	if err != nil {
		logger.Error().Err(err).Msg("doGetPartnersProducts")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkSliceResponse(c, h, &d)
}

func (rc *RestController) GetVendorProducts(c *gin.Context) ([]partnercore.PartnerProduct, *types.ReplyMessageHeaders, error) {

	value := partnercore.GetPartnerProductListPayload{
		Type:     "Vendor",
		Skip:     "0",
		Limit:    "0",
		Criteria: c.Query("criteria"),
	}

	products, h, err := rc.doGetPartnersProducts(c, value)

	if err != nil {
		logger.Error().Err(err).Msg("doGetPartnersProducts")
		return nil, nil, err
	}

	return products, h, nil
}

func (rc *RestController) GetVendorPrices(c *gin.Context, productIds []string) ([]partnercore.VendorProductPrice, error) {
	payload := partnercore.GetVendorProductsClientPayload{
		MinPrice: c.Query("minPrice"),
		MaxPrice: c.Query("maxPrice"),
		Country:  c.Query("country"),
		Skip:     c.Query("skip"),
		Limit:    c.Query("limit"),
		Sort:     c.Query("sort"),
		Criteria: c.Query("criteria"),
	}

	encodedProductIds, err := utils.EncodeToBase64(productIds)
	if err != nil {
		logger.Error().Err(err).Msg("EncodeToBase64")
		return nil, err
	}

	requestURL := fmt.Sprintf("%s/api/prices/vendor-products?maxPrice=%s&minPrice=%s&country=%s&skip=%s&limit=%s&sort=%s&criteria=%s&productIds=%s", rc.config.HIDDENProviderHost, payload.MaxPrice, payload.MinPrice, url.QueryEscape(payload.Country), payload.Skip, payload.Limit, payload.Sort, url.QueryEscape(payload.Criteria), encodedProductIds)
	res, _, err := utils.HttpGet(c, requestURL)

	if err != nil {
		logger.Error().Err(err).Msg("HttpGet")
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)

		if err != nil {
			logger.Error().Err(err).Msg("ReadAll")
			return nil, err
		}

		var vendorProductPrices []partnercore.VendorProductPrice
		err = json.Unmarshal(bodyBytes, &vendorProductPrices)
		if err != nil {
			logger.Error().Err(err).Msg("Unmarshal")
			return nil, err
		}
		return vendorProductPrices, nil
	}

	return nil, nil

}

func (rc *RestController) GetVendorProductsPrices(c *gin.Context) {
	products, _, err := rc.GetVendorProducts(c)
	if err != nil {
		logger.Error().Err(err).Msg("GetVendorProducts")
		return
	}

	productIds := make([]string, 0, len(products))
	for _, product := range products {
		productIds = append(productIds, product.ID)
	}

	vendorProductPrices, err := rc.GetVendorPrices(c, productIds)
	if err != nil {
		logger.Error().Err(err).Msg("GetVendorPrices")
		utils.SendErrorResponse(c, err)
		return
	}

	minPrice, _ := strconv.ParseFloat(c.Query("minPrice"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("maxPrice"), 64)
	skip := c.Query("skip")
	limit := c.Query("limit")

	d := internalUtils.AdaptVendorProductsV2(products, vendorProductPrices, minPrice, maxPrice)

	d = internalUtils.ApplyVendorProductsPaging(d, skip, limit)

	utils.SendOkSliceResponse(c, nil, &d)
}

func (rc *RestController) ValidatePartnerID(c *gin.Context, partnerId string) error {
	headers := utils.FormKafkaHeaders(c, constants.Partners, constants.ActionGetById, kafka.Settings.ReplyTopic)

	kfkMsg, err := json.Marshal(partnerId)

	if err != nil {
		return err
	}

	chann, err := rc.adapter.ProduceMessage(constants.Request, kfkMsg, headers)

	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		return err
	}

	_, _, err = utils.ListenReply[partnercore.Partner](chann, rc.avro, time.Second*30)

	return err

}

func (rc *RestController) Partners(c *gin.Context) {
	userID := c.GetString(constants.UserId)
	orgID := c.GetString(constants.OrgId)

	ctx := context.Background()

	draftParam := c.Query("draft")

	canViewPartners, _ := rc.fga.CheckUserForPermissions(ctx, userID, orgID, auth.AccessPartnersViewer)
	canViewDrafts, _ := rc.fga.CheckUserForPermissions(ctx, userID, orgID, auth.AccessPartnersDraftViewer)

	if draftParam == "false" && !canViewPartners {
		utils.SendEmptySliceResponse(c)
		return
	}

	if draftParam == "true" && !canViewDrafts {
		utils.SendEmptySliceResponse(c)
		return
	}

	value := partnercore.GetPartnerByPayload{
		BusinessUnit:   c.Query("businessUnit"),
		SubUnit:        c.Query("subUnit"),
		AccountManager: c.Query("accountManager"),
		ContractType:   c.Query("contractType"),
		Skip:           c.Query("skip"),
		Limit:          c.Query("limit"),
		Search:         c.Query("search"),
		Active:         c.Query("active"),
		Draft:          draftParam,
		SearchBy:       []string{"name", "alias"},
	}

	headers := utils.FormKafkaHeaders(c, constants.Partners, constants.ActionGet, kafka.Settings.ReplyTopic)

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

	d, h, err := utils.ListenReply[[]partnercore.Partner](chann, rc.avro, time.Second*30)

	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkSliceResponse(c, h, d)
}

func (rc *RestController) GetPartnerById(c *gin.Context) {
	partnerId := c.Param("partnerId")

	headers := utils.FormKafkaHeaders(c, constants.Partners, constants.ActionGetById, kafka.Settings.ReplyTopic)

	kfkMsg, err := json.Marshal(partnerId)

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

	d, _, err := utils.ListenReply[partnercore.Partner](chann, rc.avro, time.Second*30)

	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkResponse(c, nil, d)
}

func (rc *RestController) GetChannelsActiveCount(c *gin.Context) {
	productId := c.Param("productId")

	headers := utils.FormKafkaHeaders(c, constants.Channels, constants.ActionGetActiveCount, kafka.Settings.ReplyTopic)

	payload := partnercore.NewGetChannelCountPayload(productId)

	kfkMsg, err := rc.avro.Encode(payload)

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

	d, _, err := utils.ListenReply[types.EntityCount](chann, rc.avro, time.Second*30)

	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkResponse(c, nil, d)
}

func (rc *RestController) GetProductById(c *gin.Context) {
	productId := c.Param("productId")

	headers := utils.FormKafkaHeaders(c, constants.Products, constants.ActionGetById, kafka.Settings.ReplyTopic)

	kfkMsg, err := json.Marshal(productId)

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

	d, _, err := utils.ListenReply[partnercore.Product](chann, rc.avro, time.Second*30)

	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkResponse(c, nil, d)
}

func (rc *RestController) checkChannelHIDDENSets(c *gin.Context, partner partnercore.InsertPartnerWithRefs) ([]partnercore.CheckedChannel, error) {

	headers := utils.FormKafkaHeaders(c, constants.Channels, constants.ActionCheckUniqueFields, kafka.Settings.ReplyTopic)

	var payload partnercore.CheckChannelHIDDENSetsPayload

	for _, product := range partner.Products {
		for _, channel := range product.Channels {
			payload.HIDDENSets = append(payload.HIDDENSets, partnercore.ChannelHIDDENUniqueSet{
				Login:     channel.Login,
				Password:  channel.Password,
				Addresses: channel.Addresses,
			})
		}
	}

	checkChannelsKfkMsg, err := rc.avro.Encode(payload)

	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		return nil, err
	}

	checkChannelsChann, err := rc.adapter.ProduceMessage(constants.Request, checkChannelsKfkMsg, headers)

	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		return nil, err
	}

	d, _, err := utils.ListenReply[[]partnercore.CheckedChannel](checkChannelsChann, rc.avro, time.Second*30)

	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		return nil, err
	}

	return *d, err
}

func (rc *RestController) TransformProductHIDDENIds(c *gin.Context, HIDDENIds []int) ([]string, error) {
	headers := utils.FormKafkaHeaders(c, constants.Products, constants.ActionTransformHIDDENIds, kafka.Settings.ReplyTopic)

	message, err := rc.avro.Encode(HIDDENIds)
	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		return nil, fmt.Errorf("failed to encode TransformProductHIDDENIds message: %w", err)
	}

	channel, err := rc.adapter.ProduceMessage(constants.Request, message, headers)

	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		return nil, err
	}

	d, _, err := utils.ListenReply[[]string](channel, rc.avro, time.Second*30)

	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		return nil, err
	}

	return *d, nil
}

/*NEW PARTNER*/
func (rc *RestController) checkPartnerName(c *gin.Context, partner partnercore.PartnerClientPayload, partnerID string) ([]partnercore.CheckedPartner, error) {

	headers := utils.FormKafkaHeaders(c, constants.Partners, constants.ActionCheckNames, kafka.Settings.ReplyTopic)

	var partnerNamesPayload []string

	partnerNamesPayload = append(partnerNamesPayload, partner.Name)

	partnerNameKfkMsg, err := json.Marshal(partnerNamesPayload)

	if err != nil {
		logger.Error().Err(err).Msg("json.Marshal")
		return nil, err
	}

	chann, err := rc.adapter.ProduceMessage(constants.Request, partnerNameKfkMsg, headers)

	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		return nil, err
	}

	d, _, err := utils.ListenReply[[]partnercore.CheckedPartner](chann, rc.avro, time.Second*30)

	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		return nil, err
	}

	// Exclude current Partner Name
	filteredPartners := fp.Filter(*d, func(result partnercore.CheckedPartner) bool {
		return result.ID != partnerID
	})

	return filteredPartners, err
}

func (rc *RestController) CreatePartner(c *gin.Context) {
	var clientPayload partnercore.PartnerClientPayload

	if err := c.ShouldBindJSON(&clientPayload); err != nil {
		logger.Error().Err(err).Msg("ShouldBindJSON")
		utils.SendErrorResponse(c, err, customErrors.JSONMarshal)
		return
	}

	err := clientPayload.Validate()
	if err != nil {
		logger.Error().Err(err).Msg("clientPayload.Validate")
		utils.SendErrorResponse(c, err, customErrors.UnprocessableEntity)
		return
	}

	/*Create kafka payload*/
	payload := partnercore.NewCreatePartnerPayload(clientPayload)

	// GetFeatureFlags if partner name is unique in DB
	nonUniquePartnerNames, err := rc.checkPartnerName(c, clientPayload, payload.ID)

	if len(nonUniquePartnerNames) > 0 {
		log.Println("ERROR: Partner name is already exists")

		msg := fmt.Sprintf("Partner with name %v is already exists", clientPayload.Name)
		err = errors.New(msg)

		utils.SendErrorResponse(c, err, customErrors.RequestConflictError)
		return
	}

	headers := utils.FormKafkaHeaders(c, constants.Partners, constants.ActionCreate, kafka.Settings.ReplyTopic)

	kfkMsg, err := rc.avro.Encode(payload)
	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		utils.SendErrorResponse(c, err)
		return
	}

	chann, err := rc.adapter.ProduceMessage(constants.RCreate, kfkMsg, headers)
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

func (rc *RestController) UpdatePartner(c *gin.Context) {
	var clientPayload partnercore.PartnerClientPayload
	partnerId := c.Param("partnerId")

	if err := c.ShouldBindJSON(&clientPayload); err != nil {
		logger.Error().Err(err).Msg("ShouldBindJSON")
		utils.SendErrorResponse(c, customErrors.ErrRequestBadRequestError)
		return
	}

	err := clientPayload.Validate()
	if err != nil {
		logger.Error().Err(err).Msg("clientPayload.Validate")
		utils.SendErrorResponse(c, err, customErrors.UnprocessableEntity)
		return
	}

	// GetFeatureFlags if partner name is unique in DB
	nonUniquePartnerNames, err := rc.checkPartnerName(c, clientPayload, partnerId)
	if len(nonUniquePartnerNames) > 0 {
		msg := fmt.Sprintf("partner with name %v is already exists", clientPayload.Name)
		err = errors.New(msg)
		logger.Warn().Err(err).Msg("Partner name is already exists")

		utils.SendErrorResponse(c, err, customErrors.RequestConflictError)
		return
	}

	/*Create kafka payload*/
	kafkaPayload := partnercore.NewUpdatePartnerPayload(clientPayload, partnerId)

	headers := utils.FormKafkaHeaders(c, constants.Partners, constants.ActionUpdate, kafka.Settings.ReplyTopic)

	kfkMsg, err := rc.avro.Encode(kafkaPayload)
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

func (rc *RestController) DeactivatePartner(c *gin.Context) {
	partnerId := c.Param("partnerId")

	productsAndChannels, err := rc.getProductsAndChannelsByPartnerId(c, partnerId)
	if err != nil {
		logger.Error().Err(err).Msg("rc.getProductsAndChannelsByPartnerId")
		utils.SendErrorResponse(c, err)
		return
	}

	/*Create kafka payload*/
	kafkaPayload := partnercore.NewPartnerProductsAndChannelsUpdatePayload(false, partnerId, productsAndChannels)

	headers := utils.FormKafkaHeaders(c, constants.Partners, constants.ActionDeactivate, kafka.Settings.ReplyTopic)

	kfkMsg, err := rc.avro.Encode(kafkaPayload)
	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		utils.SendErrorResponse(c, err)
		return
	}

	chann, err := rc.adapter.ProduceMessage(constants.RUpdate, kfkMsg, headers)
	if err != nil {
		utils.SendErrorResponse(c, err)
		return
	}

	_, _, err = utils.ListenReply[types.IdResponse](chann, rc.avro, time.Second*30)
	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendNoContentResponse(c)
}

func (rc *RestController) TogglePartner(c *gin.Context) {
	var clientPayload partnercore.TogglePartner
	partnerId := c.Param("partnerId")

	if err := c.ShouldBindJSON(&clientPayload); err != nil {
		utils.SendErrorResponse(c, err, customErrors.JSONMarshal)
		return
	}

	logger.Info().Msgf("ID: %v", partnerId)
	logger.Info().Msgf("TOGGLE: %v", clientPayload)

	productsAndChannels, err := rc.getProductsAndChannelsByPartnerId(c, partnerId)
	if err != nil {
		logger.Error().Err(err).Msg("rc.getProductsAndChannelsByPartnerId")
		utils.SendErrorResponse(c, err)
		return
	}

	/*Create kafka payload*/
	kafkaPayload := partnercore.NewPartnerProductsAndChannelsUpdatePayload(clientPayload.Toggle, partnerId, productsAndChannels)

	headers := utils.FormKafkaHeaders(c, constants.Partners, constants.ActionToggle, kafka.Settings.ReplyTopic)

	kfkMsg, err := rc.avro.Encode(kafkaPayload)
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

	_, _, err = utils.ListenReply[types.IdResponse](chann, rc.avro, time.Second*30)
	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendNoContentResponse(c)
}

func (rc *RestController) CreatePartnerProduct(c *gin.Context) {
	var clientPayload partnercore.ProductClientPayload
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

	isExists, err := rc.checkIsProductAlreadyExists(c, partnerId, partnercore.CheckProductPayload{Name: clientPayload.Name})

	if err != nil {
		logger.Error().Err(err).Msg("rc.checkIsProductAlreadyExists")
		utils.SendErrorResponse(c, err, customErrors.BadRequest)
		return
	}

	if isExists {
		err = errors.New("make sure the product has unique name")
		logger.Error().Err(err).Msg("An error occurred")
		utils.SendErrorResponse(c, err, customErrors.UnprocessableEntity)
		return
	}

	/*Create kafka payload*/
	kafkaPayload := partnercore.NewCreateProductPayload(clientPayload, partnerId)

	headers := utils.FormKafkaHeaders(c, constants.Products, constants.ActionCreate, kafka.Settings.ReplyTopic)

	kfkMsg, err := rc.avro.Encode(kafkaPayload)
	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		utils.SendErrorResponse(c, err)
		return
	}

	channel, err := rc.adapter.ProduceMessage(constants.RUpdate, kfkMsg, headers)
	if err != nil {
		utils.SendErrorResponse(c, err)
		return
	}

	d, _, err := utils.ListenReply[partnercore.Product](channel, rc.avro, time.Second*30)
	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	product := *d

	if len(product.Tags) == 0 {
		product.Tags = make([]partnercore.Tag, 0)
	}

	utils.SendCreatedResponse(c, product)
}

func (rc *RestController) UpdatePartnerProduct(c *gin.Context) {
	var clientPayload partnercore.UpdateProductClientPayload
	partnerId := c.Param("partnerId")
	productId := c.Param("productId")

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

	isExists, err := rc.checkIsProductAlreadyExists(c, partnerId, partnercore.CheckProductPayload{Name: clientPayload.Name, ProductID: productId})

	if err != nil {
		logger.Error().Err(err).Msg("rc.checkIsProductAlreadyExists")
		utils.SendErrorResponse(c, err, customErrors.BadRequest)
		return
	}

	if isExists {
		err = errors.New("make sure the product has unique name")
		logger.Error().Err(err).Msg("rc.checkIsProductAlreadyExists")
		utils.SendErrorResponse(c, err, customErrors.UnprocessableEntity)
		return
	}

	/*Create kafka payload*/
	kafkaPayload := partnercore.NewUpdateProductPayload(clientPayload, partnerId, productId)

	headers := utils.FormKafkaHeaders(c, constants.Products, constants.ActionUpdate, kafka.Settings.ReplyTopic)

	kfkMsg, err := rc.avro.Encode(kafkaPayload)

	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		utils.SendErrorResponse(c, err)
		return
	}

	channel, err := rc.adapter.ProduceMessage(constants.RUpdate, kfkMsg, headers)
	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		utils.SendErrorResponse(c, err)
		return
	}

	d, _, err := utils.ListenReply[partnercore.Product](channel, rc.avro, time.Second*30)
	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	product := *d

	if product.Tags == nil {
		product.Tags = make([]partnercore.Tag, 0)
	}

	utils.SendOkResponse(c, nil, product)
}

func (rc *RestController) GetPartnerProductsById(c *gin.Context) {
	value := partnercore.GetProductsByPartnerIdParams{
		Skip:      c.Query("skip"),
		Limit:     c.Query("limit"),
		Search:    c.Query("search"),
		PartnerId: c.Param("partnerId"),
	}

	headers := utils.FormKafkaHeaders(c, constants.Products, constants.ActionGetByPartnerId, kafka.Settings.ReplyTopic)

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

	d, h, err := utils.ListenReply[[]partnercore.Product](chann, rc.avro, time.Second*30)

	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	data := fp.Map(*d, func(p partnercore.Product) partnercore.Product {
		if p.Tags == nil {
			p.Tags = make([]partnercore.Tag, 0)
		}
		return p
	})

	utils.SendOkSliceResponse(c, h, &data)
}

func (rc *RestController) DeletePartnerProduct(c *gin.Context) {
	partnerId := c.Param("partnerId")
	productId := c.Param("productId")

	headers := utils.FormKafkaHeaders(c, constants.Products, constants.ActionGetById, kafka.Settings.ReplyTopic)

	kfkMsg, err := json.Marshal(productId)

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

	d, _, err := utils.ListenReply[partnercore.Product](chann, rc.avro, time.Second*30)
	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	/*Create kafka payload*/
	kafkaPayload := partnercore.NewDeleteProductPayload(partnerId, productId, d.HIDDENId, d.Type)

	headers = utils.FormKafkaHeaders(c, constants.Products, constants.ActionDelete, kafka.Settings.ReplyTopic)

	kfkMsg, err = rc.avro.Encode(kafkaPayload)
	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		utils.SendErrorResponse(c, err)
		return
	}

	channel, err := rc.adapter.ProduceMessage(constants.RDelete, kfkMsg, headers)
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

func (rc *RestController) CreateProductChannel(c *gin.Context) {
	fmt.Println("CreateProductChannel")
	var clientPayload partnercore.ChannelClientPayload
	productId := c.Param("productId")

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

	var checkedChanPayload partnercore.CheckChannelPayload
	utils.ConvertStruct(clientPayload, &checkedChanPayload)

	checkedChannel, err := rc.checkIsChannelAlreadyExists(c, checkedChanPayload)
	if err != nil {
		logger.Error().Err(err).Msg("rc.checkIsChannelAlreadyExists")
		utils.SendErrorResponse(c, err, customErrors.BadRequest)
		return
	}

	if !checkedChannel.UniqueName {
		err = errors.New("make sure the channel has a unique name")
		logger.Error().Err(err).Msg("An error occurred")
		utils.SendErrorResponse(c, err, customErrors.UnprocessableEntity)
		return
	}

	if !checkedChannel.UniqueHostNameAndLogin {
		err = errors.New("make sure the channel has a unique name")
		logger.Error().Err(err).Msg("An error occurred")
		utils.SendErrorResponse(c, err, customErrors.UnprocessableEntity)
		return
	}

	headers := utils.FormKafkaHeaders(c, constants.Products, constants.ActionGetById, kafka.Settings.ReplyTopic)

	kfkMsg, err := json.Marshal(productId)

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

	prod, _, err := utils.ListenReply[partnercore.Product](chann, rc.avro, time.Second*30)

	/*Create kafka payload*/
	kafkaPayload := partnercore.NewCreateChannelPayload(clientPayload, productId, prod.Name, prod.PartnerId)

	headers = utils.FormKafkaHeaders(c, constants.Channels, constants.ActionCreate, kafka.Settings.ReplyTopic)

	kfkMsg, err = rc.avro.Encode(kafkaPayload)
	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		utils.SendErrorResponse(c, err)
		return
	}

	channel, err := rc.adapter.ProduceMessage(constants.RCreate, kfkMsg, headers)
	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		utils.SendErrorResponse(c, err)
		return
	}

	d, _, err := utils.ListenReply[partnercore.Channel](channel, rc.avro, time.Second*30)
	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	if len(d.ChannelInfo.AllowedNetworks) == 0 {
		d.ChannelInfo.AllowedNetworks = make([]string, 0)
	}

	r := d.ToChannelWithProductName(prod.Name)

	utils.SendCreatedResponse(c, r)
}

func (rc *RestController) UpdateProductChannel(c *gin.Context) {
	logger.Info().Msg("UpdateProductChannel")
	var clientPayload partnercore.UpdateChannelClientPayload
	productId := c.Param("productId")
	channelId := c.Param("channelId")

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

	checkedChannel, err := rc.checkIsChannelAlreadyExists(c, partnercore.CheckChannelPayload{
		Name:      clientPayload.Name,
		Login:     clientPayload.Login,
		Password:  clientPayload.Password,
		Addresses: clientPayload.Addresses,
		ProductId: productId,
		ChannelId: channelId,
	})
	if err != nil {
		logger.Error().Err(err).Msg("rc.checkIsChannelAlreadyExists")
		utils.SendErrorResponse(c, err, customErrors.BadRequest)
		return
	}

	if !checkedChannel.UniqueName {
		err = errors.New("make sure the channel has a unique name")
		logger.Error().Err(err).Msg("An error occurred")
		utils.SendErrorResponse(c, fmt.Errorf("Make sure the channel has a unique name"), customErrors.UnprocessableEntity)
		return
	}

	if !checkedChannel.UniqueHostNameAndLogin {
		utils.SendErrorResponse(c, fmt.Errorf("Make sure the channel has a unique hostname or login"), customErrors.UnprocessableEntity)
		return
	}

	headers := utils.FormKafkaHeaders(c, constants.Products, constants.ActionGetById, kafka.Settings.ReplyTopic)

	kfkMsg, err := json.Marshal(productId)

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

	prod, _, err := utils.ListenReply[partnercore.Product](chann, rc.avro, time.Second*30)
	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	/*GetFeatureFlags previous state*/
	request := types.IdRequest{Id: channelId}
	kfkMsg, err = rc.avro.Encode(request)
	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		utils.SendErrorResponse(c, err)
		return
	}

	headers = utils.FormKafkaHeaders(c, constants.Channels, constants.ActionGetById, kafka.Settings.ReplyTopic)
	chann, err = rc.adapter.ProduceMessage(constants.Request, kfkMsg, headers)

	prevState, _, err := utils.ListenReply[partnercore.Channel](chann, rc.avro, time.Second*30)
	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	/*GetFeatureFlags whether product activation/deactivation is needed*/
	if prevState.Active != clientPayload.Active {
		var productRequest partnercore.UpdateProductPayload
		needUpdate := false

		/*Product deactivation flow*/
		if prevState.Active && prod.Active {
			chanCountRequest := partnercore.GetChannelCountPayload{ProductId: productId}
			kfkMsg, err = rc.avro.Encode(chanCountRequest)
			if err != nil {
				logger.Error().Err(err).Msg("rc.avro.Encode")
				utils.SendErrorResponse(c, err)
				return
			}

			headers = utils.FormKafkaHeaders(c, constants.Channels, constants.ActionGetActiveCount, kafka.Settings.ReplyTopic)
			chann, err = rc.adapter.ProduceMessage(constants.Request, kfkMsg, headers)

			activeCount, _, err := utils.ListenReply[types.EntityCount](chann, rc.avro, time.Second*30)
			if err != nil {
				logger.Error().Err(err).Msg("utils.ListenReply")
				utils.SendErrorResponse(c, err)
				return
			}

			/*The channel is the last one active, so we deactivate product*/
			if activeCount.Total == 1 {
				productRequest = partnercore.NewUpdateProductActivePayload(prod, false)
				needUpdate = true
			}
		}

		/*Product activation flow*/
		if !prevState.Active && !prod.Active {
			chanCountRequest := partnercore.GetChannelCountPayload{ProductId: productId}
			kfkMsg, err = rc.avro.Encode(chanCountRequest)
			if err != nil {
				logger.Error().Err(err).Msg("rc.avro.Encode")
				utils.SendErrorResponse(c, err)
				return
			}

			headers = utils.FormKafkaHeaders(c, constants.Channels, constants.ActionCount, kafka.Settings.ReplyTopic)
			chann, err = rc.adapter.ProduceMessage(constants.Request, kfkMsg, headers)

			totalCount, _, err := utils.ListenReply[types.EntityCount](chann, rc.avro, time.Second*30)
			if err != nil {
				logger.Error().Err(err).Msg("utils.ListenReply")
				utils.SendErrorResponse(c, err)
				return
			}

			/*The channel is the only one existing, so we activate product*/
			if totalCount.Total == 1 {
				productRequest = partnercore.NewUpdateProductActivePayload(prod, true)
				needUpdate = true
			}
		}

		/*Update the product if one of the cases fired*/
		if needUpdate {
			kfkMsg, err = rc.avro.Encode(productRequest)
			if err != nil {
				logger.Error().Err(err).Msg("rc.avro.Encode")
				utils.SendErrorResponse(c, err)
				return
			}

			headers := utils.FormKafkaHeaders(c, constants.Products, constants.ActionUpdate, kafka.Settings.ReplyTopic)
			channel, err := rc.adapter.ProduceMessage(constants.RUpdate, kfkMsg, headers)
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
		}
	}

	/*Create kafka payload*/
	kafkaPayload := partnercore.NewUpdateChannelPayload(clientPayload, channelId, prod.Name, prod.PartnerId)

	headers = utils.FormKafkaHeaders(c, constants.Channels, constants.ActionUpdate, kafka.Settings.ReplyTopic)

	kfkMsg, err = rc.avro.Encode(kafkaPayload)

	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		utils.SendErrorResponse(c, err)
		return
	}

	channel, err := rc.adapter.ProduceMessage(constants.RUpdate, kfkMsg, headers)
	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		utils.SendErrorResponse(c, err)
		return
	}

	d, _, err := utils.ListenReply[partnercore.Channel](channel, rc.avro, time.Second*30)
	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	if len(d.ChannelInfo.AllowedNetworks) == 0 {
		d.ChannelInfo.AllowedNetworks = make([]string, 0)
	}

	r := d.ToChannelWithProductName(prod.Name)

	utils.SendOkResponse(c, nil, r)
}

func (rc *RestController) GetChannelsByProductId(c *gin.Context) {
	productId := c.Param("productId")

	headers := utils.FormKafkaHeaders(c, constants.Channels, constants.ActionGetByProductId, kafka.Settings.ReplyTopic)

	kfkMsg, err := json.Marshal(productId)

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

	d, _, err := utils.ListenReply[[]partnercore.Channel](chann, rc.avro, time.Second*30)

	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	data := fp.Map(*d, func(p partnercore.Channel) partnercore.Channel {
		if p.Addresses == nil {
			p.Addresses = make([]string, 0)
		}

		if p.ChannelInfo.AllowedNetworks == nil {
			p.ChannelInfo.AllowedNetworks = make([]string, 0)
		}

		return p
	})

	utils.SendOkSliceResponse(c, nil, &data)
}

func (rc *RestController) getChannelsByPartnerIdInternal(c *gin.Context, skip, limit, search, partnerId string) (*[]partnercore.ChannelWithProductName, *types.ReplyMessageHeaders, error) {
	params := partnercore.GetChannelsByPartnerId{
		Skip:      skip,
		Limit:     limit,
		Search:    search,
		PartnerId: partnerId,
	}

	headers := utils.FormKafkaHeaders(c, constants.Partners, constants.ActionFindAllPartnerChannels, kafka.Settings.ReplyTopic)

	kfkMsg, err := rc.avro.Encode(params)

	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		return nil, nil, err
	}

	chann, err := rc.adapter.ProduceMessage(constants.Request, kfkMsg, headers)

	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		return nil, nil, err
	}

	d, h, err := utils.ListenReply[[]partnercore.ChannelWithProductName](chann, rc.avro, time.Second*30)
	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		return nil, nil, err
	}

	data := fp.Map(*d, func(p partnercore.ChannelWithProductName) partnercore.ChannelWithProductName {
		if p.Addresses == nil {
			p.Addresses = make([]string, 0)
		}

		if p.ChannelInfo.AllowedNetworks == nil {
			p.ChannelInfo.AllowedNetworks = make([]string, 0)
		}

		return p
	})

	return &data, h, nil
}

func (rc *RestController) GetChannelsByPartnerId(c *gin.Context) {
	logger.Info().Msg("GetChannelsByPartnerId")

	data, h, err := rc.getChannelsByPartnerIdInternal(
		c,
		c.Query("skip"),
		c.Query("limit"),
		c.Query("search"),
		c.Param("partnerId"),
	)

	if err != nil {
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkSliceResponse(c, h, data)
}

func (rc *RestController) DeleteProductChannel(c *gin.Context) {
	productId := c.Param("productId")
	channelId := c.Param("channelId")
	partnerId := c.Param("partnerId")

	isTheLastOne, notFount, HIDDENId, err := rc.checkIsChannelTheLastOne(c, partnerId, channelId)

	if err != nil {
		logger.Error().Err(err).Msg("rc.checkIsChannelTheLastOne")
		utils.SendErrorResponse(c, err)
		return
	}

	if notFount {
		logger.Warn().Msg("not_fount")
		utils.SendErrorResponse(c, err, customErrors.BadRequest)
		return
	}

	if isTheLastOne {
		err = errors.New("the last one channel cannot be deleted")
		logger.Warn().Err(err).Msg("Warn")
		utils.SendErrorResponse(c, err, customErrors.RequestConflictError)
		return
	}
	// GET CHANNEL HIDDEN ID

	/*Create kafka payload*/
	kafkaPayload := partnercore.NewDeleteChannelPayload(productId, channelId, HIDDENId)

	headers := utils.FormKafkaHeaders(c, constants.Channels, constants.ActionDelete, kafka.Settings.ReplyTopic)

	kfkMsg, err := rc.avro.Encode(kafkaPayload)
	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		utils.SendErrorResponse(c, err)
		return
	}

	channel, err := rc.adapter.ProduceMessage(constants.RDelete, kfkMsg, headers)
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

func (rc *RestController) checkIsProductAlreadyExists(c *gin.Context, partnerID string, payload partnercore.CheckProductPayload) (bool, error) {
	headers := utils.FormKafkaHeaders(c, constants.Products, constants.ActionGetByPartnerId, kafka.Settings.ReplyTopic)

	value := partnercore.GetProductsByPartnerIdParams{
		Skip:      "0",
		Limit:     "0",
		PartnerId: partnerID,
	}

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

	d, _, err := utils.ListenReply[[]partnercore.Product](chann, rc.avro, time.Second*30)

	isProductAlreadyExists := fp.Some(*d, func(p partnercore.Product) bool {
		return p.Name == payload.Name && p.ID != payload.ProductID
	})

	return isProductAlreadyExists, nil
}

func (rc *RestController) checkIsChannelAlreadyExists(c *gin.Context, p partnercore.CheckChannelPayload) (*partnercore.IsChannelExists, error) {

	headers := utils.FormKafkaHeaders(c, constants.Channels, constants.ActionCheckIfExists, kafka.Settings.ReplyTopic)

	kfkMsg, err := json.Marshal(p)
	if err != nil {
		logger.Error().Err(err).Msg("json.Marshal")
		return nil, err
	}

	m, err := rc.adapter.ProduceMessage(constants.Request, kfkMsg, headers)
	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		return nil, err
	}

	d, _, err := utils.ListenReply[partnercore.IsChannelExists](m, rc.avro, time.Second*30)
	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		return nil, err
	}

	return d, nil
}

func (rc *RestController) checkIsChannelTheLastOne(c *gin.Context, partnerId string, channelId string) (bool, bool, *int64, error) {
	fmt.Println("checkIsChannelTheLastOne")

	headers := utils.FormKafkaHeaders(c, constants.Partners, constants.ActionFindAllPartnerChannels, kafka.Settings.ReplyTopic)

	payload := partnercore.GetChannelsByPartnerId{
		PartnerId: partnerId,
		Skip:      "0",
		Limit:     "0",
	}

	kfkMsg, err := rc.avro.Encode(payload)
	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		return false, false, nil, err
	}

	m, err := rc.adapter.ProduceMessage(constants.Request, kfkMsg, headers)
	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		return false, false, nil, err
	}

	d, _, err := utils.ListenReply[[]partnercore.Channel](m, rc.avro, time.Second*30)

	//fmt.Printf("FILTERED CHANNELS %+v\n", filteredChannels)
	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		return false, false, nil, err
	}

	filteredChannels := fp.Filter(*d, func(c partnercore.Channel) bool {
		return c.ID == channelId
	})

	fmt.Printf("FILTERED CHANNELS %+v\n", filteredChannels)

	if len(filteredChannels) == 0 {
		return false, true, nil, fmt.Errorf("cannot find channel with this id for this partner")
	}
	HIDDENID := filteredChannels[0].HIDDENId
	channelsLength := len(*d)
	fmt.Println("channelsLength", channelsLength)

	return channelsLength == 1, false, HIDDENID, nil
}

func (rc *RestController) getProductsAndChannelsByPartnerId(c *gin.Context, partnerID string) (*partnercore.PartnerProductsAndChannels, error) {

	headers := utils.FormKafkaHeaders(c, constants.Partners, constants.ActionFindAllPartnerProductsAndChannels, kafka.Settings.ReplyTopic)

	kfkMsg, err := json.Marshal(partnerID)
	if err != nil {
		logger.Error().Err(err).Msg("json.Marshal")
		return nil, err
	}

	chann, err := rc.adapter.ProduceMessage(constants.Request, kfkMsg, headers)
	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		return nil, err
	}

	d, _, err := utils.ListenReply[partnercore.PartnerProductsAndChannels](chann, rc.avro, time.Second*30)

	return d, nil
}

func (rc *RestController) GetPartnerCurrency(c *gin.Context, partnerID string) (*partnercore.PartnerCurrency, error) {
	headers := utils.FormKafkaHeaders(c, constants.Partners, constants.ActionGetCurrency, kafka.Settings.ReplyTopic)

	payload := partnercore.GetPartnerCurrencyPayload{
		PartnerId: partnerID,
	}

	kfkMsg, err := rc.avro.Encode(payload)
	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		return nil, err
	}

	chann, err := rc.adapter.ProduceMessage(constants.Request, kfkMsg, headers)
	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		return nil, err
	}

	d, _, err := utils.ListenReply[partnercore.PartnerCurrency](chann, rc.avro, time.Second*30)

	return d, nil
}

