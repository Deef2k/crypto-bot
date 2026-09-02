package bot

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func stopAuto(b *Bot, update *botapi.Update) {
	b.pointer.mutex.Lock()
	cancelFunc, exists := b.pointer.subscription[update.Message.Chat.ID] //будет true - если ключ есть
	b.pointer.mutex.Unlock()                                             //прочитали и сразу разблокировали
	if !exists {
		tickerNoHave := "Авто обновление не запущено."
		msg := botapi.NewMessage(update.Message.Chat.ID, tickerNoHave)
		b.bot.Send(msg)
	} else {
		cancelFunc()
		b.pointer.mutex.Lock()
		delete(b.pointer.subscription, update.Message.Chat.ID) //удаляем тикер
		b.pointer.mutex.Unlock()
		tickerStop := "Авто обновление успешно остановлено."
		msg := botapi.NewMessage(update.Message.Chat.ID, tickerStop)
		b.bot.Send(msg)
	}
}
