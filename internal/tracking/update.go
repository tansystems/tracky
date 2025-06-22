package tracking

import (
	"database/sql"
	"fmt"
	"log"
	"tracky/internal/db"
	"tracky/internal/notifier"
)

func CheckUpdates(database *sql.DB) {
	// Получаем все треки с user_id
	rows, err := database.Query(`SELECT id, user_id, tracking_number, carrier_code, status FROM trackings`)
	if err != nil {
		log.Println("Ошибка получения треков:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id, userID int
		var trackingNumber, carrierCode, oldStatus string
		if err := rows.Scan(&id, &userID, &trackingNumber, &carrierCode, &oldStatus); err != nil {
			continue
		}
		if carrierCode == "" {
			continue
		}
		status, err := GetTrackingStatusAfterShip(trackingNumber, carrierCode)
		if err == nil && status != "" && status != oldStatus {
			_ = db.UpdateTrackingStatus(database, id, status)
			msg := fmt.Sprintf("Статус вашей посылки %s изменился: %s", trackingNumber, status)
			_ = notifier.Notify(userID, msg)
		}
	}
}
