package db

import "time"

type Payment struct {
	ID        string  `db:"id"`
	User      string  `db:"user"`
	Paysystem string  `db:"paysystem"`
	IsPaid    bool    `db:"isPaid"`
	Amount    float64 `db:"amount"`
	CreatedAt int64   `db:"createdAt"`
}

func NewPayment(user *User, amount float64, paysystem string) *Payment {
	payment := &Payment{
		ID:        newId(),
		User:      user.ID,
		Paysystem: paysystem,
		IsPaid:    false,
		Amount:    amount,
		CreatedAt: time.Now().Unix(),
	}
	db.Exec("INSERT INTO `payments` VALUES (?, ?, ?, 0, ?, ?)", payment.ID, payment.User, payment.Paysystem, payment.Amount, payment.CreatedAt)
	return payment
}

func GetPaymentByID(id string) *Payment {
	payment := &Payment{}
	if db.Get(payment, "SELECT * FROM `payments` WHERE `id`=?", id) != nil {
		return nil
	}
	return payment
}

func (p *Payment) Pay() {
	p.IsPaid = true
	db.Exec("UPDATE `payments` SET `isPaid`=1 WHERE `id`=?", p.ID)
}
