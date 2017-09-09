package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"log"
)

const (
	// SQLPATH is for SQL link path
	SQLPATH = "root:mailboxdbs@tcp(MARIADB:3306)/mailbox?parseTime=true"
)

// GenSeed is to gen seed
func GenSeed() []byte {
	var buf = make([]byte, 4)
	rand.Read(buf)
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
		log.Fatal("[GetConn] ", err)
	}
	conn.SetMaxOpenConns(1024)
	if err := conn.Ping(); err != nil {
		log.Fatal("[GetConn][Ping] ", err)
	}

	return conn
}
