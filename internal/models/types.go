package models

type RateResponse struct { //описываем как будут выглядить данные которые мы отправим клиенту
	Symbol             string `json:"symbol"`
	Price              string `json:"price"`
	LowPrice24         string `json:"lowPrice"`
	HighPrice24        string `json:"highPrice"`
	PriceChangePercent string `json:"priceChangePercent"`
} //вывел в отдельный файл так как нужен в нескольких файлах
