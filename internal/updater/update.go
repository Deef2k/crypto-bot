package updater //обновление курсов

import (
	"context"
	"log/slog"
	"time"

	"github.com/Deef2k/crypto-bot/internal/api"
	"github.com/Deef2k/crypto-bot/storage"
)

func ValutResult(ctx context.Context, repo storage.Repository, symbols []string, baseURL string) { //вывел в отдельную функцию так как используеться в нескольких местах DRY
	for _, symbol := range symbols {
		result, err := api.GetRate(symbol, baseURL)
		if err != nil {
			slog.Warn("Ошибка в получении данных с сайта ", "err", err) //warn - потому что не критическая ошибка
			continue
		}
		if err = repo.SaveInfo(ctx, result); err != nil {
			slog.Error("Ошибка в сохранение данных в БД", "err", err)
			continue
		}
		slog.Info("Данные получены:", "symbol", symbol)
	}
}

func Update(ctx context.Context, repo storage.Repository) {

	symbols, err := repo.GetTrackedSymbols(ctx)
	if err != nil {
		slog.Error("Ошибка в получении символов из БД", "err", err)
		return // так как дальнейшее смысл выполнения безсмыслен
	}
	baseURL := "https://api.binance.com"

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	ValutResult(ctx, repo, symbols, baseURL)
	for {
		select {
		case <-ctx.Done():
			slog.Info("Получен сигнал,коректно завершаем работу програмы")
			return //ctx - это канал открытый только на чтение и при срабатывании cancel канал закрывается и select понимает что он отработал
		case <-ticker.C:
			symbols, err := repo.GetTrackedSymbols(ctx)
			if err != nil {
				slog.Error("Ошибка в получении символов из БД", "err", err)
				continue
			}
			ValutResult(ctx, repo, symbols, baseURL)
		}
	}
}
