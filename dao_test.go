package error_Test

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type userData struct {
	id   int32
	name string
}

func dao(id int) (userData, error) {
	var res userData

	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/test")

	//@todo:err...

	sqlScript := "select * from user where id = ?"
	err = db.QueryRow(sqlScript, id).Scan(&res.id, &res.name)

	if err != nil {
		if err == sql.ErrNoRows {
			return res, errors.Wrap(err, fmt.Sprintf("sqlScript: %s\r\nerror: %v\r\n", sqlScript, err))
		} else {
			return res, errors.Wrap(err, "other sqlErr")
		}
	}

	return res, nil
}

func TestDao(t *testing.T) {

	res, err := dao(999)

	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println(fmt.Sprintf("%+v", err))
		return
	}

	if err != nil {
		fmt.Println(fmt.Sprintf("%+v", err))
		return
	}

	t.Log(res)
	return
}
