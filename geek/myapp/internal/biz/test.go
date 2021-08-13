package biz

import (
	"errors"
	"myapp/internal/data"
)

type TestRequest interface {
	SetId(string)
	GetId() string
	BizDemo() TestResponse
}

type TestResponse interface {
	GetMessage() string
}

type Test interface {
	test(string) string
}

type TestReq struct {
	Test
	id string
}

type TestRes struct {
	name string
}

func NewBizDemo(id string) TestReq {
	return TestReq{id: id}
}

func (req *TestReq) BizDemo() TestRes {
	mysql := data.NewMySQL("127.0.0.1")
	name, err := mysql.MysqlConn.GetNameById(req.id)

	if errors.Is(err, data.NoUser) {
		return TestRes{name: "no user"}
	}

	return TestRes{name: name}

}

func (req *TestReq) GetId() string {
	return req.id
}

func (req *TestReq) setId(id string) {
	req.id = id
}

func (res *TestRes) GetMessage() string {
	return res.name
}
