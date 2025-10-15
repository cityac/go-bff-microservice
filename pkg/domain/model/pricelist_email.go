package model

import "time"

type PriceListRow struct {
	MCC           string    `json:"mcc"`
	MNC           string    `json:"mnc"`
	Country       string    `json:"country"`
	Network       string    `json:"network"`
	EffectiveDate time.Time `json:"effectiveDate"`
	Rate          string    `json:"rate"`
	ChangeType    string    `json:"changeType"`
}

type FormattedPriceList struct {
	Rows []PriceListRow `json:"rows"`
}
