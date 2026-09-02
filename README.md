# 🤖 Crypto Bot

Telegram-бот для отслеживания курсов криптовалют с Binance.

## ✨ Возможности

### Telegram-команды

| Команда                      | Описание                             | Пример                  |
|------------------------------|--------------------------------------|-------------------------|
| `/start`                     | Запуск бота + получение списка команд| `/start`                |
| `/rates`                     | Курсы всех отслеживаемых валют       | `/rates`                |
| `/rates <SYMBOL>`            | Курс конкретной пары                 | `/rates BTCUSDT`        |
| `/start_auto <минуты>`       | Авто-рассылка курсов                 | `/start_auto 10`        |
| `/stop_auto`                 | Остановка авто-рассылки              | `/stop_auto`            |
| `/add_symbol <SYMBOL>`       | Добавить валюту в отслеживаемые      | `/add_symbol BTCUSDT`   |
| `/remove_symbol <SYMBOL>`    | Убрать валюту из отслеживания        | `/remove_symbol BTCUSDT`|

### HTTP API 

| Метод | Endpoint              | Описание             | Пример ответа |
|-------|-----------------------|----------------------|---------------|
| GET   | `/rates`              | Все курсы валют      | JSON массив   |
| GET   | `/rates/{symbol}`     | Курс конкретной пары | JSON объект   |

#### Пример ответа `/rates/ADAUSDT`:
```json


  {
    "symbol": "ADAUSDT",
    "price": "0.18340000",
    "lowPrice": "0.18060000",
    "highPrice": "0.18670000",
    "priceChangePercent24h": "-0.380",
    "priceChangePercent1h": "-0.06 %"
  }

```
## Запуск через Docker:

  ### 1) git clone https://github.com/Deef2k/crypto-bot.git
  ### 2) cd crypto-bot
  ### 3) создайте и заполните файл .env
  ### 4) docker-compose up -d --build

## Локальный запуск(для разработки):

  ### 1) git clone https://github.com/Deef2k/crypto-bot.git
  ### 2) cd crypto-bot
  ### 3)Установить и запустить PosgreSQL
  ### 4)создайте и заполнить файл .env
  ### 5)go run main.go 

## Пример заполнения файла .env:
  ### .env.example  