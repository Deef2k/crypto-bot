package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Deef2k/crypto-bot/internal/models"
)

type FakeRepository struct {
	modRate []models.RateResponse
}

func (f *FakeRepository) GetAllRates(ctx context.Context) ([]models.RateResponse, error) { //прописыаем метод для структуры что-бы структура подходила под требования фукнции куда будет отправляться
	return f.modRate, nil //f *.. - этот метод принадлежит этой структуре   !!Не делал проверку так как пустой слайс как данные может обработаться невыдав ошибку!!
}
func (f *FakeRepository) GetLastRate(ctx context.Context, symbol string) (models.RateResponse, error) { //добавляем 2 метод что-бы стурктура FakeRepository полность соответсвовала RateResponse
	if len(f.modRate) > 0 { //проверяем есть ли хоть один элемент в текстовом файле
		return f.modRate[0], nil
	}
	return models.RateResponse{}, nil
}
func (f *FakeRepository) SaveInfo(ctx context.Context, rate models.RateResponse) error {
	return nil
}
func TestGetRateAllHandlers(t *testing.T) { //тестовые функции ничего не возвращают

	fakejson := FakeRepository{
		modRate: []models.RateResponse{ //поле с тестовыми данными
			{
				Symbol: "BTCUSDT",
				Price:  "112.123",
			},
		},
	}
	handler := GetRateAllHandlers(&fakejson)             //создаем некую фукнцию которая уже работает с &fakejson
	req := httptest.NewRequest("Get", "/api/rates", nil) //Get-тип запроса,/api/rates-путь,nil-потому что для get запроса ничего тут не передаем
	rr := httptest.NewRecorder()                         //фейковая запись в блокнот(БД)

	handler(rr, req)
	if rr.Code != http.StatusOK { //проверяем не выдал ли запрос статус код ошибку
		t.Errorf("Ожидали статус-%d,получили статус - %d", http.StatusOK, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "BTCUSDT") {
		t.Errorf("Не тот ответ по price,ожидали - 'BTCUSDT',получили -%s", rr.Body.String())
	}
}
func TestGetRateHandlers(t *testing.T) {
	fakejson := FakeRepository{
		modRate: []models.RateResponse{
			{
				Symbol: "BTCUSDT",
				Price:  "112.123",
			},
		},
	}
	handler := GetRateHandler(&fakejson)
	req := httptest.NewRequest("Get", "/api/rates/BTCUSDT", nil)
	rr := httptest.NewRecorder() //обьект ловит все что хендлер пытаеться отправть клиенту

	handler(rr, req)
	if rr.Code != http.StatusOK { //проверяем не выдал ли запрос статус код ошибку
		t.Errorf("Ожидали статус-%d,получили статус - %d", http.StatusOK, rr.Code)
	}
	var result models.RateResponse
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		t.Errorf("Был получен невалидный json-файл,ошибка - %v", err)
	}
	if result.Symbol != "BTCUSDT" {
		t.Errorf("Ожидали увидеть название - 'BTCUSDT',а получили - %s", result.Symbol)
	}
	if result.Price != "112.123" {
		t.Errorf("Ожидали курс - '112.123',получил-%s", result.Price)
	}
}
