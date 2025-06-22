package models

type Tracking struct {
	ID             int
	UserID         int
	TrackingNumber string
	CarrierCode    string
	Status         string
	LastUpdate     string
	CreatedAt      string
}
