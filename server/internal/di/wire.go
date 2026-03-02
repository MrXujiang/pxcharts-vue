//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"mvtable/internal/app"
	"mvtable/internal/pkg/config"
	"mvtable/internal/server"
)

// App 应用结构
type App struct {
	HTTPServer *server.HTTPServer
}

// InitializeApp 初始化应用
func InitializeApp(cfg *config.Config) (*App, func(), error) {
	panic(wire.Build(
		app.Set,
		server.NewRouter,
		server.NewHTTPServer,
		wire.Struct(new(App), "*"),
	))
}
