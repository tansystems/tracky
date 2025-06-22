package notifier

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func Notify(userID int, message string) error {
	// Получаем telegram_id пользователя
	dbConnStr := os.Getenv("DB_DSN")
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		return err
	}
	defer db.Close()
	var telegramID int64
	err = db.QueryRow("SELECT telegram_id FROM users WHERE id=$1", userID).Scan(&telegramID)
	if err != nil {
		return err
	}
	// Отправляем сообщение через Telegram Bot API
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN не задан")
	}
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	resp, err := http.PostForm(apiURL, url.Values{
		"chat_id": {fmt.Sprintf("%d", telegramID)},
		"text":    {message},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("Ошибка Telegram API: %s", resp.Status)
	}
	return nil
}
