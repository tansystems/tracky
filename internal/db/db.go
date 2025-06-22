package db

import (
	"database/sql"

	"tracky/internal/models"

	_ "github.com/lib/pq"
)

func InitDB(dataSourceName string) (*sql.DB, error) {
	return sql.Open("postgres", dataSourceName)
}

func AddTracking(db *sql.DB, t *models.Tracking) error {
	query := `INSERT INTO trackings (user_id, tracking_number, carrier_code, status, last_update) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	return db.QueryRow(query, t.UserID, t.TrackingNumber, t.CarrierCode, t.Status, t.LastUpdate).Scan(&t.ID, &t.CreatedAt)
}

func GetTracking(db *sql.DB, id int) (*models.Tracking, error) {
	var t models.Tracking
	query := `SELECT id, user_id, tracking_number, carrier_code, status, last_update, created_at FROM trackings WHERE id=$1`
	row := db.QueryRow(query, id)
	if err := row.Scan(&t.ID, &t.UserID, &t.TrackingNumber, &t.CarrierCode, &t.Status, &t.LastUpdate, &t.CreatedAt); err != nil {
		return nil, err
	}
	return &t, nil
}

func ListTrackings(db *sql.DB, userID int) ([]models.Tracking, error) {
	query := `SELECT id, user_id, tracking_number, carrier_code, status, last_update, created_at FROM trackings WHERE user_id=$1`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Tracking
	for rows.Next() {
		var t models.Tracking
		if err := rows.Scan(&t.ID, &t.UserID, &t.TrackingNumber, &t.CarrierCode, &t.Status, &t.LastUpdate, &t.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func DeleteTracking(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM trackings WHERE id=$1`, id)
	return err
}

func UpdateTrackingStatus(db *sql.DB, id int, status string) error {
	_, err := db.Exec(`UPDATE trackings SET status=$1, last_update=NOW() WHERE id=$2`, status, id)
	return err
}

func InsertOrUpdateUser(db *sql.DB, telegramID int64, username string) error {
	_, err := db.Exec(`INSERT INTO users (telegram_id, username) VALUES ($1, $2) ON CONFLICT (telegram_id) DO UPDATE SET username=EXCLUDED.username`, telegramID, username)
	return err
}
