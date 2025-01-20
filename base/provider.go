package base

import (
	"github.com/fenggwsx/filego/base/conf"
	"github.com/fenggwsx/filego/base/flags"
	"github.com/fenggwsx/filego/base/log"
	"github.com/fenggwsx/filego/base/server"
	"github.com/google/wire"
)

var ProviderSetBase = wire.NewSet(
	conf.NewConfig,
	flags.NewFlagSet,
	log.NewLogger,
	server.NewHttpServer,
)
