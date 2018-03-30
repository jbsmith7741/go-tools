package sqlh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStmt_Query(t *testing.T) {
	db, _ := Open(Config{Dialect: Mock})
	_, err := db.Query("Select * from table where column = $1 and column = $2", 1, 2)
	assert.Nil(t, err)

	_, err = db.Query("Select * from table where column = ? and column = ?", 1, 2)
	assert.Nil(t, err)
}
