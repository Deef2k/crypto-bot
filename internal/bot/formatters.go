package bot

import (
	"fmt"

	"github.com/Deef2k/crypto-bot/internal/models"
)

func FormatRates(rates []models.RateResponse, headerText string, time string) string {
	text := fmt.Sprintf("%s%s", headerText, time)
	for _, rate := range rates {
		text += fmt.Sprintf("Symbol:%s,\n Price:%s,\n LowPrice24h:%s,\n HighPrice24:%s,\n PriceChangePercent24h:%s,\n PriceChangePercent1h:%s\n\n",
			rate.Symbol, rate.Price, rate.LowPrice24, rate.HighPrice24,
			rate.PriceChangePercent24h, rate.PriceChangePercent1h)
	}
	return text
}
