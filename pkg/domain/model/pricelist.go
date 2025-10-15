package model

import (
	pricecore "gitlab.HIDDEN.com/HIDDEN/sms-core/price"
)

type PriceList struct {
	pricecore.PriceList
	Prices *[]PriceListPrice `json:"prices"`
}

type PriceListPrice struct {
	Price
	RatePlanType string `json:"ratePlanType"`
}
