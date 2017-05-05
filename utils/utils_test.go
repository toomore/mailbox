package utils

import (
	"log"
	"testing"
)

func TestGenSeed(*testing.T) {
	log.Printf("%s", GenSeed())
	log.Printf("%s", GenSeed())
}
