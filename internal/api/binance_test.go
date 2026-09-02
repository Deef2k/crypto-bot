package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetChangePrice1h(t *testing.T) { //Пишим с большой буквы что бы test runer нашел эту функцию

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // создаем фейк сервер ( заместо banance)
		//w - для отправки ответов, r - полученный запрос
		response := `[ 
		[169000000,"29000","29200","28500","29000"],
		[169003000,"29100","29300","29050","29100"]
		]` //данные которые отдаст наш фейк сервер
		w.WriteHeader(http.StatusOK) //говорим клиенту всё ок = 200
		w.Write([]byte(response))    //отправляем json клиенту, []byte -превращаем строку в байт так как принимает только байты (слайс байт)
	}))
	defer server.Close()
	expected := "+0.34 %"
	symbol := "BTCUSDT"      //любое название
	client := &http.Client{} //стандартный http.client работает с mock
	baseURL := server.URL
	result, err := GetChangePrice1h(symbol, client, baseURL)
	if err != nil {
		t.Fatal("Функция не выполнилась правильно") //используем t. что бы ошибка выводилась в тестах а не пропадала при окончании его выполнения (если бы использовали fmt)
	}
	if result != expected {
		t.Errorf("Ожидали - %s, а получили - %s", expected, result)
	}
}

func TestGetChangePrice1h_Negative(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `[ 
		[169003000,"29100","29300","29050","29100"],
		[169000000,"29000","29200","28500","29000"]
		]` //данные которые отдаст наш фейк сервер
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()
	expected := "-0.34 %"
	symbol := "BTCUSDT"
	client := &http.Client{}
	baseURL := server.URL
	result, err := GetChangePrice1h(symbol, client, baseURL)
	if err != nil {
		t.Fatal("Функция не выполнилась правильно") //используем t. что бы ошибка выводилась в тестах а не пропадала при окончании его выполнения (если бы использовали fmt)
	}
	if result != expected {
		t.Errorf("Ожидали - %s, а получили - %s", expected, result)
	}
}
func TestGetRate_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //создание мок сервера

		if r.URL.Path == "/api/v3/klines" {
			response := `[
  [169000000, "29000", "29200", "28500", "29000"],
  [169003000, "29100", "29300", "29050", "29100"]
]`
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(response))

		} else if r.URL.Path == "/api/v3/ticker/24hr" {
			response := `{
  "symbol": "BTCUSDT",
  "lastPrice": "836453.23",
  "lowPrice": "834433.12",
  "highPrice": "900000.15",
  "priceChangePercent": "1.09"
}`
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(response))

		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	result, err := GetRate("BTCUSDT", server.URL) //делаем запрос на мок-сервер
	if err != nil {
		t.Error("Ошибка от мок-сервера")
	}
	if result.Symbol != "BTCUSDT" {
		t.Errorf("Ожидали - 'BTCUSDT',а получили -%s", result.Symbol)
	}

}
