package main

import (
	"encoding/json"
	"net/http"
	"threehead/db"
)

var donate = NewUserHandler().Ok(func(w http.ResponseWriter, r *http.Request, user *db.User) {
	tmpl.ExecuteTemplate(w, "donate", nil)
})

type AvatarPage struct {
	User        AvatarPageUser          `json:"user"`
	TopDonaters []AvatarPageTopDonaters `json:"topDonaters"`
	LatestHeads []string                `json:"latestHeads"`
}
type AvatarPageUser struct {
	DisplayName string `json:"displayName"`
}
type AvatarPageTopDonaters struct {
	DisplayName string `json:"displayName"`
	Amount      int    `json:"amount"`
}

var avatar = NewUserHandler().Ok(func(w http.ResponseWriter, r *http.Request, user *db.User) {
	if r.URL.Query().Get("json") == "true" {
		w.Header().Add("Content-Type", "application/json")
		j := json.NewEncoder(w)
		d := AvatarPage{
			User: AvatarPageUser{
				DisplayName: user.DisplayName,
			},
			TopDonaters: []AvatarPageTopDonaters{},
			LatestHeads: []string{},
		}
		for _, v := range db.GetTopDonaters() {
			d.TopDonaters = append(d.TopDonaters, AvatarPageTopDonaters{
				DisplayName: v.DisplayName,
				Amount:      int(v.FormattedAmount()),
			})
		}
		for _, v := range user.GetLastHeads(15, 0) {
			d.LatestHeads = append(d.LatestHeads, "/head/"+v.ID)
		}
		j.Encode(d)
		return
	}
	tmpl.ExecuteTemplate(w, "avatar", struct {
		User         *db.User
		TopDontaters []*db.Donater
	}{user, db.GetTopDonaters()})
})

var settings = NewUserHandler().Ok(func(w http.ResponseWriter, r *http.Request, user *db.User) {
	if r.URL.Query().Get("json") == "true" {
		w.Header().Add("Content-Type", "application/json")
		j := json.NewEncoder(w)
		j.Encode(struct {
			User struct {
				DisplayName           string `json:"displayName"`
				IsHiddenInTopDonaters bool   `json:"isHiddenInTopDonaters"`
			} `json:"user"`
		}{User: struct {
			DisplayName           string `json:"displayName"`
			IsHiddenInTopDonaters bool   `json:"isHiddenInTopDonaters"`
		}{
			DisplayName:           user.DisplayName,
			IsHiddenInTopDonaters: user.IsHiddenInTopDonaters,
		}})
		return
	}

	tmpl.ExecuteTemplate(w, "settings", struct {
		User *db.User
	}{user})
})

var create = NewUserHandler().NotOk(func(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "create", nil)
})

var index = NewUserHandler().NotOk(func(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "login", nil)
}).Ok(func(w http.ResponseWriter, r *http.Request, user *db.User) {
	http.Redirect(w, r, "/avatar", http.StatusTemporaryRedirect)
})
