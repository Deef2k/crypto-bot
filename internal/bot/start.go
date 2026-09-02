package bot

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func startBot(b *Bot, update *botapi.Update) {
	text := "Бот зупущен!\nСписок команд:\n/rates - курсы всех отслеживаемых валют\n/rates ...(BTCUSDT) - курс конкретной валюты(например: BTCUSDT)\n/start_auto ...(5) - запуск авторассылки курсов валют (например: на 5 минут)\n/stop_auto - остановка авторассылки курсов валют,\n/add_symbol ... (BTCUSDT) -добавить валюту в отслеживаемые (например:BTCUSDT)\n/remove_symbol ... (BTCUSDT) - удалить валюту из отслеживаемых (например:BTCUSDT)"
	msg := botapi.NewMessage(update.Message.Chat.ID, text)
	b.bot.Send(msg)
}
