package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"log"
	"testing"
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
