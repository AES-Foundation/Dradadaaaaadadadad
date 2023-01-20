package db

import (
	"crypto/sha256"
	"encoding/hex"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	HOST     = "127.0.0.1"
	USER     = "root"
	PASSWORD = "root"
	DBNAME   = "root"
)

var db *sqlx.DB

func OpenDB() {
	var err error
	db, err = sqlx.Connect("mysql", USER+":"+PASSWORD+"@tcp("+HOST+")/"+DBNAME)
	if err != nil {
		log.Fatalln("MySQL error: " + err.Error())
	}
}

func hash(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

func newId() string {
	i, _ := uuid.NewRandom()
	return i.String()
}
