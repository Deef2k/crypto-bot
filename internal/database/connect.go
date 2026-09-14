package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB() (*pgxpool.Pool, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		password := os.Getenv("DB_PASSWORD")
		port := os.Getenv("DB_PORT")
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		dbName := os.Getenv("DB_NAME")
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // создание контекста который прервет запрос к бд через 5 секуд либо при выходе из функции либо при выходе из функции
	defer cancel()                                                          //прервет внутренний таймер который будет занимать ресурсы
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		slog.Warn("Не удалось подключиться к БД:", "err", err)
		return nil, err
	}
	return pool, nil // лучше возвращать явно nil так как функция если дошла до этого места то значит все проверки прошли удачно
}
