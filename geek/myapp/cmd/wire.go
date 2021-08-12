// +build wireinject

package main

import (
	"myapp/internal/biz"
	"myapp/internal/data"
	"myapp/internal/service"

	"github.com/google/wire"
)

// initApp init application.
func initApp() {
	panic(wire.Build(data.Test, biz.NewTestRequst, service.NewTest, NewApp()))
}
