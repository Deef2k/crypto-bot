package database

import (
	"context"
	"log/slog"

	"github.com/Deef2k/crypto-bot/internal/models"
	"github.com/jackc/pgx/v5"
)

func GetLastRate(ctx context.Context, symbol string, conn *pgx.Conn) (models.RateResponse, error) {
	query := `Select symbol,price,low_price,high_price,price_change_percent
			From rates
			where symbol = $1
			Order by created_at desc
			Limit 1`
	var rate models.RateResponse

	if err := conn.QueryRow(ctx, query, symbol).Scan(&rate.Symbol, &rate.Price, &rate.LowPrice24, &rate.HighPrice24, &rate.PriceChangePercent); err != nil { //тут не нужен defer close потому что при выдачи одной строки он сразу сам закрывает
		slog.Warn("Ошибка в получении данных с БД при запросе от клиента", "err", err)
		return models.RateResponse{}, err //QueryRow-в ответе получаем только одну строку а Query- в ответе получаем несколько строк
	}
	return rate, nil
}
