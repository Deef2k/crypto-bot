package updater

import (
	"context"
	"log/slog"
	"time"

	"github.com/Deef2k/crypto-bot/internal/api"
	"github.com/Deef2k/crypto-bot/internal/handlers"
)

func Update(ctx context.Context, repo handlers.RatesRepository) {

	valute := []string{"BTCUSDT", "ETHUSDT", "SOLUSDT"}

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for _, valut := range valute { //нужно что-бы выдать курсы сразу при запуске а потом уйти на 5 минутное ожидание
		result, err := api.GetRate(valut)
		if err != nil {
			slog.Warn("Ошибка в получении данных с сайта ", "err", err) //warn - потому что не критическая ошибка
			continue
		}
		if err = repo.SaveInfo(ctx, result); err != nil {
			slog.Error("Ошибка в сохранение данных в БД", "err", err)
			continue
		}
		slog.Info("Данные получены:", "valute", valut, "result", result)
	}
	for {
		select {
		case <-ctx.Done():
			slog.Info("Получен сигнал,коректно завершаем работу програмы")
			return //ctx - это канал открытый только на чтение и при срабатывании cancel канал закрывается и select понимает что он отработал
		case <-ticker.C:
			for _, valut := range valute {
				result, err := api.GetRate(valut)
				if err != nil {
					slog.Warn("Ошибка в получении данных с сайта ", "err", err) //warn - потому что не критическая ошибка
					continue
				}
				if err = repo.SaveInfo(ctx, result); err != nil {
					slog.Error("Ошибка в сохранение данных в БД", "err", err)
					continue
				}
				slog.Info("Данные получены:", "valute", valut, "result", result)
			}
		}
	}
}
