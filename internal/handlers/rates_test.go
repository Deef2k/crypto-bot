package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Deef2k/crypto-bot/internal/models"
)

type FakeRepository struct {
	modRate []models.RateResponse

	getAllErr error //ошибка репозитория для всех курсов
}

func (f *FakeRepository) GetAllRates(ctx context.Context) ([]models.RateResponse, error) { //прописыаем метод для структуры что-бы структура подходила под требования фукнции куда будет отправляться
	if f.getAllErr != nil {
		return nil, f.getAllErr //если в тесте задать f,getAllErr - то fake имитирует падение бд
	}

	return f.modRate, nil //f *.. - этот метод принадлежит этой структуре   !!Не делал проверку так как пустой слайс как данные может обработаться невыдав ошибку!!
}
func (f *FakeRepository) GetLastRate(ctx context.Context, symbol string) (models.RateResponse, error) { //добавляем 2 метод что-бы стурктура FakeRepository полность соответсвовала RateResponse
	if len(f.modRate) > 0 { //проверяем есть ли хоть один элемент в текстовом файле
		return f.modRate[0], nil
	}
	return models.RateResponse{}, errors.New("no rates found")
}
func (f *FakeRepository) SaveInfo(ctx context.Context, rate models.RateResponse) error {
	return nil
}
func (f *FakeRepository) AddTrackedSymbol(ctx context.Context, symbol string) error {
	return nil
}
func (f *FakeRepository) RemoveTrackedSymbol(ctx context.Context, symbol string) error {
	return nil
}
func (f *FakeRepository) GetTrackedSymbols(ctx context.Context) ([]string, error) {
	return nil, nil
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
	handler := GetRateAllHandlers(&fakejson)         //создаем некую фукнцию которая уже работает с &fakejson
	req := httptest.NewRequest("GET", "/rates", nil) //Get-тип запроса,/rates-путь,nil-потому что для get запроса ничего тут не передаем
	rr := httptest.NewRecorder()                     //фейковая запись в блокнот(БД)

	handler(rr, req)
	if rr.Code != http.StatusOK { //проверяем не выдал ли запрос статус код ошибку
		t.Errorf("Ожидали статус-%d,получили статус - %d", http.StatusOK, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "BTCUSDT") {
		t.Errorf("Не тот ответ по price,ожидали - 'BTCUSDT',получили -%s", rr.Body.String())
	}
}
func TestGetRateHandler_OK(t *testing.T) {
	fake := &FakeRepository{
		modRate: []models.RateResponse{
			{
				Symbol: "BTCUSDT",
				Price:  "112.123",
			},
		},
	}
	handler := GetRateHandler(fake)
	req := httptest.NewRequest(http.MethodGet, "/rates", nil)
	req.SetPathValue("symbol", "BTCUSDT") // обязательно!
	rr := httptest.NewRecorder()
	handler(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rr.Code, http.StatusOK)
	}
	var result models.RateResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if result.Symbol != "BTCUSDT" || result.Price != "112.123" {
		t.Errorf("got %+v", result)
	}
}

func TestGetRateAllHandler_DBError(t *testing.T) {

	fake := &FakeRepository{
		getAllErr: errors.New("database error"),
	}
	handler := GetRateAllHandlers(fake)

	req := httptest.NewRequest(http.MethodGet, "/rates", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("оиждался статус - %d,а получили статус - %d", http.StatusInternalServerError, rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "Ошибка сервера") {
		t.Errorf("Ожидали ошибку сервера,а получили - %q", body)
	}
}
func TestGetRateAllHandlers_Empty(t *testing.T) {
	fake := &FakeRepository{
		modRate: []models.RateResponse{},
	}
	handler := GetRateAllHandlers(fake)
	req := httptest.NewRequest(http.MethodGet, "/rates", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Статус код ошибки не соответствует ожидаемому - %d", rr.Code)
	}
	body := rr.Body.String()
	if body != "[]" && body != "[]\n" {
		t.Errorf("Ожидали пустой ответ,а получили - %q", body)
	}
}
