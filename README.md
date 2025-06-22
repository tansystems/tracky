# Tracky — Telegram Mini App для отслеживания посылок

## Структура
- `/cmd/server` — запуск Go-сервера
- `/internal` — бизнес-логика, обработчики, интеграции
- `/ui/webapp` — фронтенд для Telegram Mini App

## Запуск
1. `go run cmd/server/main.go`
2. Настройте переменные окружения в `.env`

## Описание
Мини-приложение для отслеживания посылок через Telegram с интеграцией AfterShip. 

## Миграции БД
1. Установите [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate):
   
   brew install golang-migrate

2. Примените миграции:
   
   migrate -database "$DB_DSN" -path ./migrations up 