package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"tracky/internal/db"
	"tracky/internal/models"
	"tracky/internal/tracking"
)

var DBConn *sql.DB

func InitHandlers(dbConn *sql.DB) {
	DBConn = dbConn
}

func AddTrackingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		UserID         int    `json:"user_id"`
		TrackingNumber string `json:"tracking_number"`
		CarrierCode    string `json:"carrier_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	t := &models.Tracking{
		UserID:         req.UserID,
		TrackingNumber: req.TrackingNumber,
		CarrierCode:    req.CarrierCode,
		Status:         "new",
	}
	if err := db.AddTracking(DBConn, t); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if t.CarrierCode != "" {
		_ = tracking.AddTrackingAfterShip(t.TrackingNumber, t.CarrierCode)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func GetTrackingStatusHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	t, err := db.GetTracking(DBConn, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if t.CarrierCode != "" {
		status, err := tracking.GetTrackingStatusAfterShip(t.TrackingNumber, t.CarrierCode)
		if err == nil && status != "" && status != t.Status {
			t.Status = status
			_ = db.UpdateTrackingStatus(DBConn, t.ID, status)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func ListTrackingsHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tracks, err := db.ListTrackings(DBConn, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tracks)
}

func DeleteTrackingHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := db.DeleteTracking(DBConn, id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func CarriersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	carriers, err := tracking.GetCarriersFromAfterShip()
	if err != nil {
		http.Error(w, "Ошибка загрузки перевозчиков", http.StatusInternalServerError)
		return
	}
	var result []map[string]string
	for _, c := range carriers {
		result = append(result, map[string]string{
			"code": c.Slug,
			"name": c.Name,
		})
	}
	json.NewEncoder(w).Encode(result)
}
