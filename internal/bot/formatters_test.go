package bot

import (
	"testing"

	"github.com/Deef2k/crypto-bot/internal/models"
)

func TestFormatRates(t *testing.T) {
	header := "Курсы валют на -"
	time := "2026-08-24\n"
	rates := []models.RateResponse{ //создаем слайс структур
		{
			Symbol:                "BTCUSDT",
			Price:                 "836453.23",
			LowPrice24:            "834433.12",
			HighPrice24:           "900000.15",
			PriceChangePercent24h: "+1.09 %",
			PriceChangePercent1h:  "+6.7 %",
		},
		{
			Symbol:                "SIXSEVEN",
			Price:                 "676767.21",
			LowPrice24:            "7676767",
			HighPrice24:           "122112.54",
			PriceChangePercent24h: "+1.11 %",
			PriceChangePercent1h:  "+6.7 %",
		},
	}
	result := FormatRates(rates, header, time)
	expected := `Курсы валют на -2026-08-24
Symbol:BTCUSDT,
 Price:836453.23,
 LowPrice24h:834433.12,
 HighPrice24:900000.15,
 PriceChangePercent24h:+1.09 %,
 PriceChangePercent1h:+6.7 %

Symbol:SIXSEVEN,
 Price:676767.21,
 LowPrice24h:7676767,
 HighPrice24:122112.54,
 PriceChangePercent24h:+1.11 %,
 PriceChangePercent1h:+6.7 %

`

	if result != expected {
		t.Errorf("Ожидалось - %q, Получили -%q", expected, result) //%s - Обычная строка, %q -строка с кавычками
	}
}
