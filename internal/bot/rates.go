package bot

import (
	"log/slog"
	"strings"
	"time"

	"github.com/Deef2k/crypto-bot/internal/api"
	"github.com/Deef2k/crypto-bot/internal/models"
	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func rates(b *Bot, update *botapi.Update) {
	b.setPending(update.Message.Chat.ID, pendingRates)
	b.bot.Send(botapi.NewMessage(update.Message.Chat.ID, "Введите пару, например BTCUSDT. Чтобы увидеть все отслеживаемые курсы — напишите all"))
}

func applyRates(b *Bot, chatID int64, symbol string) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		b.bot.Send(botapi.NewMessage(chatID, "Укажите название пары валют или all"))
		return
	}
	if symbol == "ALL" || symbol == "ВСЕ" {
		showAllTrackedRates(b, chatID)
		return
	}
	showRateBySymbol(b, chatID, symbol)
}

func showAllTrackedRates(b *Bot, chatID int64) {
	times := time.Now().Format("2006-01-02 15:04:05\n")
	rates, err := b.repo.GetAllRates(b.ctx)
	if err != nil {
		slog.Warn("Не удалось получить курсы всех валют с бота -", "err", err)
		b.bot.Send(botapi.NewMessage(chatID, "Функция временно не работает"))
		return
	}
	text := FormatRates(rates, "Курсы валют на -", times)
	msg := botapi.NewMessage(chatID, text)
	b.bot.Send(msg)
}

func showRateBySymbol(b *Bot, chatID int64, symbol string) {
	times := time.Now().Format("2006-01-02 15:04:05\n")

	rate, err := b.repo.GetLastRate(b.ctx, symbol)
	if err != nil {
		slog.Info("такой пары нет в БД,идем в API")
		rate, err = fetchRate(b, symbol)
		if err != nil {
			slog.Warn("Ошибка в получении курса из API", "err", err)
			b.bot.Send(botapi.NewMessage(chatID, "Произошла ошибка (Убедитесь что ввели правильно название пары!)"))
			return
		}
	}
	text := FormatRates([]models.RateResponse{rate}, "Курсы валют на -", times)
	msg := botapi.NewMessage(chatID, text)
	b.bot.Send(msg)
}

func fetchRate(b *Bot, symbol string) (models.RateResponse, error) {
	rate, err := api.GetRate(symbol, "https://api.binance.com")
	if err != nil { //ошибка уже с API то возвращяем её юзеру
		slog.Warn("Ошибка в получении курса из API", "err", err)
		return models.RateResponse{}, err
	}
	if err = b.repo.SaveInfo(b.ctx, rate); err != nil { //В APi нашлось сохраняем в бд
		slog.Warn("Ошибка в сохранении данных в БД", "err", err)
	}
	return rate, nil
}
