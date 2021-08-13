// +build wireinject

package main

import (
	demo_v1 "myapp/api/demo/v1"
	"myapp/internal/biz"
	"myapp/internal/data"
	"myapp/internal/service"

	"github.com/google/wire"
)

func initApp() (*demo_v1.TestResponse, error) {
	panic(wire.Build(data.NewMySQL, biz.NewBizDemo, service.NewServiceDemo))
}
