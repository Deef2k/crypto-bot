package bot

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func startAuto(b *Bot, update *botapi.Update) {
	if update.Message.CommandArguments() != "" {
		argTime := update.Message.CommandArguments()

		ctxTimer, cancel := context.WithCancel(b.ctx) //используеться для того что-бы остановить горутину(так как орутины нельзя останавливать с наружи) нужно иметь что-то внутри что остановит
		timeTick, err := strconv.Atoi(argTime)
		if err != nil {
			slog.Warn("Не удалось приобразовать строку в число.Ошибка - ", "err", err)
			b.bot.Send(botapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка в запросе(точно ли стоит число?)"))
			cancel()
			//не стоит добавлять Return так как из за маленького отключения от бд может выключиться авто обновление
		}
		if timeTick <= 0 {
			slog.Warn("Пользователь ввел число меньше 0 или 0.")
			b.bot.Send(botapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка в запросе(точно ли число выше или не равно 0?)"))
			cancel()
		} else {
			b.pointer.mutex.Lock()
			cancelFunc, exists := b.pointer.subscription[update.Message.Chat.ID]
			if exists {
				cancelFunc() //delete не нужен перезапись сама очистит старую запись
			}
			b.pointer.subscription[update.Message.Chat.ID] = cancel // записываем в неё cancel= стоп функция для ctxTimer
			b.pointer.mutex.Unlock()
			textStart := fmt.Sprintf("Рассылка сообщений запущена каждые -%s минут ", argTime)
			msg := botapi.NewMessage(update.Message.Chat.ID, textStart)
			b.bot.Send(msg)

			go func() {
				ticker := time.NewTicker(time.Duration(timeTick) * time.Minute)
				defer ticker.Stop()
				var msg botapi.MessageConfig //обьявил локально что-бы именно с ней работал 1 человек а следующий создовал для себя локальную переменную
				for {
					select {
					case <-ctxTimer.Done():
						return
					case <-ticker.C:
						times := time.Now().Format("2006-01-02 15:04:05\n")
						textTicker := fmt.Sprintf("(Самые важные)Курсы валют на данный промежуток времени -%s", times) //стоит тут что-бы вызывалось в начале каждого срабатывания кейса
						rates, err := b.repo.GetAllRates(b.ctx)
						if err != nil {
							slog.Warn("Не удалось получить курсы всех валют с бота -", "err", err)
							b.bot.Send(botapi.NewMessage(update.Message.Chat.ID, "Функция временно не работает"))
							return
						}
						for _, rate := range rates {
							textTicker += fmt.Sprintf("Symbol:%s,\n Price:%s,\n LowPrice24h:%s,\n HighPrice24:%s,\n PriceChangePercent24h:%s\n\n", rate.Symbol, rate.Price, rate.LowPrice24, rate.HighPrice24, rate.PriceChangePercent)
						}
						msg = botapi.NewMessage(update.Message.Chat.ID, textTicker)
						b.bot.Send(msg)
					}
				}
			}()
		}
	} else {
		slog.Warn("Пользователь не ввел число после start_auto")
		b.bot.Send(botapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка в запросе(не указано время после start_auto)"))
		return
	}
}
