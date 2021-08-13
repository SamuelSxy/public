package service

import (
	demo_v1 "myapp/api/demo/v1"
	"myapp/internal/biz"
)

type Demo struct {
	biz.TestRequest
}

func NewServiceDemo() (*demo_v1.TestResponse, error) {
	request := biz.NewBizDemo("1")
	respone := request.BizDemo()
	return &demo_v1.TestResponse{Name: respone.GetMessage()}, nil
}
