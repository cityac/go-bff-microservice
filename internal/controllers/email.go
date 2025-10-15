package controllers

import (
	"bff/pkg/domain/model"
	"bff/pkg/files"
	kafka "bff/pkg/kfk"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/constants"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/customErrors"
	izonet "gitlab.HIDDEN.com/HIDDEN/sms-core/izonet"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/mailservice"
	partnercore "gitlab.HIDDEN.com/HIDDEN/sms-core/partner"
	pricecore "gitlab.HIDDEN.com/HIDDEN/sms-core/price"
	sms_HIDDEN "gitlab.HIDDEN.com/HIDDEN/sms-core/HIDDEN"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/utils"
)

func (rc *RestController) SendEmail(c *gin.Context) {
	logger.Trace().Msg("SendEmail")

	var request mailservice.SendEmailPriceListPayload
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error().Err(err).Msg("Failed to bind request body")
		utils.SendErrorResponse(c, err)
		return
	}
	logger.Trace().Interface("request", request).Msg("Successfully parsed request body")

	// Get product Info
	logger.Trace().Msg("Starting to fetch product info")
	headers := utils.FormKafkaHeaders(c, constants.Products, constants.ActionGetById, kafka.Settings.ReplyTopic)

	kfkMsg, err := json.Marshal(request.ProductID)

	if err != nil {
		logger.Error().Err(err).Msg("json.Marshal")
		utils.SendErrorResponse(c, err)
		return
	}
	logger.Trace().Str("productID", request.ProductID).Msg("Successfully marshaled ProductID")

	chann, err := rc.adapter.ProduceMessage(constants.Request, kfkMsg, headers)

	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		utils.SendErrorResponse(c, err)
		return
	}
	logger.Trace().Msg("Successfully produced Kafka message for product info")

	product, _, err := utils.ListenReply[partnercore.Product](chann, rc.avro, time.Second*30)

	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}
	logger.Trace().Interface("product", product).Msg("Successfully received product info")

	// Get partner Info
	logger.Trace().Msg("Starting to fetch partner info")
	headers = utils.FormKafkaHeaders(c, constants.Partners, constants.ActionGetById, kafka.Settings.ReplyTopic)

	kfkMsg, err = json.Marshal(request.PartnerID)

	if err != nil {
		logger.Error().Err(err).Msg("json.Marshal")
		utils.SendErrorResponse(c, err)
		return
	}
	logger.Trace().Str("partnerID", request.PartnerID).Msg("Successfully marshaled PartnerID")

	chann, err = rc.adapter.ProduceMessage(constants.Request, kfkMsg, headers)

	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		utils.SendErrorResponse(c, err)
		return
	}
	logger.Trace().Msg("Successfully produced Kafka message for partner info")

	partner, _, err := utils.ListenReply[partnercore.Partner](chann, rc.avro, time.Second*30)

	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}
	logger.Trace().Interface("partner", partner).Msg("Successfully received partner info")

	formattedPriceList, err := rc.getPriceListInfo(c, request.ProductID)
	if err != nil {
		logger.Error().Err(err).Msg("getPriceListInfo")
		utils.SendErrorResponse(c, err)
		return
	}

	// Some test info
	//formattedPriceList := model.FormattedPriceList{
	//	Rows: []model.PriceListRow{
	//		{
	//			MCC:           "284",
	//			MNC:           "01",
	//			Country:       "Bulgaria",
	//			Network:       "A1",
	//			EffectiveDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	//			Rate:          "0.15",
	//			ChangeType:    "increase",
	//		},
	//		{
	//			MCC:           "284",
	//			MNC:           "05",
	//			Country:       "Bulgaria",
	//			Network:       "Yettel",
	//			EffectiveDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	//			Rate:          "0.18",
	//			ChangeType:    "decrease",
	//		},
	//		{
	//			MCC:           "284",
	//			MNC:           "03",
	//			Country:       "Bulgaria",
	//			Network:       "Vivacom",
	//			EffectiveDate: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
	//			Rate:          "0.20",
	//			ChangeType:    "new",
	//		},
	//		{
	//			MCC:           "255",
	//			MNC:           "01",
	//			Country:       "Ukraine",
	//			Network:       "Kyivstar",
	//			EffectiveDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	//			Rate:          "0.25",
	//			ChangeType:    "increase",
	//		},
	//		{
	//			MCC:           "425",
	//			MNC:           "01",
	//			Country:       "Israel",
	//			Network:       "Partner",
	//			EffectiveDate: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
	//			Rate:          "0.22",
	//			ChangeType:    "new",
	//		},
	//	},
	//}

	// Create XLS file
	xlsBuffer, err := files.PriceListToXLSBuffer(formattedPriceList, partner.Name, product.Name)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create XLS buffer")
		utils.SendErrorResponse(c, err)
		return
	}

	systemID := "[]"
	logger.Trace().Msg("Starting to fetch channels info")

	if product.HIDDENId != nil {
		logger.Trace().Int64("HIDDENId", *product.HIDDENId).Msg("Found HIDDENId, fetching channels")
		channels, err := rc.getChannels(c, *product.HIDDENId)
		if err != nil {
			logger.Error().Err(err).Msg("rc.getChannels")
		}
		if len(channels) > 0 {
			systemID = channels[0].SystemID
			logger.Trace().Str("systemID", systemID).Msg("Successfully got systemID from channels")
		}
	}

	logger.Trace().Msg("Starting to fetch contract company info")
	contractCompanyName := partner.ContractCompany

	contractCompany, err := rc.getContractCompany(c, request.ProductID)
	if err != nil {
		logger.Error().Err(err).Msg("error to get contractCompany from Izonet")
	} else {
		contractCompanyName = contractCompany.Name
	}
	logger.Trace().Interface("contractCompany", contractCompany).Msg("Successfully received contract company info")

	serviceRequest := mailservice.NewEmailRequest(
		partner.Name,
		product.Name,
		contractCompanyName,
		request.Email,
		systemID,
		xlsBuffer.Bytes(),
	)
	logger.Trace().Interface("serviceRequest", serviceRequest).Msg("Created email service request")

	// Encode message using Avro
	encodedData, err := rc.avro.Encode(serviceRequest)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to encode email message")
		utils.SendErrorResponse(c, err)
		return
	}
	logger.Trace().Msg("Successfully encoded email message")

	// Send message to Kafka
	headers = []sarama.RecordHeader{
		{Key: []byte("ReplyTopic"), Value: []byte(kafka.Settings.ReplyTopic)},
	}

	channel, err := rc.adapter.ProduceMessage("Requests.Email", encodedData, headers)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to send email message to Kafka")
		utils.SendErrorResponse(c, err)
		return
	}
	logger.Trace().Msg("Successfully sent email message to Kafka")

	// Wait for response
	response, _, err := utils.ListenReply[mailservice.EmailResponse](channel, rc.avro, time.Second*30)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get response from email service")
		utils.SendErrorResponse(c, err)
		return
	}
	logger.Trace().Interface("response", response).Msg("Received response from email service")

	if response == nil || !response.Success {
		if response != nil && response.Error != nil {
			logger.Error().Str("response.Error", *response.Error).Msg("Email service temporarily unavailable")
		} else {
			logger.Error().Err(err).Msg("Email service temporarily unavailable")
		}
		utils.SendErrorResponse(c, errors.New(customErrors.StatusServiceUnavailable))
		return
	}

	logger.Trace().Msg("Email sent successfully")
	utils.SendNoContentResponse(c)
}

func (rc *RestController) getPriceListInfo(c *gin.Context, productID string) (model.FormattedPriceList, error) {
	// Get price list info
	logger.Trace().Msg("Starting to fetch price list info")

	var value = pricecore.GetPriceListByProductPayload{
		ProductId: productID,
	}

	priceList, _, err := rc.GetPriceListByProductId(c, value)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get price list")
		return model.FormattedPriceList{}, err
	}

	// Format price list data
	formattedPriceList := model.FormattedPriceList{
		Rows: make([]model.PriceListRow, 0),
	}

	now := time.Now()

	if priceList.Prices != nil {
		for _, price := range *priceList.Prices {
			// Split mccmnc into mcc and mnc
			mcc := price.Price.Mccmnc[:3]
			mnc := price.Price.Mccmnc[3:]

			// Format rate
			rateStr := ""
			if price.Price.Rate.Value != nil {
				rateStr = fmt.Sprintf("%.3f %s", *price.Price.Rate.Value, price.Price.Rate.Currency)
			}

			// Determine change type
			changeType := "Future"
			if price.Price.EffectiveFrom.Before(now) && price.Price.EffectiveTo.After(now) {
				changeType = "Current"
			}

			row := model.PriceListRow{
				MCC:           mcc,
				MNC:           mnc,
				Country:       price.Price.SystemCountry,
				Network:       price.Price.SystemNet,
				EffectiveDate: price.Price.EffectiveFrom,
				Rate:          rateStr,
				ChangeType:    changeType,
			}

			formattedPriceList.Rows = append(formattedPriceList.Rows, row)
		}
	}
	return formattedPriceList, nil
}

func (rc *RestController) getContractCompany(c *gin.Context, productId string) (*izonet.IzonetContractCompanyResponse, error) {
	logger.Trace().Str("productId", productId).Msg("Starting getContractCompany request")

	requestURL := fmt.Sprintf("%s/api/izonet/contract-company/%s", rc.config.HIDDENProviderHost, productId)

	logger.Trace().Str("requestURL", requestURL).Msg("Formed request URL")

	res, _, err := utils.HttpGet(c, requestURL)
	if err != nil {
		logger.Error().Err(err).Msg("utils.HttpGet")
		return nil, err
	}
	logger.Trace().Int("statusCode", res.StatusCode).Msg("Received HTTP response")

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error().Err(err).Msg("Body.Close")
		}
	}(res.Body)

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			logger.Error().Err(err).Msg("io.ReadAll")
			return nil, err
		}
		logger.Trace().Int("bodyLength", len(bodyBytes)).Msg("Successfully read response body")

		var contractCompany izonet.IzonetContractCompanyResponse
		err = json.Unmarshal(bodyBytes, &contractCompany)
		if err != nil {
			logger.Error().Err(err).Msg("json.Unmarshal")
			return nil, err
		}
		logger.Trace().Interface("contractCompany", contractCompany).Msg("Successfully unmarshaled contract company")

		return &contractCompany, nil
	}

	return nil, fmt.Errorf("failed to get contract company: status code %d", res.StatusCode)
}

func (rc *RestController) getChannels(c *gin.Context, productId int64) ([]sms_HIDDEN.Channel, error) {
	logger.Info().Int64("productId", productId).Msg("Starting getChannels request")

	organizationId := c.Request.Header.Get(constants.OrgIdHeader)
	logger.Trace().Str("organizationId", organizationId).Msg("Got organization ID from header")

	requestURL := fmt.Sprintf("%s/channels/%s/%d?organizationId=%s",
		rc.config.HIDDENBrokerHost,
		"client",
		productId,
		organizationId,
	)
	logger.Trace().Str("requestURL", requestURL).Msg("Formed request URL")

	res, _, err := utils.HttpGet(c, requestURL)
	if err != nil {
		logger.Error().Err(err).Msg("utils.HttpGet")
		return nil, err
	}
	logger.Trace().Int("statusCode", res.StatusCode).Msg("Received HTTP response")

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error().Err(err).Msg("Body.Close")
		}
	}(res.Body)

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			logger.Error().Err(err).Msg("io.ReadAll")
			return nil, err
		}
		logger.Trace().Int("bodyLength", len(bodyBytes)).Msg("Successfully read response body")

		var channels []sms_HIDDEN.Channel
		err = json.Unmarshal(bodyBytes, &channels)
		if err != nil {
			logger.Error().Err(err).Msg("json.Unmarshal")
			return nil, err
		}
		logger.Trace().Int("channelsCount", len(channels)).Msg("Successfully unmarshaled channels")

		return channels, nil
	}

	return nil, fmt.Errorf("failed to get channels: status code %d", res.StatusCode)
}
