package api //Клиент Binance: HTTP → JSON → RateResponse

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Deef2k/crypto-bot/internal/models"
)

type BinanceResponce struct { //структура для хранения ответов от Binance по API
	Symbol                string `json:"symbol"` //с большой буквы что-бы было публичным
	Price                 string `json:"lastPrice"`
	LowPrice24            string `json:"lowPrice"` //слева моё личное название, справо- тег от Binance что-бы было понятно какое значение от binance куда положить
	HighPrice24           string `json:"highPrice"`
	PriceChangePercent24h string `json:"priceChangePercent"` //string потому что binance отдает как строку
}

func GetChangePrice1h(symbol string, client *http.Client, baseURL string) (string, error) {
	url := fmt.Sprintf("%s/api/v3/klines?symbol=%s&interval=1h&limit=2", baseURL, symbol)
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("Ошибка в запросе от клиента по этому url - %s,ошибка - %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { //проверяем статус код
		return "", fmt.Errorf("Ошибка в загрузки страницы- %d", resp.StatusCode) // Errorf- создает ошибку
	}
	result, err := io.ReadAll(resp.Body) //записали данные полученные resp используя io.ReadAll
	if err != nil {
		return "", fmt.Errorf("Ошибка в данных полученных с сервера,ошибка %v", err)
	}
	var klines [][]interface{} //массив,где каждый элемент - массив из значений любого типа

	err = json.Unmarshal(result, &klines) // используем & потому что без него мы использовали бы копию а так записываем в сам crypto
	if err != nil {
		return "", fmt.Errorf("ошибка в полученных данных (не тот формат)-%v", err)
	}
	if len(klines) < 2 { // Проверка на то что сайт вернул результаты именно по 2 свечам
		return "", fmt.Errorf("Сайт вернул мало свечей")
	}

	price1, status1 := klines[0][4].(string) //получаем 4 значение(как раз цену за прошлый час)
	if !status1 {
		return "", fmt.Errorf("Ошибка в полученном типа переменной 1 (ожидалась строка)")
	}
	price2, status2 := klines[1][4].(string)
	if !status2 {
		return "", fmt.Errorf("Ошибка в полученном типа переменной 1 (ожидалась строка)")
	}

	pr1, err := strconv.ParseFloat(price1, 64)
	if err != nil {
		return "", fmt.Errorf("Не удалось конвертировать переменную - %s", price1) //Возвращаем price 1 - что-бы понять какая переменная не конвертировалась
	}

	pr2, err := strconv.ParseFloat(price2, 64)
	if err != nil {
		return "", fmt.Errorf("Не удалось конвертировать переменную - %s", price2)
	}

	Percent := ((pr2 - pr1) / pr1) * 100

	if Percent >= 0 {
		PriceChangePercent1h := fmt.Sprintf("%+.2f %%", Percent)
		return PriceChangePercent1h, nil
	}

	PriceChangePercent1h := fmt.Sprintf("%.2f %%", Percent)

	return PriceChangePercent1h, nil
}

func GetRate(symbol string, baseURL string) (models.RateResponse, error) { //функция получения курса криптовалют
	symbol = strings.ToUpper(symbol)
	var Crypto BinanceResponce
	client := &http.Client{ // настраиваем http.Client
		Timeout: 5 * time.Second, //время выбрал как примерное время обработки маленького запроса
	}
	PriceChangePercent1h, err := GetChangePrice1h(symbol, client, baseURL)
	if err != nil {
		return models.RateResponse{}, fmt.Errorf("Ошибка в счислении процентов изменения за час: %v", err)
	}

	url := fmt.Sprintf("%s/api/v3/ticker/24hr?symbol=%s", baseURL, symbol)

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
	err = json.Unmarshal(result, &Crypto) //переводим данные из байт в переменную типа структуры,err = ... -нужен что-бы если сервер вернет невалидные данные выдасться ошибка
	if err != nil {                       // используем & потому что без него мы использовали бы копию а так записываем в сам crypto
		return models.RateResponse{}, fmt.Errorf("ошибка в полученных данных (не тот формат)-%v", err)
	}
	i, err := strconv.ParseFloat(Crypto.PriceChangePercent24h, 64)
	if err != nil {
		return models.RateResponse{}, fmt.Errorf("Не удалось приобразовать PriceChangePercent24h в число")
	}
	if i >= 0 {
		Crypto.PriceChangePercent24h = fmt.Sprintf("%+.2f %%", i)
		return models.RateResponse{
			Symbol:                Crypto.Symbol,
			Price:                 Crypto.Price,
			LowPrice24:            Crypto.LowPrice24,
			HighPrice24:           Crypto.HighPrice24,
			PriceChangePercent24h: Crypto.PriceChangePercent24h,
			PriceChangePercent1h:  PriceChangePercent1h,
		}, nil
	}
	Crypto.PriceChangePercent24h = fmt.Sprintf("%.2f %%", i)
	return models.RateResponse{
		Symbol:                Crypto.Symbol,
		Price:                 Crypto.Price,
		LowPrice24:            Crypto.LowPrice24,
		HighPrice24:           Crypto.HighPrice24,
		PriceChangePercent24h: Crypto.PriceChangePercent24h,
		PriceChangePercent1h:  PriceChangePercent1h,
	}, nil
}
