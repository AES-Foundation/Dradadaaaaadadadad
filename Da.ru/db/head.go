package db

import (
	"time"
)

type Head struct {
	ID        string `db:"id"`
	User      string `db:"user"`
	Skin      string `db:"skin"`
	Head      string `db:"head"`
	Params    string `db:"params"`
	CreatedAt int64  `db:"createdAt"`
}

func NewHead(user *User, skin string, head string, params string) *Head {
	h := &Head{
		ID:        newId(),
		User:      user.ID,
		Skin:      skin,
		Head:      head,
		Params:    params,
		CreatedAt: time.Now().Unix(),
	}
	db.Exec("INSERT INTO `heads` VALUES (?, ?, ?, ?, ?, ?)", h.ID, h.User, h.Skin, h.Head, h.Params, h.CreatedAt)
	return h
}

func GetHeadByID(id string) *Head {
	h := &Head{}
	if db.Get(h, "SELECT * FROM `heads` WHERE `id`=?", id) != nil {
		return nil
	}
	return h
}
