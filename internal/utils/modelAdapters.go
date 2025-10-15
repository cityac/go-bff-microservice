package utils

import (
	"bff/pkg/domain/model"
	"fmt"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	"sort"
	"strconv"

	partnercore "gitlab.HIDDEN.com/HIDDEN/sms-core/partner"
	pricecore "gitlab.HIDDEN.com/HIDDEN/sms-core/price"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/utils"
)

func AdaptRatePlanPrices(prices []pricecore.RatePlanPrice) []model.Price {
	result := make([]model.Price, len(prices))

	for i, price := range prices {
		var p model.Price
		utils.ConvertStruct(price.Price, &p)
		result[i] = p
	}

	return result
}
func AdaptVendorProducts(products *[]partnercore.VendorProduct, minPrice float64, maxPrice float64) []model.VendorProduct {
	result := make([]model.VendorProduct, len(*products))
	for i, product := range *products {
		var p model.VendorProduct
		utils.ConvertStruct(product, &p)
		p.MccMnc = fmt.Sprintf("%d", product.Count)

		network := sort.StringSlice(UniqueSliceByProp[partnercore.VendorProductPrice, string](product.ProductPrices, "Network"))
		sort.Sort(network)
		p.Network = network

		vendorPrices := make([]model.VendorProductPrice, 0, len(p.VendorPrices))

		if len(product.ProductPrices) == 1 {
			price := product.ProductPrices[0]
			mergedProduct := model.VendorProduct{
				MccMnc:       price.Mccmnc,
				Network:      []string{price.Network},
				Rate:         fmt.Sprintf("%.3f %s", *price.Rate.Value, price.Rate.Currency),
				Name:         product.Name,
				VendorPrices: make([]model.VendorProductPrice, 0),
			}
			result[i] = mergedProduct
		} else {
			for _, price := range product.ProductPrices {
				vendorPrices = append(vendorPrices, model.VendorProductPrice{
					ID:      price.ID,
					Mccmnc:  price.Mccmnc,
					Network: price.Network,
					Rate:    fmt.Sprintf("%.3f %s", *price.Rate.Value, price.Rate.Currency),
				})
			}
			p.VendorPrices = vendorPrices
			p.Rate = fmt.Sprintf("%.3f - %.3f EUR", minPrice, maxPrice)
			result[i] = p
		}
	}

	return result
}

func AdaptVendorProductsV2(products []partnercore.PartnerProduct, vendorPrices []partnercore.VendorProductPrice, minPrice float64, maxPrice float64) []model.VendorProduct {
	result := make([]model.VendorProduct, 0)

	vendorPricesToProcess := make([]partnercore.VendorProductPrice, len(vendorPrices))
	copy(vendorPricesToProcess, vendorPrices)

	for _, product := range products {
		prices := make([]partnercore.VendorProductPrice, 0)

		processedIdsSet := make(map[string]struct{}, 0)
		for _, price := range vendorPricesToProcess {
			if price.ProductId == product.ID {

				prices = append(prices, price)
				processedIdsSet[price.ID] = struct{}{}

			}
		}
		filteredVendorPrices := vendorPricesToProcess[:0]
		for _, vendorPrice := range vendorPricesToProcess {
			if _, found := processedIdsSet[vendorPrice.ID]; !found {
				filteredVendorPrices = append(filteredVendorPrices, vendorPrice)
			}
		}

		vendorPricesToProcess = filteredVendorPrices

		network := sort.StringSlice(UniqueSliceByProp[partnercore.VendorProductPrice, string](prices, "Network"))
		sort.Sort(network)

		if len(prices) > 0 {
			if len(prices) == 1 {
				price := prices[0]
				mergedProduct := model.VendorProduct{
					Id:           product.ID,
					MccMnc:       price.Mccmnc,
					Network:      []string{price.Network},
					Rate:         fmt.Sprintf("%.3f %s", *price.Rate.Value, price.Rate.Currency),
					Name:         product.Name,
					VendorPrices: make([]model.VendorProductPrice, 0),
				}
				result = append(result, mergedProduct)
			} else {
				p := model.VendorProduct{
					Id:      product.ID,
					Name:    product.Name,
					Network: network,
					MccMnc:  fmt.Sprintf("%d", len(prices)),
				}
				vendorPrices := make([]model.VendorProductPrice, 0, len(p.VendorPrices))
				for _, price := range prices {
					vendorPrices = append(vendorPrices, model.VendorProductPrice{
						ID:      price.ID,
						Mccmnc:  price.Mccmnc,
						Network: price.Network,
						Rate:    fmt.Sprintf("%.3f %s", *price.Rate.Value, price.Rate.Currency),
					})
				}
				p.VendorPrices = vendorPrices
				p.Rate = fmt.Sprintf("%.3f - %.3f EUR", minPrice, maxPrice)
				result = append(result, p)
			}
		}
	}

	return result
}

func ApplyVendorProductsPaging(vendorProducts []model.VendorProduct, skip, limit string) []model.VendorProduct {
	skipInt, err := strconv.Atoi(skip)
	if err != nil {
		logger.Warn().Err(err).Msgf("skipInt %v", skipInt)
		skipInt = 0
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		logger.Warn().Err(err).Msgf("limitInt %v", limitInt)
		limitInt = len(vendorProducts)
	}

	if skipInt < 0 {
		skipInt = 0
	}

	if limitInt > len(vendorProducts) || limitInt < 0 {
		limitInt = len(vendorProducts)
	}

	if skipInt > len(vendorProducts) {
		return []model.VendorProduct{}
	}

	end := skipInt + limitInt
	if end > len(vendorProducts) {
		end = len(vendorProducts)
	}

	return vendorProducts[skipInt:end]
}
