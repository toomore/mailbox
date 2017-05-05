package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log"

	"github.com/google/uuid"
)

const (
	// SQLPATH is for SQL link path
	SQLPATH = "root:mailboxdbs@tcp(MARIADB:3306)/mailbox?parseTime=true"
)

// GenSeed is to gen seed
func GenSeed() [8]byte {
	var buf [8]byte
	u := uuid.Must(uuid.NewRandom())
	hex.Encode(buf[:], u[:4])
	return buf
}

// GenHmac is to gen hmac
func GenHmac(key, message []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}

// GetConn DB conn
func GetConn() *sql.DB {
	var err error
	var conn *sql.DB
	if conn, err = sql.Open("mysql", SQLPATH); err != nil {
		log.Fatal(err)
	}
	conn.SetMaxOpenConns(1024)
	if err := conn.Ping(); err != nil {
		log.Fatal(err)
	}

	return conn
}
