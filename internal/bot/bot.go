package bot

import (
	"context"
	"log/slog"
	"os"
	"sync"

	"github.com/Deef2k/crypto-bot/storage"
	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SubscriptionManager struct {
	subscription map[int64]context.CancelFunc
	mutex        sync.Mutex
}
type Bot struct {
	bot     *botapi.BotAPI
	ctx     context.Context
	pointer *SubscriptionManager
	repo    storage.Repository
}

func Start(ctx context.Context, repo storage.Repository) {
	token := os.Getenv("TG_TOKEN")
	bot, err := botapi.NewBotAPI(token) //создание бота
	if err != nil {
		slog.Error("Бот не создался -", "err", err)
		return
	}
	subscriber := SubscriptionManager{
		subscription: make(map[int64]context.CancelFunc), //int64 - ключ User Телеграм,context -функция которая остновит авто рассылку
		mutex:        sync.Mutex{},
	}

	myBot := Bot{
		bot:     bot,
		ctx:     ctx,
		pointer: &subscriber,
		repo:    repo,
	}

	slog.Info("Бот создан")
	updateConfig := botapi.NewUpdate(0) //настройки канала
	updateConfig.Timeout = 15

	updates := bot.GetUpdatesChan(updateConfig) // канал с настройками

	for update := range updates { //update - принимает всю информацию (кто оправил,текс сообщения,команда,арументы команды и тп)
		if update.Message == nil {
			continue
		}

		if update.Message.Command() == "start" {
			startBot(&myBot, &update)
			continue
		}

		if update.Message.Command() == "rates" {
			separation(&myBot, &update)
			continue
		}

		if update.Message.Command() == "start_auto" {
			startAuto(&myBot, &update)
			continue
		}

		if update.Message.Command() == "stop_auto" {
			stopAuto(&myBot, &update)
			continue
		}

		if update.Message.Command() == "remove_symbol" {
			removeSymbol(&myBot, &update)
			continue
		}
		if update.Message.Command() == "add_symbol" {
			addSymbol(&myBot, &update)
			continue
		}
	}
}
