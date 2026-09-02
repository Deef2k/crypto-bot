package database

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTable(pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	table := `CREATE TABLE IF NOT EXISTS rates (
		id SERIAL PRIMARY KEY,
		symbol VARCHAR(20) NOT NULL,
		price DECIMAL(20, 8) NOT NULL,
		low_price DECIMAL(20, 8) NOT NULL,
		high_price DECIMAL(20, 8) NOT NULL,
		price_change_percent_24h TEXT NOT NULL,
		price_change_percent_1h TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS tracked_symbol (
		symbol VARCHAR(20) PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	INSERT INTO tracked_symbol (symbol) VALUES ('BTCUSDT'), ('ETHUSDT')
	ON CONFLICT (symbol) DO NOTHING;`

	_, err := pool.Exec(ctx, table)
	if err != nil {
		slog.Warn("Ошибка при создании таблицы:", "err", err) //сначало создаем ошибку
		return err                                            // затем её возвращаем
	}

	return nil
}
