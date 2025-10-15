package files

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	"strconv"

	pricecore "gitlab.HIDDEN.com/HIDDEN/sms-core/price"
)

func RatePlanToCSVBuffer(ratePlan pricecore.RatePlan) (bytes.Buffer, error) {
	var buf bytes.Buffer
	var totalRows [][]string

	headerRow := []string{"Mccmnc", "Country", "Net", "Effective from", "Rate", "Currency"}
	totalRows = append(totalRows, headerRow)

	for _, price := range ratePlan.Prices {
		var rate string
		if price.Price.Rate.Value != nil {
			rate = strconv.FormatFloat(float64(*price.Price.Rate.Value), 'g', -1, 32)
		}
		effectiveFrom := price.Price.EffectiveFrom.Format("2006/01/02 15:04")

		totalRows = append(totalRows, []string{
			price.Price.Mccmnc,
			price.Price.SystemCountry,
			price.Price.SystemNet,
			effectiveFrom,
			rate,
			price.Price.Rate.Currency,
		})
	}

	csvWriter := csv.NewWriter(&buf)

	err := csvWriter.WriteAll(totalRows)
	if err != nil {
		logger.Error().Err(err).Msg("csvWriter.WriteAll")
		return buf, err
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		logger.Error().Err(err).Msg("csvWriter.Error")
		return buf, err
	}

	return buf, err
}

func PriceListToCSVBuffer(priceList pricecore.PriceList) (bytes.Buffer, error) {
	var buf bytes.Buffer
	var totalRows [][]string

	headerRow := []string{"Mccmnc", "Country", "Net", "Effective from", "Rate", "Currency", "Plan type"}
	totalRows = append(totalRows, headerRow)

	if priceList.Prices != nil {
		for _, priceListSection := range *priceList.Prices {
			rate := strconv.FormatFloat(float64(*priceListSection.Price.Rate.Value), 'g', -1, 32)
			effectiveFrom := priceListSection.Price.EffectiveFrom.Format("2006/01/02 15:04")

			totalRows = append(totalRows, []string{
				priceListSection.Price.Mccmnc,
				priceListSection.Price.SystemCountry,
				priceListSection.Price.SystemNet,
				effectiveFrom,
				rate,
				priceListSection.Price.Rate.Currency,
				priceListSection.RatePlanType,
			})
		}
	}

	csvWriter := csv.NewWriter(&buf)

	err := csvWriter.WriteAll(totalRows)
	if err != nil {
		logger.Error().Err(err).Msg("csvWriter.WriteAll")
		return buf, err
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		logger.Error().Err(err).Msg("csvWriter.Error")
		return buf, err
	}

	return buf, err
}

func PriceListToCSVBufferV2(priceList pricecore.PriceListV2) (bytes.Buffer, error) {
	var buf bytes.Buffer
	var totalRows [][]string

	headerRow := []string{"Mccmnc", "Country", "Net", "Effective from", "Rate", "Currency", "Plan type"}
	totalRows = append(totalRows, headerRow)

	if priceList.Prices != nil {
		for _, priceListSection := range *priceList.Prices {
			rate := strconv.FormatFloat(float64(*priceListSection.Price.Rate.Value), 'g', -1, 32)
			effectiveFrom := priceListSection.Price.EffectiveFrom.Format("2006/01/02 15:04")

			totalRows = append(totalRows, []string{
				priceListSection.Price.Mccmnc,
				priceListSection.Price.SystemCountry,
				priceListSection.Price.SystemNet,
				effectiveFrom,
				rate,
				priceListSection.Price.Rate.Currency,
				priceListSection.RatePlanType,
			})
		}
	}

	csvWriter := csv.NewWriter(&buf)

	err := csvWriter.WriteAll(totalRows)
	if err != nil {
		logger.Error().Err(err).Msg("csvWriter.WriteAll")
		return buf, err
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		logger.Error().Err(err).Msg("csvWriter.Error")
		return buf, err
	}

	return buf, err
}

func PricesToCSVBuffer(prices []pricecore.QuestionPrice) (bytes.Buffer, error) {
	var buf bytes.Buffer
	var totalRows [][]string

	headerRow := []string{"Mccmnc", "Country", "Net", "Old Rate", "New Rate", "Type", "Effective from"}
	totalRows = append(totalRows, headerRow)

	for _, item := range prices {
		var rate string
		var oldRate string
		if item.Price.Rate.Value != nil {
			rate = fmt.Sprintf("%s %s", strconv.FormatFloat(float64(*item.Price.Rate.Value), 'g', -1, 32), item.Price.Rate.Currency)
		}
		if item.OldRate != nil && item.OldRate.Value != nil {
			oldRate = fmt.Sprintf("%s %s", strconv.FormatFloat(float64(*item.OldRate.Value), 'g', -1, 32), item.OldRate.Currency)
		}
		effectiveFrom := item.Price.EffectiveFrom.Format("2006/01/02 15:04")

		totalRows = append(totalRows, []string{
			item.Price.Mccmnc,
			item.Price.SystemCountry,
			item.Price.SystemNet,
			oldRate,
			rate,
			item.Price.Rate.Type,
			effectiveFrom,
		})
	}

	csvWriter := csv.NewWriter(&buf)

	err := csvWriter.WriteAll(totalRows)
	if err != nil {
		logger.Error().Err(err).Msg("csvWriter.WriteAll")
		return buf, err
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		logger.Error().Err(err).Msg("csvWriter.Error")
		return buf, err
	}

	return buf, err
}
