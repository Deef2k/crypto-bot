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
type pendingKind string

const (
	pendingAdd    pendingKind = "add"
	pendingRemove pendingKind = "remove"
	pendingRates  pendingKind = "rates"
)

type Bot struct {
	bot       *botapi.BotAPI
	ctx       context.Context
	pointer   *SubscriptionManager
	repo      storage.Repository
	pending   map[int64]pendingKind //int64 - ключ ChatID
	pendingMu sync.Mutex
}

func (b *Bot) setPending(chatID int64, kind pendingKind) {
	b.pendingMu.Lock()
	defer b.pendingMu.Unlock()
	b.pending[chatID] = kind
}

func (b *Bot) takePending(chatID int64) (pendingKind, bool) {
	b.pendingMu.Lock()
	defer b.pendingMu.Unlock()
	kind, ok := b.pending[chatID]
	if !ok {
		return "", false
	}
	delete(b.pending, chatID)
	return kind, true
}

func (b *Bot) clearPending(chatID int64) {
	b.pendingMu.Lock()
	defer b.pendingMu.Unlock()
	delete(b.pending, chatID)
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
		pending: make(map[int64]pendingKind),
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
			myBot.clearPending(update.Message.Chat.ID)
			startBot(&myBot, &update)
			continue
		}

		if update.Message.Command() == "rates" {
			rates(&myBot, &update)
			continue
		}

		if update.Message.Command() == "start_auto" {
			myBot.clearPending(update.Message.Chat.ID)
			startAuto(&myBot, &update)
			continue
		}

		if update.Message.Command() == "stop_auto" {
			myBot.clearPending(update.Message.Chat.ID)
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

		kind, ok := myBot.takePending(update.Message.Chat.ID)
		if !ok {
			continue
		}

		switch kind {
		case pendingAdd:
			applyAddSymbol(&myBot, update.Message.Chat.ID, update.Message.Text)
		case pendingRemove:
			applyRemoveSymbol(&myBot, update.Message.Chat.ID, update.Message.Text)
		case pendingRates:
			applyRates(&myBot, update.Message.Chat.ID, update.Message.Text)
		}
	}
}
