package models

import "github.com/google/uuid"

type Candle struct {
	Open  float64 `json:"open"`
	Close float64 `json:"close"`
}

type ChartSegment struct {
	ID      uuid.UUID `json:"id"`
	Candles []Candle  `json:"candles"`
	Future  []Candle  `json:"future"`
}
