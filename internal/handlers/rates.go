package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Deef2k/crypto-bot/internal/api"
	"github.com/Deef2k/crypto-bot/internal/models"
)

type RatesRepository interface {
	GetAllRates(ctx context.Context) ([]models.RateResponse, error)
	GetLastRate(ctx context.Context, symbol string) (models.RateResponse, error)
	SaveInfo(ctx context.Context, rate models.RateResponse) error
}

func GetRateHandler(repo RatesRepository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) { //(Функция замыкание)пишем внутри функции потому что при использовании этой функции полученные данные запишуться в w и r
		symbol := r.PathValue("symbol") //так как у нас в binance.go указан symbol в http запросе то что-бы его взять нам нужно
		var rateData models.RateResponse
		rateData, err := repo.GetLastRate(r.Context(), symbol)
		var BDerr error
		if err != nil {
			slog.Warn("Не удалось найти данные из бд по валюте,запрашиваем данные из Binance", "symbol", symbol, "err", err)
			rateData, err = api.GetRate(symbol)
			if err != nil {
				slog.Warn("Ошибка в получении данных с сайта при запросе от пользователя", "symbol", symbol, "err", err)
				http.Error(w, "Проблема со сторонним провайдером", http.StatusBadGateway)
				return
			}
			if BDerr = repo.SaveInfo(r.Context(), rateData); BDerr != nil {
				slog.Error("Ошибка в сохранение данных в БД", "err", BDerr)
			}
		}
		//r.Context-контекст который при выключении отдаст команду другим процессам (не выполняться),symbol- вызываем функцию с этим символом,conn- коннект который используем
		w.Header().Set("Content-type", "application/json") //говорим браузеру что мы отдаем JSON
		json.NewEncoder(w).Encode(rateData)                //куда писать, и каким методом расшифровывать
	} //w-инструмент для ответа, r - данные запроса
}

func GetRateAllHandlers(rate RatesRepository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		rates, err := rate.GetAllRates(r.Context())
		if err != nil {
			slog.Warn("Проблема при запросе всех актуальных курсов", "err", err)
			http.Error(w, "Ошибка сервера ", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-type", "application/json") //говорим браузеру что мы отдаем JSON
		json.NewEncoder(w).Encode(rates)
	}
}
