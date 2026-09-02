package bot

import (
	"fmt"
	"log/slog"
	"strings"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func removeSymbol(b *Bot, update *botapi.Update) {

	symbol := strings.ToUpper(strings.TrimSpace(update.Message.CommandArguments())) //toUpper- нужен что-бы если передали название в нижнем регистре он поднял до верхнео
	if symbol == "" {
		b.bot.Send(botapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, введите название валюты"))
		return
	}

	err := b.repo.RemoveTrackedSymbol(b.ctx, symbol)
	if err != nil {
		slog.Error("Ошибка в удалении символа из БД", "err", err)
		b.bot.Send(botapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при удалении символа из БД"))
		return
	}
	text := fmt.Sprintf("Символ %s успешно удален из отслеживаемых валют", symbol)
	b.bot.Send(botapi.NewMessage(update.Message.Chat.ID, text))
}
