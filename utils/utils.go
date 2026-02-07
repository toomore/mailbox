package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

const (
	// SQLPATH is for SQL link path
	SQLPATH = "root:mailboxdbs@tcp(MARIADB:3306)/mailbox?parseTime=true"
)

var (
	dbOnce     sync.Once
	dbInstance *sql.DB
)

// getDSN returns DSN from env MAILBOX_DB_DSN or default SQLPATH
func getDSN() string {
	if s := os.Getenv("MAILBOX_DB_DSN"); s != "" {
		return s
	}
	return SQLPATH
}

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
	dbOnce.Do(func() {
		var err error
		dbInstance, err = sql.Open("mysql", getDSN())
		if err != nil {
			log.Fatal("[GetConn] ", err)
		}
		dbInstance.SetMaxOpenConns(50)
		dbInstance.SetMaxIdleConns(10)
		if err := dbInstance.Ping(); err != nil {
			log.Fatal("[GetConn][Ping] ", err)
		}
	})
	return dbInstance
}

// FormatEmail is to make mail address unique
func FormatEmail(email string) string {
	re := regexp.MustCompile(`\+[^@]*`)
	email = re.ReplaceAllString(strings.TrimSpace(strings.ToLower(email)), "")

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return strings.TrimSpace(strings.ToLower(email))
	}
	localPart, domain := parts[0], parts[1]

	localPart = strings.ReplaceAll(localPart, ".", "")

	return localPart + "@" + domain
}
