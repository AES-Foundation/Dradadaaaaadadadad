package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"threehead/db"
	"threehead/payments/qiwi"
	"time"
)

var actionDonate = NewUserHandler().Ok(func(w http.ResponseWriter, r *http.Request, user *db.User) {
	r.ParseForm()
	if !r.Form.Has("amount") || !r.Form.Has("paysystem") {
		http.Redirect(w, r, "/donate?error=invalid_request", http.StatusTemporaryRedirect)
		return
	}
	amount, err := strconv.ParseFloat(r.Form.Get("amount"), 64)
	if err != nil || amount < 0 || amount > 10000 {
		http.Redirect(w, r, "/donate?error=invalid_request", http.StatusTemporaryRedirect)
		return
	}
	paysystem := r.Form.Get("paysystem")
	if paysystem != "qiwi" {
		http.Redirect(w, r, "/donate?error=invalid_paysystem", http.StatusTemporaryRedirect)
		return
	}
	payment := db.NewPayment(user, amount, paysystem)
	bill, err := qiwi.NewBillRequest(amount, "Донат").Send(payment.ID)
	if err != nil {
		http.Redirect(w, r, "/donate?error=invalid_paysystem", http.StatusTemporaryRedirect)
	}
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintln(w, "<script>window.location.href = \""+bill.PayURL+"\";</script>")
})

var actionGenerate = NewUserHandler().Ok(func(w http.ResponseWriter, r *http.Request, user *db.User) {
	w.Header().Add("Content-Type", "application/json")
	data, err := base64.StdEncoding.DecodeString(r.URL.Query().Get("data"))
	params := r.URL.Query().Get("params")
	if err != nil || params == "" {
		json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{"invalid request"})
		return
	}
	if params != "default" && params != "voxel" {
		json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{"invalid params"})
		return
	}
	skin, head, err := AddRequest(string(data), params)
	if err != nil {
		json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{err.Error()})
		return
	}
	h := db.NewHead(user, skin, head, params)
	json.NewEncoder(w).Encode(struct {
		Url string `json:"url"`
	}{
		"/head/" + h.ID,
	})
})

var actionCreate = func(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !r.Form.Has("email") || !r.Form.Has("password") || !r.Form.Has("repeat-password") {
		http.Redirect(w, r, "/create?error=invalid_request", http.StatusTemporaryRedirect)
		return
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	repeatPassword := r.Form.Get("repeat-password")
	if len(email) == 0 || len(password) < 8 || password != repeatPassword {
		http.Redirect(w, r, "/create?error=invalid_request", http.StatusTemporaryRedirect)
		return
	}
	if db.GetUserByEmail(email) != nil {
		http.Redirect(w, r, "/create?error=email_is_used", http.StatusTemporaryRedirect)
		return
	}
	db.NewUser(email, password)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

var actionLogin = func(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !r.Form.Has("email") || !r.Form.Has("password") {
		http.Redirect(w, r, "/?error=invalid_request", http.StatusTemporaryRedirect)
		return
	}
	user, key := db.Login(r.Form.Get("email"), r.Form.Get("password"))
	if user == nil {
		http.Redirect(w, r, "/?error=invalid_creditinals", http.StatusTemporaryRedirect)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "key",
		Value:   key,
		Path:    "/",
		Expires: time.Now().Add(time.Hour * 24 * 30),
	})
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

var actionSettings = NewUserHandler().Ok(func(w http.ResponseWriter, r *http.Request, user *db.User) {
	if r.URL.Query().Has("hideInTopDonaters") {
		hideInTopDonaters := false
		val := r.URL.Query().Get("hideInTopDonaters")
		if val != "true" && val != "false" {
			json.NewEncoder(w).Encode(struct {
				Error string `json:"error"`
			}{"invalid request"})
			return
		}
		if val == "true" {
			hideInTopDonaters = true
		}
		user.SetHideIntTopDonaters(hideInTopDonaters)

		return
	}
	if r.URL.Query().Has("displayName") {
		displayName := r.URL.Query().Get("displayName")
		if len(displayName) == 0 {
			json.NewEncoder(w).Encode(struct {
				Error string `json:"error"`
			}{"invalid request"})
			return
		}
		user.SetDisplayName(displayName)
		return
	}
	json.NewEncoder(w).Encode(struct {
		Ok bool `json:"bool"`
	}{
		true,
	})
})

var actionSettingsClearHistory = NewUserHandler().Ok(func(w http.ResponseWriter, r *http.Request, user *db.User) {
	user.ClearHeadsHistory()

	json.NewEncoder(w).Encode(struct {
		Ok bool `json:"bool"`
	}{
		true,
	})
})
