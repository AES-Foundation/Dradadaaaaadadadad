package db

import (
	"math"
	"math/rand"
	"strings"
	"time"
)

type User struct {
	ID                    string `db:"id"`
	Email                 string `db:"email"`
	Password              string `db:"password"`
	DisplayName           string `db:"displayName"`
	IsVerified            bool   `db:"isVerified"`
	IsHiddenInTopDonaters bool   `db:"isHiddenInTopDonaters"`
}

func generateKey() string {
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	key := ""
	for i := 0; i < 64; i++ {
		key += string(alphabet[r.Intn(len(alphabet))])
	}
	return key
}

func Login(email string, password string) (*User, string) {
	user := &User{}
	if err := db.Get(user, "SELECT * FROM `users` WHERE `email`=? AND `password`=?", strings.ToLower(email), hash(password)); err != nil {
		return nil, ""
	}
	key := generateKey()
	db.Exec("INSERT INTO `sessions` VALUES (?, ?, ?, ?)", newId(), key, user.ID, time.Now().Unix()+30*24*60*60)
	return user, key
}

func GetUserByKey(key string) *User {
	user := &User{}
	if err := db.Get(user, "SELECT `users`.* FROM `sessions` JOIN `users` ON `users`.`id` = `sessions`.`user` WHERE `key`=?", key); err != nil {
		return nil
	}
	return user
}

func GetUserByEmail(email string) *User {
	user := &User{}
	if err := db.Get(user, "SELECT * FROM `users` WHERE `email`=?", strings.ToLower(email)); err != nil {
		return nil
	}
	return user
}

func NewUser(email string, password string) *User {
	user := &User{
		ID:                    newId(),
		Email:                 strings.ToLower(email),
		Password:              hash(password),
		DisplayName:           strings.Split(email, "@")[0],
		IsVerified:            true,
		IsHiddenInTopDonaters: false,
	}
	db.Exec("INSERT INTO `users` VALUES (?, ?, ?, ?, ?, ?)", user.ID, user.Email, user.Password, user.DisplayName, user.IsVerified, user.IsHiddenInTopDonaters)
	return user
}

func (u *User) GetPaymentAll() float64 {
	amount := float64(0.0)
	db.Get(&amount, "SELECT SUM(amount) FROM `payments` WHERE `user`=? AND `isPaid`=1", u.ID)
	return amount
}

type Donater struct {
	ID          string  `db:"id"`
	DisplayName string  `db:"displayName"`
	Amount      float64 `db:"amount"`
}

func GetTopDonaters() []*Donater {
	donators := []*Donater{}
	db.Select(&donators, "SELECT `users`.`displayName`, `users`.`id`, SUM(`payments`.`amount`) as amount FROM `payments` JOIN `users` ON `users`.`id` = `payments`.`user` WHERE `payments`.`isPaid`=1 AND `users`.`isHiddenInTopDonaters` = 0 GROUP BY `user` HAVING amount > 0 ORDER BY amount DESC;")
	return donators
}

func (d *Donater) FormattedAmount() float64 {
	return math.Floor(d.Amount)
}

func (u *User) GetLastHeads(max int, offset int) []*Head {
	heads := []*Head{}
	db.Select(&heads, "SELECT * FROM `heads` WHERE `user`=? ORDER BY `createdAt` DESC LIMIT ?", u.ID, max)
	return heads
}

func (u *User) SetHideIntTopDonaters(val bool) {
	u.IsHiddenInTopDonaters = val
	db.Exec("UPDATE `users` SET `isHiddenInTopDonaters`=? WHERE `id` = ?;", val, u.ID)
}

func (u *User) SetDisplayName(val string) {
	u.DisplayName = val
	db.Exec("UPDATE `users` SET `displayName`=? WHERE `id` = ?;", val, u.ID)
}

func (u *User) ClearHeadsHistory() {
	db.Exec("DELETE FROM `heads` WHERE `user` = ?;", u.ID)

}
