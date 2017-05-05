package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/toomore/mailbox/utils"
)

var conn *sql.DB

func create() ([8]byte, [8]byte) {
	id, seed := utils.GenSeed(), utils.GenSeed()
	_, err := conn.Query(fmt.Sprintf(`INSERT INTO campaign(id,seed) VALUES('%s', '%s')`, id, seed))
	if err != nil {
		log.Fatal(err)
	}
	return id, seed
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) >= 1 && args[0] == "create" {
		conn = utils.GetConn()
		id, seed := create()
		log.Printf("id: %s, seed: %s", id, seed)
	} else {
		fmt.Println("mailbox_campaign [cmd]\ncmd: `create`")
	}
}
