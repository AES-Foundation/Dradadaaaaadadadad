package main

import (
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"threehead/db"
)

const (
	HOST = ":21098"
)

var tmpl, _ = template.ParseGlob("templates/*.html")

func main() {
	db.OpenDB()
	RunOrderThread()

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/a/login", actionLogin)
	http.HandleFunc("/a/create", actionCreate)
	http.Handle("/a/donate", actionDonate)
	http.HandleFunc("/a/paysystem/qiwi", actionPaymentQiwi)
	http.Handle("/a/generate", actionGenerate)
	http.Handle("/a/settings/clearHistory", actionSettingsClearHistory)
	http.Handle("/a/settings", actionSettings)
	http.Handle("/head/", NewUserHandler().Ok(func(w http.ResponseWriter, r *http.Request, user *db.User) {
		id := regexp.MustCompile("/head/(.+)/raw").FindStringSubmatch(r.URL.Path)
		if len(id) != 0 {
			h := db.GetHeadByID(id[1])
			if h != nil {
				http.Redirect(w, r, strings.Replace(h.Head, "heads/", "/files/", 1), http.StatusTemporaryRedirect)
			} else {
				http.ServeFile(w, r, "static/fallback.png")
			}
			return
		}

		id = regexp.MustCompile("/head/(.+)").FindStringSubmatch(r.URL.Path)
		if len(id) == 0 {
			http.ServeFile(w, r, "static/fallback.png")
			return
		}
		h := db.GetHeadByID(id[1])
		if h == nil {
			http.ServeFile(w, r, "static/fallback.png")
			return
		}
		http.ServeFile(w, r, h.Head)
	}))
	http.Handle("/donate", donate)
	http.Handle("/avatar", avatar)
	http.Handle("/create", create)
	http.Handle("/settings", settings)
	http.Handle("/", index)
	http.ListenAndServe(HOST, nil)
}
