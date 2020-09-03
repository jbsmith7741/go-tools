package sqlh

import (
	"testing"
)

func TestStmt_Query(t *testing.T) {
	db, _ := Open(Config{Dialect: Mock})
	_, err := db.Query("Select * from table where column = $1 and column = $2", 1, 2)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Query("Select * from table where column = ? and column = ?", 1, 2)
	if err != nil {
		t.Fatal(err)
	}
}
