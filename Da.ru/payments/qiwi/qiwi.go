package qiwi

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"time"
)

const (
	SECRET_KEY = "HA"
	KEY        = "POPKA"
	BILL       = "https://api.qiwi.com/partner/bill/v1/bills/"
)

type BillRequest struct {
	Amount struct {
		Currency string `json:"currency"`
		Value    string `json:"value"`
	} `json:"amount"`
	Comment string `json:"comment"`
}

type BillResponse struct {
	SiteID string `json:"siteId"`
	BillID string `json:"billId"`
	Amount struct {
		Currency string `json:"currency"`
		Value    string `json:"value"`
	} `json:"amount"`
	Status struct {
		Value           string    `json:"value"`
		ChangedDateTime time.Time `json:"changedDateTime"`
	} `json:"status"`
	Customer struct {
		Phone   string `json:"phone"`
		Email   string `json:"email"`
		Account string `json:"account"`
	} `json:"customer"`
	CustomFields struct {
		PaySourcesFilter string `json:"paySourcesFilter"`
		ThemeCode        string `json:"themeCode"`
		YourParam1       string `json:"yourParam1"`
		YourParam2       string `json:"yourParam2"`
	} `json:"customFields"`
	Comment            string    `json:"comment"`
	CreationDateTime   time.Time `json:"creationDateTime"`
	ExpirationDateTime time.Time `json:"expirationDateTime"`
	PayURL             string    `json:"payUrl"`
}

type BillNotification struct {
	Bill struct {
		BillID string `json:"billId"`
		Status struct {
			Value           string `json:"value"`
			ChangedDateTime string `json:"changedDateTime"`
		} `json:"status"`
	} `json:"bill"`
	Version string `json:"version"`
}

func NewBillRequest(amount float64, comment string) *BillRequest {
	return &BillRequest{
		Amount: struct {
			Currency string "json:\"currency\""
			Value    string "json:\"value\""
		}{
			Currency: "RUB",
			Value:    strconv.FormatFloat(math.Floor(amount*100)/100.0, 'f', -1, 64),
		},
		Comment: comment,
	}
}

func (r *BillRequest) Send(id string) (*BillResponse, error) {
	data := bytes.NewBuffer(nil)
	json.NewEncoder(data).Encode(r)
	client := http.Client{}
	req, _ := http.NewRequest("PUT", BILL+id, data)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+SECRET_KEY)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	response := &BillResponse{}
	err = json.NewDecoder(res.Body).Decode(response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
