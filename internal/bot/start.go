package bot

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func startBot(b *Bot, update *botapi.Update) {
	text := "Бот запущен!\nСписок команд:\n/rates - затем пару (BTCUSDT) или all для всех отслеживаемых курсов\n/start_auto ...(5) - запуск авторассылки курсов валют (например: на 5 минут)\n/stop_auto - остановка авторассылки курсов валют\n/add_symbol - затем пару (например BTCUSDT)\n/remove_symbol - затем пару (например BTCUSDT)"
	msg := botapi.NewMessage(update.Message.Chat.ID, text)
	b.bot.Send(msg)
}
