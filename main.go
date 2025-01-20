package main

import (
	"github.com/fenggwsx/filego/base/server"
)

type App struct {
	server *server.HttpServer
}

func newApp(server *server.HttpServer) *App {
	return &App{
		server: server,
	}
}

func (m *App) start() {
	m.server.Run()
}

func main() {
	// Use `wire` to update the `initApp` func.
	app, err := initApp()
	if err != nil {
		panic(err)
	}

	// Start the app.
	app.start()
}
