package bot

import (
	"log/slog"
	"strings"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func addSymbol(b *Bot, update *botapi.Update) {
	b.setPending(update.Message.Chat.ID, pendingAdd)
	b.bot.Send(botapi.NewMessage(update.Message.Chat.ID, "Укажите название пары валют"))
}

func applyAddSymbol(b *Bot, chatID int64, symbol string) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))

	if symbol == "" {
		b.bot.Send(botapi.NewMessage(chatID, "Укажите название пары валют"))
		return
	}

	_, err := fetchRate(b, symbol)
	if err != nil {
		slog.Warn("Ошибка в получении курса валют с API -", "err", err)
		b.bot.Send(botapi.NewMessage(chatID, "Произошла ошибка в получении курса(проверьте на правильность названия криптовалюты)"))
		return
	}

	if err := b.repo.AddTrackedSymbol(b.ctx, symbol); err != nil {
		slog.Warn("Ошибка в добавлении символа в БД", "err", err)
		b.bot.Send(botapi.NewMessage(chatID, "Произошла ошибка (Повторите попытку попозже)"))
		return
	}
	b.bot.Send(botapi.NewMessage(chatID, "Пара успешно добавлена в список отслеживаемых"))
}
