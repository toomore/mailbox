package reader

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestSave(t *testing.T) {
	// Smoke test: Save does not panic when called with valid args.
	// Requires MariaDB (run sh ./dev-run-mariadb.sh).
	Save("00000001", "1", "127.0.0.1", "test-agent")
}
