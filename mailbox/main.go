package main

import (
	"github.com/toomore/mailbox/mailbox/cmd"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cmd.Execute()
}
