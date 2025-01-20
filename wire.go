//go:build wireinject
// +build wireinject

package main

import (
	"github.com/fenggwsx/filego/base"
	"github.com/fenggwsx/filego/middleware"
	"github.com/google/wire"
)

func initApp() (*App, error) {
	wire.Build(
		newApp,
		base.ProviderSetBase,
		middleware.ProviderSetMiddleware,
	)
	return &App{}, nil
}
