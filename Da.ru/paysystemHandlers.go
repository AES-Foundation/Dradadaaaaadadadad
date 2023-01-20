package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"threehead/db"
	"threehead/payments/qiwi"
)

var actionPaymentQiwi = func(w http.ResponseWriter, r *http.Request) {
	notification := &qiwi.BillNotification{}
	err := json.NewDecoder(r.Body).Decode(notification)
	if err != nil {
		fmt.Fprintln(w, "Invalid request")
		return
	}
	payment := db.GetPaymentByID(notification.Bill.BillID)
	if payment == nil || notification.Bill.Status.Value != "PAID" {
		fmt.Fprintln(w, "Invalid id")
		return
	}
	payment.Pay()
	fmt.Fprintln(w, "ok")
}
