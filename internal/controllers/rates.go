package controllers

import (
	kafka "bff/pkg/kfk"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/constants"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/fp"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	pricecore "gitlab.HIDDEN.com/HIDDEN/sms-core/price"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/utils"
)

func (rc *RestController) getClientProductRate(c *gin.Context, cpId string, mccMnc []string) (*pricecore.CurrentRatesForProduct, error) {
	list, _, err := rc.GetPriceListByProductId(c, pricecore.GetPriceListByProductPayload{ProductId: cpId})

	if err != nil {
		logger.Error().Err(err).Msg("rc.GetPriceListByProductId")
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, err
	}

	logger.Info().Msgf("LIST IS %+v\n", list)

	rate := fp.Filter(*list.Prices, func(p pricecore.ClientPrice) bool { return fp.Includes(mccMnc, p.Price.Mccmnc) })
	if len(rate) == 0 {
		return nil, nil
	}

	logger.Info().Msgf("Rate IS %+v\n", rate)

	return &pricecore.CurrentRatesForProduct{
		ProductId: cpId,
		Rate:      rate[0].Price.Rate,
		MccMnc:    rate[0].Price.Mccmnc,
	}, nil
}

func (rc *RestController) getVendorProductsRates(c *gin.Context, payload pricecore.GetCurrentVendorProductRatesPayload) ([]pricecore.CurrentRatesForProduct, error) {
	// Get vp rates from prices collection
	//headers := utils.FormKafkaHeaders(c, constants.Prices, constants.GetCurrentRatesForProducts, kafka.Settings.ReplyTopic)
	headers := utils.FormKafkaHeaders(c, constants.Prices, constants.ActionGetCurrentRates, kafka.Settings.ReplyTopic)

	message, err := rc.avro.Encode(payload)
	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		return nil, fmt.Errorf("failed to encode GetCurrentRatesForProducts message: %w", err)
	}

	channel, err := rc.adapter.ProduceMessage(constants.Request, message, headers)

	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		return nil, err
	}

	d, _, err := utils.ListenReply[[]pricecore.CurrentRatesForProduct](channel, rc.avro, time.Second*30)

	if err != nil {
		logger.Error().Err(err).Msg("ListenReply")
		return nil, err
	}

	return *d, nil
}

func (rc *RestController) GetCurrentRatesForProducts(c *gin.Context, payload pricecore.GetCurrentRatesForProductsPayload) ([]pricecore.CurrentRatesForProduct, error) {
	logger.Info().Msgf("GetCurrentRatesForProducts %+v\n", payload)

	cpRate, err := rc.getClientProductRate(c, payload.CPId, payload.MccMnc)
	if err != nil {
		logger.Error().Err(err).Msg("rc.getClientProductRate")
		return nil, err
	}

	vpRates, err := rc.getVendorProductsRates(c, pricecore.GetCurrentVendorProductRatesPayload{MccMnc: payload.MccMnc, VPIds: payload.VPIds})
	if err != nil {
		logger.Error().Err(err).Msg("rc.getVendorProductsRates")
		return nil, err
	}

	rates := make([]pricecore.CurrentRatesForProduct, len(vpRates))
	copy(rates, vpRates)

	logger.Info().Msgf("rates %+v\n", rates)

	if cpRate != nil {
		rates = append(rates, *cpRate)
	}

	return rates, err
}
