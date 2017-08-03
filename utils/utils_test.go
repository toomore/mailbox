package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestGenSeed(*testing.T) {
	log.Printf("%s", GenSeed())
	log.Printf("%s", GenSeed())
}

func TestGenHmac(t *testing.T) {
	funcs := GenHmac([]byte("toomore"), []byte("Let it go"))
	t.Logf("%x", funcs)

	mac := hmac.New(sha256.New, []byte("toomore"))
	mac.Write([]byte("Let it go"))
	if hmac.Equal(funcs, mac.Sum(nil)) {
		t.Logf("%x", mac.Sum(nil))
	} else {
		t.Fatalf("%x", mac.Sum(nil))
	}
}

func TestGetConn(t *testing.T) {
	rows, err := GetConn().Query("select id from campaign;")
	defer rows.Close()

	if err != nil {
		t.Log(err)
	}
	var id string
	for rows.Next() {
		rows.Scan(&id)
		t.Log(id)
	}
}

func BenchmarkGenSeed(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenSeed()
	}
}

func BenchmarkGenHmac(b *testing.B) {
	seed := GenSeed()
	msg := []byte("Toomore")
	for i := 0; i < b.N; i++ {
		GenHmac(seed[:], msg)
	}
}
