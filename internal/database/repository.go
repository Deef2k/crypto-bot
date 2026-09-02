package database //Реализация storage.Repository на SQL

import (
	"context"
	"log/slog"

	"github.com/Deef2k/crypto-bot/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRateRepository struct {
	pool *pgxpool.Pool
}

func (r *PostgresRateRepository) GetAllRates(ctx context.Context) ([]models.RateResponse, error) { //r - reciever-получатель- говорит что этот метод именно этой структуры
	var rates []models.RateResponse
	query := `SELECT DISTINCT ON (r.symbol)
    r.symbol, r.price, r.low_price, r.high_price,
    r.price_change_percent_24h, r.price_change_percent_1h
FROM rates r
INNER JOIN tracked_symbol t ON r.symbol = t.symbol
ORDER BY r.symbol, r.created_at DESC ` //INNER JOIN - объединяет таблицы по условию, так как у нас в таблице rates есть символы которые мы не хотим видеть то мы объединяем их с таблицей tracked_symbol и фильтруем по условию
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		slog.Warn("Запрос от клиента на получения всех курсов валют не сработал", "err", err)
		return nil, err
	}
	defer rows.Close() //делаем defer после проверки,так как если будет ошибка то rows будет nil и нечего закрывать

	for rows.Next() {
		var r models.RateResponse
		if err := rows.Scan(&r.Symbol, &r.Price, &r.LowPrice24, &r.HighPrice24, &r.PriceChangePercent24h, &r.PriceChangePercent1h); err != nil {
			slog.Warn("Ошикба в читке с бд", "err", err)
			return nil, err
		}
		rates = append(rates, r)
	}
	if err = rows.Err(); err != nil { //если будет ошибка то не отдает недозаполненный файл а выдает ошибку
		return nil, err
	}

	return rates, nil //отправляем nil потому что код если дошел до сюда то отработал без ошибок
}

//это мы написали ту самую реализацию GetAllRates(интерфеса) с помощью структуры PostgresRateRepository

func (r *PostgresRateRepository) SaveInfo(ctx context.Context, rate models.RateResponse) error {
	query := `INSERT INTO rates (symbol, price, low_price, high_price, price_change_percent_24h,price_change_percent_1h) 
              VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, query, rate.Symbol, rate.Price, rate.LowPrice24, rate.HighPrice24, rate.PriceChangePercent24h, rate.PriceChangePercent1h)
	return err
}

func NewPostgresRateRepository(pool *pgxpool.Pool) *PostgresRateRepository { //функция конструктор который создает эту структуру
	return &PostgresRateRepository{
		pool: pool,
	}
}
func (r *PostgresRateRepository) GetLastRate(ctx context.Context, symbol string) (models.RateResponse, error) {
	query := `Select symbol,price,low_price,high_price,price_change_percent_24h,price_change_percent_1h
			From rates
			where symbol = $1
			Order by created_at desc
			Limit 1`
	var rate models.RateResponse

	if err := r.pool.QueryRow(ctx, query, symbol).Scan(&rate.Symbol, &rate.Price, &rate.LowPrice24, &rate.HighPrice24, &rate.PriceChangePercent24h, &rate.PriceChangePercent1h); err != nil { //тут не нужен defer close потому что при выдачи одной строки он сразу сам закрывает
		slog.Warn("Ошибка в получении данных с БД при запросе от клиента", "err", err)
		return models.RateResponse{}, err //QueryRow-в ответе получаем только одну строку а Query- в ответе получаем несколько строк
	}
	return rate, nil
}

func (r *PostgresRateRepository) AddTrackedSymbol(ctx context.Context, symbol string) error {
	query := `INSERT INTO tracked_symbol (symbol) VALUES ($1)
			ON CONFLICT (symbol) DO NOTHING`
	_, err := r.pool.Exec(ctx, query, symbol)
	return err
}

func (r *PostgresRateRepository) RemoveTrackedSymbol(ctx context.Context, symbol string) error {
	query := `DELETE FROM tracked_symbol WHERE symbol = $1`
	_, err := r.pool.Exec(ctx, query, symbol)
	return err
}

func (r *PostgresRateRepository) GetTrackedSymbols(ctx context.Context) ([]string, error) {
	query := `SELECT symbol FROM tracked_symbol ORDER BY symbol`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var symbols []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		symbols = append(symbols, s)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return symbols, nil
}
