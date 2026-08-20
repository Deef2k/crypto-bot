package bot

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Deef2k/crypto-bot/internal/api"
	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func allRates(b *Bot, update *botapi.Update) {
	var text string
	times := time.Now().Format("2006-01-02 15:04:05\n")
	if update.Message.CommandArguments() == "" {
		rates, err := b.repo.GetAllRates(b.ctx)
		if err != nil {
			slog.Warn("Не удалось получить курсы всех валют с бота -", "err", err)
			b.bot.Send(botapi.NewMessage(update.Message.Chat.ID, "Функция временно не работает")) //update.Message.Chat.ID-заменяет messageID
			return
		}
		text += fmt.Sprintf("Курсы валют на -%s", times)
		for _, rate := range rates {
			text += fmt.Sprintf("Symbol:%s,\n Price:%s,\n LowPrice24h:%s,\n HighPrice24:%s,\n PriceChangePercent24h:%s\n\n", rate.Symbol, rate.Price, rate.LowPrice24, rate.HighPrice24, rate.PriceChangePercent)
		}
		msg := botapi.NewMessage(update.Message.Chat.ID, text)
		b.bot.Send(msg)
		return
	} else {
		symbol := update.Message.CommandArguments()
		rate, err := b.repo.GetLastRate(b.ctx, symbol)
		if err != nil {
			slog.Info("Такой пары либо не существует либо её нет в БД-", "err", err)

			rate, err = api.GetRate(symbol)
			if err != nil {
				slog.Warn("Ошибка в получении курса.Ошибка -", "err", err)
				b.bot.Send(botapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка в получение курса(Убедитесь что ввели правильно название валютной пары!)"))
				return
			} else {
				if err = b.repo.SaveInfo(b.ctx, rate); err != nil {
					slog.Warn("Ошибка в сохранении данных в БД.Ошибка -", "err", err)
					return
				}
			}
		}
		text += fmt.Sprintf("Курсы валют на данный промежуток времени -%s", times)
		text += fmt.Sprintf("Symbol:%s,\n Price:%s,\n LowPrice24h:%s,\n HighPrice24:%s,\n PriceChangePercent24h:%s\n\n", rate.Symbol, rate.Price, rate.LowPrice24, rate.HighPrice24, rate.PriceChangePercent)
		msg := botapi.NewMessage(update.Message.Chat.ID, text)
		b.bot.Send(msg)
	}
}
