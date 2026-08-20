package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Deef2k/crypto-bot/internal/models"
)

type BinanceResponce struct { //структура для хранения ответов от Binance по API
	Symbol             string `json:"symbol"` //с большой буквы что-бы было публичным
	Price              string `json:"lastPrice"`
	LowPrice24         string `json:"lowPrice"`
	HighPrice24        string `json:"highPrice"`
	PriceChangePercent string `json:"priceChangePercent"` //string потому что binance отдает как строку
}

func GetRate(symbol string) (models.RateResponse, error) { //функция получения курса криптовалют

	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/24hr?symbol=%s", symbol) // Sprintf-

	client := &http.Client{ // настраиваем http.Client
		Timeout: 1 * time.Second, //время выбрал как примерное время обработки маленького запроса
	}
	resp, err := client.Get(url) //делаем get запрос от клиента
	if err != nil {
		return models.RateResponse{}, fmt.Errorf("Ошибка в запросе от клиента по этому url - %s,ошибка - %v", url, err)
	}
	defer resp.Body.Close() //закрываем resp что бы не держать открытым

	if resp.StatusCode != http.StatusOK { //проверяем статус код
		return models.RateResponse{}, fmt.Errorf("Ошибка в загрузки страницы- %d", resp.StatusCode) // Errorf- создает ошибку
	}

	result, err := io.ReadAll(resp.Body) //записали данные полученные resp используя io.ReadAll
	if err != nil {
		return models.RateResponse{}, fmt.Errorf("Ошибка в данных полученных с сервера,ошибка %v", err)
	}
	var Crypto BinanceResponce
	err = json.Unmarshal(result, &Crypto) //переводим данные из байт в переменную типа структуры,err = ... -нужен что-бы если сервер вернет невалидные данные выдасться ошибка
	if err != nil {                       // используем & потому что без него мы использовали бы копию а так записываем в сам crypto
		return models.RateResponse{}, fmt.Errorf("ошибка в полученных данных (не тот формат)-%v", err)
	}

	return models.RateResponse{
		Symbol:             Crypto.Symbol,
		Price:              Crypto.Price,
		LowPrice24:         Crypto.LowPrice24,
		HighPrice24:        Crypto.HighPrice24,
		PriceChangePercent: Crypto.PriceChangePercent,
	}, nil
}
