package model

import (
	"time"

	pricecore "gitlab.HIDDEN.com/HIDDEN/sms-core/price"
)

type Price struct {
	ID            string         `json:"id"`
	Mccmnc        string         `json:"mccmnc"`
	Country       string         `json:"country"`
	Net           string         `json:"net"`
	EffectiveFrom time.Time      `json:"effectiveFrom"`
	EffectiveTo   time.Time      `json:"effectiveTo"`
	Rate          pricecore.Rate `json:"rate"`
	ProductId     string         `json:"productId"`
	SystemCountry string         `json:"systemCountry"`
	SystemNet     string         `json:"systemNet"`
}

type TimelinePrice struct {
	EffectiveFrom time.Time `json:"effectiveFrom"`
	EffectiveTo   time.Time `json:"effectiveTo"`
	Rate          string    `json:"rate"`
	Current       bool      `json:"current"`
	ProductId     string    `json:"productId"`
}

type OldCurrentRatePrice struct {
	Price
	OldRate *pricecore.Rate `json:"oldRate"`
}
