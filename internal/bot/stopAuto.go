package bot

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func stopAuto(b *Bot, update *botapi.Update) {
	b.pointer.mutex.Lock()
	cancelFunc, exists := b.pointer.subscription[update.Message.Chat.ID] //будет true - если ключ есть
	defer b.pointer.mutex.Unlock()
	if !exists {
		tickerNoHave := "Авто обновление не запущено."
		msg := botapi.NewMessage(update.Message.Chat.ID, tickerNoHave)
		b.bot.Send(msg)
	} else {
		cancelFunc()
		delete(b.pointer.subscription, update.Message.Chat.ID) //удаляем тикер
		tickerStop := "Авто обновление успешно остановлено."
		msg := botapi.NewMessage(update.Message.Chat.ID, tickerStop)
		b.bot.Send(msg)
	}
}
