package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"threehead/db"
)

func GetByIndex(i int, a []string) string {
	if i >= 0 {
		return a[i]
	}
	return a[len(a)+i]
}

func DownloadFile(out string, url string) (string, error) {
	r, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	buffer := bytes.NewBuffer(nil)
	io.Copy(buffer, r.Body)

	hash := sha1.Sum(buffer.Bytes())
	hx := hex.EncodeToString(hash[:])

	output, err := os.Create(out + hx + ".png")
	if err != nil {
		return "", err
	}
	defer output.Close()
	io.Copy(output, buffer)
	return hx, nil
}

type UserHandler struct {
	ifNotOk func(w http.ResponseWriter, r *http.Request)
	ifOk    func(w http.ResponseWriter, r *http.Request, user *db.User)
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		ifNotOk: func(w http.ResponseWriter, r *http.Request) {
			tmpl.ExecuteTemplate(w, "401", nil)
		},
		ifOk: func(w http.ResponseWriter, r *http.Request, user *db.User) {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		},
	}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	key := ""
	for _, v := range cookies {
		if v.Name == "key" {
			key = v.Value
		}
	}
	user := db.GetUserByKey(key)
	if user != nil {
		h.ifOk(w, r, user)
		return
	}
	h.ifNotOk(w, r)
}

func (h *UserHandler) Ok(handler func(w http.ResponseWriter, r *http.Request, user *db.User)) *UserHandler {
	h.ifOk = handler
	return h
}

func (h *UserHandler) NotOk(handler func(w http.ResponseWriter, r *http.Request)) *UserHandler {
	h.ifNotOk = handler
	return h
}

func (h *UserHandler) Anyway(handler func(w http.ResponseWriter, r *http.Request, user *db.User)) *UserHandler {
	h.ifOk = handler
	h.ifNotOk = func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, nil)
	}
	return h
}
