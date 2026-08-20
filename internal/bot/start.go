package bot

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func startBot(b *Bot, update *botapi.Update) {
	msg := botapi.NewMessage(update.Message.Chat.ID, "Бот запущен") //создаем сообщение
	b.bot.Send(msg)                                                 //образаемся к полям переменно, отправляем созданное сообщение
}
