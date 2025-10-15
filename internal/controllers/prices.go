package controllers

import (
	"bff/pkg/domain/model"
	"bff/pkg/files"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"gitlab.HIDDEN.com/HIDDEN/sms-core/customErrors"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/fp"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/types"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/utils"

	pricecore "gitlab.HIDDEN.com/HIDDEN/sms-core/price"

	"github.com/gin-gonic/gin"
)

func (rc *RestController) GetPrices(c *gin.Context, action string) ([]pricecore.QuestionPrice, *types.ReplyMessageHeaders, error) {

	productId := c.Query("productId")
	if productId == "" {
		productId = c.Param("productId")
	}
	typeParam := c.Query("type")
	country := c.Query("country")
	tdq := utils.ParseTableDataRequest(c)

	requestURL := fmt.Sprintf("%s/api/prices/%s?countryName=%s&productId=%s&type=%s&criteria=%s&skip=%s&limit=%s", rc.config.HIDDENProviderHost, action, country, productId, typeParam, url.QueryEscape(tdq.Criteria), tdq.Skip, tdq.Limit)
	res, headers, err := utils.HttpGet(c, requestURL)

	if err != nil {
		logger.Error().Err(err).Msg("HttpGet")
		return nil, nil, err
	}

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
		}

		var prices []pricecore.QuestionPrice
		err = json.Unmarshal(bodyBytes, &prices)
		if err != nil {
			logger.Error().Err(err).Msg("json.Unmarshal")
			return nil, nil, err
		}

		prices = fp.Map(prices, func(p pricecore.QuestionPrice) pricecore.QuestionPrice {
			p.Price.ID = p.Price.PriceUploadId
			return p
		})

		prices = fp.Filter(prices, func(p pricecore.QuestionPrice) bool {
			return p.Price.PriceType != "Global"
		})
		headers.Total = res.Header["X-Total-Count"][0]

		return prices, headers, err
	}

	return nil, headers, nil
}

func (rc *RestController) DownloadLastPrices(c *gin.Context) {
	prices, headers, err := rc.GetPrices(c, "last")

	if headers != nil && headers.Error != nil {
		logger.Error().Err(err).Msg("headers.Error")
		utils.SendErrorResponse(c, errors.New(*headers.Error), customErrors.BadRequest)
		return
	}

	if err != nil {
		logger.Error().Err(err).Msg("rc.GetPrices")
		utils.SendErrorResponse(c, err)
		return
	}

	productName := c.Query("product")
	partnerName := c.Query("partner")

	if prices == nil {
		utils.SendOkSliceResponse[pricecore.QuestionPrice](c, nil, nil)
		return
	}

	buffer, err := files.PricesToCSVBuffer(prices)

	if err != nil {
		logger.Error().Err(err).Msg("files.PricesToCSVBuffer")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkFileResponse(c, map[string]string{
		"Accept-Length":       fmt.Sprintf("%d", len(buffer.Bytes())),
		"Content-Disposition": fmt.Sprintf(`attachment; filename="Last Price changes for the %s-%s.csv"`, partnerName, productName),
	}, "text/csv", buffer.Bytes())
}

func (rc *RestController) LastPrices(c *gin.Context) {
	prices, headers, err := rc.GetPrices(c, "last")

	if err != nil {
		logger.Error().Err(err).Msg("rc.GetPrices")
		utils.SendErrorResponse(c, err)
		return
	}

	if headers != nil && headers.Error != nil {
		logger.Error().Err(err).Msg("headers.Error")
		utils.SendErrorResponse(c, errors.New(*headers.Error), customErrors.BadRequest)
		return
	}

	result := fp.Map(prices, func(p pricecore.QuestionPrice) model.OldCurrentRatePrice {
		var price model.Price
		utils.ConvertStruct(p.Price, &price)
		return model.OldCurrentRatePrice{
			OldRate: p.OldRate,
			Price:   price,
		}
	})

	utils.SendOkSliceResponse(c, headers, &result)
}

func (rc *RestController) CurrentPrices(c *gin.Context) {
	prices, headers, err := rc.GetPrices(c, "current")

	if err != nil {
		logger.Error().Err(err).Msg("rc.GetPrices")
		utils.SendErrorResponse(c, err)
		return
	}

	if headers != nil && headers.Error != nil {
		logger.Error().Err(err).Msg("headers.Error")
		utils.SendErrorResponse(c, errors.New(*headers.Error), customErrors.BadRequest)
		return
	}

	result := fp.Map(prices, func(p pricecore.QuestionPrice) model.Price {
		var price model.Price
		utils.ConvertStruct(p.Price, &price)
		return price
	})

	utils.SendOkSliceResponse(c, headers, &result)
}

func (rc *RestController) PriceTimeline(c *gin.Context) {
	productId := c.Query("productId")
	typeParam := c.Query("type")
	forward := c.Query("forward")
	backward := c.Query("backward")
	mccmnc := c.Query("mccmnc")

	var forwardPayload = 1
	var backwardPayload = 2

	if forward != "" {
		f, err := strconv.Atoi(forward)
		if err == nil {
			forwardPayload = f
		} else {
			logger.Error().Err(err).Msg("strconv.Atoi(forward)")
		}
	}
	if backward != "" {
		b, err := strconv.Atoi(backward)
		if err == nil {
			backwardPayload = b
		} else {
			logger.Error().Err(err).Msg("strconv.Atoi(backward)")
		}
	}

	requestURL := fmt.Sprintf("%s/api/prices/timeline?productId=%s&type=%s&mccmnc=%s&forward=%d&backward=%d", rc.config.HIDDENProviderHost, productId, typeParam, mccmnc, forwardPayload, backwardPayload)

	res, headers, err := utils.HttpGet(c, requestURL)

	if err != nil {
		logger.Error().Err(err).Msg("utils.HttpGet")
		utils.SendErrorResponse(c, err)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error().Err(err).Msg("Body.Close()")
		}
	}(res.Body)

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			logger.Error().Msgf("Get PriceListByProductId Error: %v", err)
			utils.SendErrorResponse(c, err)
			return
		}

		var prices []pricecore.QuestionPrice
		err = json.Unmarshal(bodyBytes, &prices)
		if err != nil {
			logger.Error().Err(err).Msg("json.Unmarshal")
			utils.SendErrorResponse(c, err)
			return
		}

		timeLine := fp.Map(prices, func(p pricecore.QuestionPrice) model.TimelinePrice {
			var rate string
			if p.Price.Rate.Type == "CLOSE" {
				rate = "Closed price"
			} else if p.Price.Rate.Value != nil {
				logger.Info().Msgf("p.Price.Rate.Value")
				rate, err = utils.FormatRate(p.Price.Rate.Value, p.Price.Rate.Currency)
				if err != nil {
					logger.Error().Err(err).Msg("FormatRate")
				}
			}

			return model.TimelinePrice{
				EffectiveFrom: p.Price.EffectiveFrom,
				EffectiveTo:   p.Price.EffectiveTo,
				Rate:          rate,
				Current:       p.Current,
				ProductId:     p.Price.ProductId,
			}
		})
		utils.SendOkSliceResponse(c, headers, &timeLine)
		return
	}

	utils.SendErrorResponse(c, errors.New(res.Status), customErrors.BadRequest)
}
