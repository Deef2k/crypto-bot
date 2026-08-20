package database

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
	query := `Select Distinct On (symbol)symbol , price , low_price, high_price, price_change_percent 
	From rates
	Order by symbol,created_at Desc` //Distinct on -берет последнюю запись с каждой группы, так как перед этим мы отсортировали их с помощью desc
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		slog.Warn("Запрос от клиента на получения всех курсов валют не сработал", "err", err)
		return nil, err
	}
	defer rows.Close() //делаем defer после проверки,так как если будет ошибка то rows будет nil и нечего закрывать
	for rows.Next() {
		var r models.RateResponse
		if err := rows.Scan(&r.Symbol, &r.Price, &r.LowPrice24, &r.HighPrice24, &r.PriceChangePercent); err != nil {
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
	query := `INSERT INTO rates (symbol, price, low_price, high_price, price_change_percent) 
              VALUES ($1, $2, $3, $4, $5)`
	_, err := r.pool.Exec(ctx, query, rate.Symbol, rate.Price, rate.LowPrice24, rate.HighPrice24, rate.PriceChangePercent)
	return err
}

func NewPostgresRateRepository(pool *pgxpool.Pool) *PostgresRateRepository { //функция конструктор который создает эту структуру
	return &PostgresRateRepository{
		pool: pool,
	}
}
func (r *PostgresRateRepository) GetLastRate(ctx context.Context, symbol string) (models.RateResponse, error) {
	query := `Select symbol,price,low_price,high_price,price_change_percent
			From rates
			where symbol = $1
			Order by created_at desc
			Limit 1`
	var rate models.RateResponse

	if err := r.pool.QueryRow(ctx, query, symbol).Scan(&rate.Symbol, &rate.Price, &rate.LowPrice24, &rate.HighPrice24, &rate.PriceChangePercent); err != nil { //тут не нужен defer close потому что при выдачи одной строки он сразу сам закрывает
		slog.Warn("Ошибка в получении данных с БД при запросе от клиента", "err", err)
		return models.RateResponse{}, err //QueryRow-в ответе получаем только одну строку а Query- в ответе получаем несколько строк
	}
	return rate, nil
}
