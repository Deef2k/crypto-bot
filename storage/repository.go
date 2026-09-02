package storage

import (
	"context"

	"github.com/Deef2k/crypto-bot/internal/models"
)

type Repository interface {
	GetAllRates(ctx context.Context) ([]models.RateResponse, error)
	GetLastRate(ctx context.Context, symbol string) (models.RateResponse, error)
	SaveInfo(ctx context.Context, rate models.RateResponse) error

	AddTrackedSymbol(ctx context.Context, symbol string) error
	RemoveTrackedSymbol(ctx context.Context, symbol string) error
	GetTrackedSymbols(ctx context.Context) ([]string, error)
}
