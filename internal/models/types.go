package models //Общая структура ответа RateResponse (DTO) для API/бота/БД

type RateResponse struct { //описываем как будут выглядить данные которые мы отправим клиенту
	Symbol                string `json:"symbol"` //заполнение полей - слева- их название в структуре,а справа значение которое в это поле кладём
	Price                 string `json:"price"`
	LowPrice24            string `json:"lowPrice"`
	HighPrice24           string `json:"highPrice"`
	PriceChangePercent24h string `json:"priceChangePercent24h"`
	PriceChangePercent1h  string `json:"priceChangePercent1h"`
} //вывел в отдельный файл так как нужен в нескольких файлах
