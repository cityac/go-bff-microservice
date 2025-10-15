package controllers

import (
	internalUtils "bff/internal/utils"
	"bff/pkg/domain/model"
	"bff/pkg/files"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	"io"
	"net/http"
	"net/url"

	"gitlab.HIDDEN.com/HIDDEN/sms-core/customErrors"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/utils"

	"gitlab.HIDDEN.com/HIDDEN/sms-core/constants"
	pricecore "gitlab.HIDDEN.com/HIDDEN/sms-core/price"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/types"

	"github.com/gin-gonic/gin"
)

func (rc *RestController) RatePlans(c *gin.Context) {
	planType := c.Query("type")
	partnerId := c.Query("partnerId")

	tdq := utils.ParseTableDataRequest(c)

	if planType == "custom" {
		planType = constants.Custom
	} else if planType == "global" {
		planType = constants.Global
	}

	var value = pricecore.GetRatePlansPayload{
		PlanType: planType,
		Skip:     tdq.Skip,
		Limit:    tdq.Limit,
		Sort:     tdq.Sort,
		SortBy:   tdq.SortBy,
		SearchBy: tdq.SearchBy,
		Criteria: tdq.Criteria,
	}

	if partnerId != "" {
		partnerCurrency, err := rc.GetPartnerCurrency(c, partnerId)

		if err != nil {
			logger.Error().Msgf("rc.GetPartnerCurrency Error: %v", err)
			utils.SendErrorResponse(c, err)
			return
		}

		value.Currency = partnerCurrency.Currency
	}

	list, h, err := rc.GetRatePlans(c, value)

	if err != nil {
		logger.Error().Msgf("Get PriceListByProductId Error: %v", err)
		utils.SendErrorResponse(c, err)
		return
	}

	if list == nil {
		utils.SendOkSliceResponse[pricecore.RatePlan](c, nil, nil)
		return
	}

	if planType == constants.Global {
		result := []model.GlobalRatePlan{}
		utils.ConvertSlice(list, &result)

		for i, plan := range result {
			plan.PlanBackwardId = list[i].BackwardRef.Id
			result[i] = plan
		}

		utils.SendOkSliceResponse(c, h, &result)
		return
	}

	result := []model.CustomRatePlan{}
	utils.ConvertSlice(list, &result)

	for i, plan := range result {
		plan.ProductId = list[i].BackwardRef.Id
		plan.PartnerId = list[i].BackwardRef.MetaId
		result[i] = plan
	}

	utils.SendOkSliceResponse(c, h, &result)
	return

}

func (rc *RestController) RatePlan(c *gin.Context) {
	productId := c.Query("productId")
	id := c.Query("id")
	if productId != "" {
		rc.RatePlanByProductId(c, productId)
	} else if id != "" {
		rc.RatePlanById(c, id)
	}
}

func (rc *RestController) RatePlanById(c *gin.Context, id string) {
	tdq := utils.ParseTableDataRequest(c)

	var value = pricecore.GetRatePlanByIdPayload{
		Id:       id,
		Skip:     tdq.Skip,
		Limit:    tdq.Limit,
		Sort:     tdq.Sort,
		SortBy:   tdq.SortBy,
		SearchBy: tdq.SearchBy,
		Criteria: tdq.Criteria,
	}

	plan, h, err := rc.GetRatePlanById(c, value)

	if err != nil {
		logger.Error().Err(err).Msg("rc.GetRatePlanById")
		utils.SendErrorResponse(c, err)
		return
	}

	if plan.PlanType != constants.Custom && plan.PlanType != constants.Global {
		logger.Error().Msgf("Failed to Get RatePlanByI %s: Undefined or unknown planType value %s", plan.ID, plan.PlanType)
		utils.SendErrorResponse(c, customErrors.ErrRequestWrongEnumValues)
		return
	}

	if plan.PlanType == constants.Custom {
		var customRatePlan model.CustomRatePlan
		utils.ConvertStruct(plan, &customRatePlan)
		customRatePlan.ProductId = plan.BackwardRef.Id
		customRatePlan.PartnerId = plan.BackwardRef.MetaId
		customRatePlan.Prices = internalUtils.AdaptRatePlanPrices(plan.Prices)

		utils.SendOkResponse(c, h, customRatePlan)
		return
	}

	var globalRatePlan model.GlobalRatePlan
	utils.ConvertStruct(plan, &globalRatePlan)
	globalRatePlan.PlanBackwardId = plan.BackwardRef.Id
	globalRatePlan.Prices = internalUtils.AdaptRatePlanPrices(plan.Prices)

	utils.SendOkResponse(c, h, globalRatePlan)
	return

}

func (rc *RestController) RatePlanByProductId(c *gin.Context, productId string) {
	planType := c.Query("type")
	tdq := utils.ParseTableDataRequest(c)

	if planType == "custom" {
		planType = constants.Custom
	} else if planType == "global" {
		planType = constants.Global
	}

	var value = pricecore.GetRatePlanByProductPayload{
		ProductId: productId,
		PlanType:  planType,
		Skip:      tdq.Skip,
		Limit:     tdq.Limit,
		Sort:      tdq.Sort,
		SortBy:    tdq.SortBy,
		SearchBy:  tdq.SearchBy,
		Criteria:  tdq.Criteria,
	}

	plan, h, err := rc.GetRatePlanByProductId(c, value)

	if err != nil {
		logger.Error().Err(err).Msg("rc.GetRatePlanByProductId")
		utils.SendErrorResponse(c, err)
		return
	}

	if plan == nil {
		utils.SendOkResponse(c, h, plan)
		return
	}

	if plan.PlanType != constants.Custom && plan.PlanType != constants.Global {
		txt := fmt.Sprintf("failed to get RatePlanByID %s: undefined or unknown planType value %s", plan.ID, plan.PlanType)
		logger.Error().Msg(txt)
		utils.SendErrorResponse(c, errors.New(txt))
		return
	}

	if plan.PlanType == constants.Custom {
		var customRatePlan model.CustomRatePlan
		utils.ConvertStruct(plan, &customRatePlan)
		customRatePlan.ProductId = plan.BackwardRef.Id
		customRatePlan.PartnerId = plan.BackwardRef.MetaId
		customRatePlan.Prices = internalUtils.AdaptRatePlanPrices(plan.Prices)

		utils.SendOkResponse(c, h, customRatePlan)
		return
	}

	var globalRatePlan model.GlobalRatePlan

	utils.ConvertStruct(plan, &globalRatePlan)
	globalRatePlan.PlanBackwardId = plan.BackwardRef.Id
	globalRatePlan.Prices = internalUtils.AdaptRatePlanPrices(plan.Prices)

	utils.SendOkResponse(c, h, globalRatePlan)

}

func (rc *RestController) DownloadRatePlan(c *gin.Context) {
	productId := c.Query("productId")
	planType := c.Query("type")
	id := c.Query("id")
	var plan *pricecore.RatePlan
	var err error

	if productId != "" {
		var payload = pricecore.GetRatePlanByProductPayload{
			ProductId: productId,
			PlanType:  planType,
			Skip:      "0",
			Limit:     "0",
		}
		plan, _, err = rc.GetRatePlanByProductId(c, payload)

		if err != nil {
			logger.Error().Err(err).Msg("rc.GetRatePlanByProductId")
			utils.SendErrorResponse(c, err)
			return
		}

	} else if id != "" {
		var payload = pricecore.GetRatePlanByIdPayload{
			Id:    id,
			Skip:  "0",
			Limit: "0",
		}
		plan, _, err = rc.GetRatePlanById(c, payload)

		if err != nil {
			logger.Error().Err(err).Msg("rc.GetRatePlanById")
			utils.SendErrorResponse(c, err)
			return
		}
	}

	if plan == nil {
		logger.Error().Msgf(fmt.Sprintf("Failed to prepare ratePlan for download"))
		utils.SendErrorResponse(c, err)
		return
	}

	buffer, err := files.RatePlanToCSVBuffer(*plan)

	if err != nil {
		logger.Error().Err(err).Msg("files.RatePlanToCSVBuffer")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkFileResponse(c, map[string]string{
		"Accept-Length": fmt.Sprintf("%d", len(buffer.Bytes())),
	}, "text/csv", buffer.Bytes())
	return

}

func (rc *RestController) GetRatePlanById(c *gin.Context, payload pricecore.GetRatePlanByIdPayload) (*pricecore.RatePlan, *types.ReplyMessageHeaders, error) {

	requestURL := fmt.Sprintf("%s/api/rate-plans/plan/%s?&sort=%s&sortBy=%s&criteria=%s&skip=%s&limit=%s", rc.config.HIDDENProviderHost, payload.Id, payload.Sort, payload.SortBy, url.QueryEscape(payload.Criteria), payload.Skip, payload.Limit)
	res, headers, err := utils.HttpGet(c, requestURL)

	if err != nil {
		logger.Error().Err(err).Msg("utils.HttpGet")
		return nil, nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error().Err(err).Msg("Body.Close()")
		}
	}(res.Body)

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)

		var rateplan *pricecore.RatePlan
		err = json.Unmarshal(bodyBytes, &rateplan)
		if err != nil {
			logger.Error().Err(err).Msg("json.Unmarshal")
			return nil, nil, err
		}

		if totalCount, ok := res.Header["X-Total-Count"]; ok && len(totalCount) > 0 {
			headers.Total = totalCount[0]
		} else {
			logger.Warn().Msg("Header 'X-Total-Count' is missing")
		}

		return rateplan, headers, nil
	}

	return nil, headers, nil
}

func (rc *RestController) GetRatePlanByProductId(c *gin.Context, payload pricecore.GetRatePlanByProductPayload) (*pricecore.RatePlan, *types.ReplyMessageHeaders, error) {

	requestURL := fmt.Sprintf("%s/api/rate-plans/plan/?&productId=%s&excludeGlobals=%t&type=%s&sort=%s&sortBy=%s&criteria=%s&skip=%s&limit=%s", rc.config.HIDDENProviderHost, payload.ProductId, true, payload.PlanType, payload.Sort, payload.SortBy, url.QueryEscape(payload.Criteria), payload.Skip, payload.Limit)
	res, headers, err := utils.HttpGet(c, requestURL)

	if err != nil {
		logger.Error().Err(err).Msg("utils.HttpGet")
		return nil, nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error().Err(err).Msg("Body.Close()")
		}
	}(res.Body)

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)

		var rateplan *pricecore.RatePlan
		err = json.Unmarshal(bodyBytes, &rateplan)

		if err != nil {
			logger.Error().Err(err).Msg("json.Unmarshal")
			return nil, nil, err
		}

		if totalCount, ok := res.Header["X-Total-Count"]; ok && len(totalCount) > 0 {
			headers.Total = totalCount[0]
		} else {
			logger.Warn().Msg("Header 'X-Total-Count' is missing")
		}

		return rateplan, headers, nil
	}

	return nil, headers, nil
}

func (rc *RestController) GetRatePlans(c *gin.Context, payload pricecore.GetRatePlansPayload) ([]pricecore.RatePlan, *types.ReplyMessageHeaders, error) {
	requestURL := fmt.Sprintf("%s/api/rate-plans?type=%s&sort=%s&sortBy=%s&criteria=%s&skip=%s&limit=%s&currency=%s", rc.config.HIDDENProviderHost, payload.PlanType, payload.Sort, payload.SortBy, url.QueryEscape(payload.Criteria), payload.Skip, payload.Limit, payload.Currency)
	res, headers, err := utils.HttpGet(c, requestURL)

	if err != nil {
		return nil, nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)

		var rateplans []pricecore.RatePlan
		json.Unmarshal(bodyBytes, &rateplans)

		if totalCount, ok := res.Header["X-Total-Count"]; ok && len(totalCount) > 0 {
			headers.Total = totalCount[0]
		} else {
			logger.Warn().Msg("Header 'X-Total-Count' is missing")
		}

		return rateplans, headers, err
	}

	return nil, headers, nil
}

func (rc *RestController) LinkRatePlan(c *gin.Context, payload pricecore.LinkProductRatePlanPayload) (*types.IdResponse, *types.ReplyMessageHeaders, error) {

	requestURL := fmt.Sprintf("%s/api/rate-plans/link", rc.config.HIDDENProviderHost)
	res, headers, err := utils.HttpPost(c, requestURL, payload)

	if err != nil {
		logger.Error().Err(err).Msg("utils.HttpPost")
		return nil, nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error().Err(err).Msg("Body.Close()")
		}
	}(res.Body)

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)

		var idObject *types.IdResponse
		err = json.Unmarshal(bodyBytes, &idObject)
		if err != nil {
			logger.Error().Err(err).Msg("json.Unmarshal")
			return nil, nil, err
		}

		return idObject, headers, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		msg := "Internal Server Error"
		err := errors.New(msg)

		return nil, nil, err
	}

	return nil, headers, nil
}

func (rc *RestController) LinkRatePlanToProduct(c *gin.Context) {

	var payload pricecore.LinkProductRatePlanPayload

	if err := c.ShouldBind(&payload); err != nil {
		logger.Error().Err(err).Msg("ShouldBind")
		utils.SendErrorResponse(c, err)
		return
	}

	resp, h, err := rc.LinkRatePlan(c, payload)

	if err != nil {
		logger.Error().Err(err).Msg("rc.LinkRatePlan")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkResponse(c, h, resp)
}
