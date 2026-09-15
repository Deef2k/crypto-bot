package bot

import (
	"fmt"
	"log/slog"
	"strings"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func removeSymbol(b *Bot, update *botapi.Update) {
	b.setPending(update.Message.Chat.ID, pendingRemove)
	b.bot.Send(botapi.NewMessage(update.Message.Chat.ID, "Введите пару для удаления, например BTCUSDT"))
}

func applyRemoveSymbol(b *Bot, chatID int64, symbol string) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		b.bot.Send(botapi.NewMessage(chatID, "Пожалуйста, введите название валюты"))
		return
	}

	err := b.repo.RemoveTrackedSymbol(b.ctx, symbol)
	if err != nil {
		slog.Error("Ошибка в удалении символа из БД", "err", err)
		b.bot.Send(botapi.NewMessage(chatID, "Произошла ошибка при удалении символа из БД"))
		return
	}
	text := fmt.Sprintf("Символ %s успешно удален из отслеживаемых валют", symbol)
	b.bot.Send(botapi.NewMessage(chatID, text))
}
