package tracking

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type AfterShipTracking struct {
	Tracking struct {
		TrackingNumber string `json:"tracking_number"`
		CarrierCode    string `json:"carrier_code,omitempty"`
	} `json:"tracking"`
}

type AfterShipStatusResponse struct {
	Data struct {
		Tracking struct {
			Tag    string `json:"tag"`
			Status string `json:"subtag_message"`
		} `json:"tracking"`
	} `json:"data"`
}

type Carrier struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type carriersResponse struct {
	Data struct {
		Carriers []Carrier `json:"carriers"`
	} `json:"data"`
}

func AddTrackingAfterShip(trackingNumber, carrierCode string) error {
	apiKey := os.Getenv("AFTERSHIP_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("AFTERSHIP_API_KEY не задан")
	}
	url := "https://api.aftership.com/v4/trackings"
	body, _ := json.Marshal(AfterShipTracking{
		Tracking: struct {
			TrackingNumber string `json:"tracking_number"`
			CarrierCode    string `json:"carrier_code,omitempty"`
		}{
			TrackingNumber: trackingNumber,
			CarrierCode:    carrierCode,
		},
	})
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("aftership-api-key", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("AfterShip error: %s", string(b))
	}
	return nil
}

func GetTrackingStatusAfterShip(trackingNumber, carrierCode string) (string, error) {
	apiKey := os.Getenv("AFTERSHIP_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("AFTERSHIP_API_KEY не задан")
	}
	url := fmt.Sprintf("https://api.aftership.com/v4/trackings/%s/%s", carrierCode, trackingNumber)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("aftership-api-key", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("AfterShip error: %s", string(b))
	}
	var statusResp AfterShipStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return "", err
	}
	return statusResp.Data.Tracking.Status, nil
}

func GetCarriersFromAfterShip() ([]Carrier, error) {
	apiKey := os.Getenv("AFTERSHIP_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("AFTERSHIP_API_KEY не задан")
	}
	url := "https://api.aftership.com/v4/carriers"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("aftership-api-key", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("AfterShip error: %s", resp.Status)
	}
	var cr carriersResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return nil, err
	}
	return cr.Data.Carriers, nil
}
