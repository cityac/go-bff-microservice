package controllers

import (
	"bff/pkg/domain/model"
	"bff/pkg/files"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"

	"gitlab.HIDDEN.com/HIDDEN/sms-core/types"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/utils"

	pricecore "gitlab.HIDDEN.com/HIDDEN/sms-core/price"

	"github.com/gin-gonic/gin"
)

func (rc *RestController) PriceListByProductId(c *gin.Context) {
	productId := c.Param("productId")
	tdq := utils.ParseTableDataRequest(c)

	var value = pricecore.GetPriceListByProductPayload{
		ProductId: productId,
		Skip:      tdq.Skip,
		Limit:     tdq.Limit,
		Sort:      tdq.Sort,
		SortBy:    tdq.SortBy,
		SearchBy:  tdq.SearchBy,
		Criteria:  tdq.Criteria,
	}

	list, h, err := rc.GetPriceListByProductId(c, value)

	if err != nil {
		logger.Error().Err(err).Msg("GetPriceListByProductId")
		utils.SendErrorResponse(c, err)
		return
	}

	var result model.PriceList
	decoded, err := json.Marshal(list)

	if err != nil {
		logger.Error().Err(err).Msg("Marshal")
		utils.SendErrorResponse(c, err)
		return
	}

	err = json.Unmarshal(decoded, &result)

	if err != nil {
		logger.Error().Err(err).Msg("Unmarshal")
		utils.SendErrorResponse(c, err)
		return
	}
	result.DailyPriceData = nil

	if list.Prices != nil {
		for i, price := range *list.Prices {
			p := model.PriceListPrice{
				RatePlanType: price.RatePlanType,
				Price: model.Price{
					ID:            price.Price.PriceUploadId,
					Mccmnc:        price.Price.Mccmnc,
					EffectiveFrom: price.Price.EffectiveFrom,
					EffectiveTo:   price.Price.EffectiveTo,
					Rate:          price.Price.Rate,
					SystemCountry: price.Price.SystemCountry,
					SystemNet:     price.Price.SystemNet,
					ProductId:     price.Price.ProductId,
				},
			}

			(*result.Prices)[i] = p
		}
	}

	utils.SendOkResponse(c, h, result)
}

func (rc *RestController) DownloadPriceList(c *gin.Context) {
	productId := c.Param("productId")
	tdq := utils.ParseTableDataRequest(c)

	var value = pricecore.GetPriceListByProductPayload{
		ProductId: productId,
		Skip:      tdq.Skip,
		Limit:     tdq.Limit,
		Sort:      tdq.Sort,
		SortBy:    tdq.SortBy,
		SearchBy:  tdq.SearchBy,
		Criteria:  tdq.Criteria,
	}

	list, _, err := rc.GetPriceListByProductId(c, value)

	if err != nil {
		logger.Error().Err(err).Msg("GetPriceListByProductId")
		utils.SendErrorResponse(c, err)
		return
	}

	buffer, err := files.PriceListToCSVBufferV2(*list)

	if err != nil {
		logger.Error().Err(err).Msg("PriceListToCSVBuffer")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkFileResponse(c, map[string]string{
		"Accept-Length": fmt.Sprintf("%d", len(buffer.Bytes())),
	}, "text/csv", buffer.Bytes())
}

func (rc *RestController) GetPriceListByProductId(c *gin.Context, payload pricecore.GetPriceListByProductPayload) (*pricecore.PriceListV2, *types.ReplyMessageHeaders, error) {

	requestURL := fmt.Sprintf("%s/api/prices/pricelist/%s?sort=%s&sortBy=%s&criteria=%s&skip=%s&limit=%s", rc.config.HIDDENProviderHost, payload.ProductId, payload.Sort, payload.SortBy, url.QueryEscape(payload.Criteria), payload.Skip, payload.Limit)

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

		var pricelist *pricecore.PriceListV2
		err = json.Unmarshal(bodyBytes, &pricelist)

		if err != nil {
			logger.Error().Err(err).Msg("Unmarshal")
			return nil, nil, err
		}

		headers.Total = res.Header["X-Total-Count"][0]

		return pricelist, headers, err
	}

	return nil, headers, nil
}
