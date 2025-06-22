package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"tracky/internal/db"
)

var UserDB *sql.DB

func InitUserHandlers(dbConn *sql.DB) {
	UserDB = dbConn
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		TelegramID int64  `json:"telegram_id"`
		Username   string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.TelegramID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Добавляем пользователя, если его нет
	err := db.InsertOrUpdateUser(UserDB, req.TelegramID, req.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
