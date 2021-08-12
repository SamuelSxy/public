package data

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type mysql struct {
	db *sql.DB
}

func NewMySQL(addr string) *mysql {
	mydb, _ := sql.Open("mysql", "root@tcp("+addr+")/test")
	return &mysql{db: mydb}
}

func (m *mysql) GetById(id string) (string, error) {
	var name string
	sqlScript := "select name from user where id = ?"
	err := m.db.QueryRow(sqlScript, id).Scan(name)

	if err != nil {
		if err == sql.ErrNoRows {
			return name, errors.Wrap(err, fmt.Sprintf("sqlScript: %s\r\nerror: %v\r\n", sqlScript, err))
		} else {
			return name, errors.Wrap(err, "other sqlErr")
		}
	}
	return name, nil

}