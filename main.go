package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"        //работа с os
	"os/signal" //перехват сигналов
	"syscall"   //перехват сигналов

	"github.com/Deef2k/crypto-bot/internal/bot"
	"github.com/Deef2k/crypto-bot/internal/database"
	"github.com/Deef2k/crypto-bot/internal/handlers"
	"github.com/Deef2k/crypto-bot/internal/updater"
	"github.com/joho/godotenv"

	_ "github.com/Deef2k/crypto-bot/docs"        // ← импорт docks, _ - импорт без использования напрямую (для функции init)
	httpSwagger "github.com/swaggo/http-swagger" // ← Swagger UI
)

// @title Crypto Bot API
// @version 1.0
// @description API для получения курсов криптовалют с Binance
// @description
// @description Этот API предоставляет доступ к актуальным курсам криптовалют.
// @description Данные обновляются каждые 5 минут.

// @contact.name Deef2k
// @contact.url https://github.com/Deef2k

// @host localhost:8080
// @BasePath /

func main() {
	stop := make(chan os.Signal, 1)                         //создаем канал для ловли системных вызовов
	signal.Notify(stop, syscall.SIGINT)                     //обьясняет что при получении SIGINT-Ctrl + с , отправь в канал stop\
	ctx, cancel := context.WithCancel(context.Background()) //ctx-Новый контекст который можно отменить,cancel-функция которая отменяет контекст
	go func() {
		<-stop            //так как нам важно само сабытие то мы просто читаем из канала но не записываем а понимаем что там что-то есть
		signal.Stop(stop) //останавливает получение сигналов в stop
		cancel()
	}()

	if err := godotenv.Load(); err != nil {
		slog.Warn("Ошибка в чтение файла env", "err", err)
	} //читает все фалы .env и записывает их в память программы

	pool, err := database.ConnectDB()
	if err != nil {
		slog.Error("Проблема в подключении к бд", "err", err)
		return
	}
	defer pool.Close()
	// в pgx v5 нужно передавать контекст вместе с закрытием, в v4 не нужно было

	if err := database.CreateTable(pool); err != nil {
		slog.Warn("Таблица либо создана,либо ошибка в создании таблицы", "err", err)
	}

	repo := database.NewPostgresRateRepository(pool)

	allRateHandler := handlers.GetRateAllHandlers(repo)
	rateHandler := handlers.GetRateHandler(repo) //запускаем

	mux := http.NewServeMux()

	mux.Handle("/rates/{symbol}", rateHandler) //говорим что если постучаться с таким запросом то передавай его в rateHandler
	mux.Handle("/rates", allRateHandler)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	go func() { //запускаем горутину
		err := http.ListenAndServe(":8080", mux)
		if err != nil {
			slog.Error("Проблема в запуске http сервера.Ошибка:", "err", err)
		}
	}()

	go bot.Start(ctx, repo)

	go func() {
		updater.Update(ctx, repo)
	}()

	select {
	case <-ctx.Done():
		return
	}
}
