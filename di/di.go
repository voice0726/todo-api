package di

import (
	"github.com/labstack/echo/v4"
	"github.com/voice0726/todo-app-api/config"
	"github.com/voice0726/todo-app-api/feats/address"
	"github.com/voice0726/todo-app-api/feats/todo"
	"github.com/voice0726/todo-app-api/infra"
	"github.com/voice0726/todo-app-api/server"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func PrepareApp() *fx.App {
	return fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		config.Module,
		fx.Provide(zap.NewDevelopment),
		fx.Provide(infra.NewDB),
		fx.Provide(echo.New),
		fx.Provide(
			fx.Annotate(todo.NewRepositoryImpl, fx.As(new(todo.Repository))),
			todo.NewHandler,
		),
		fx.Provide(
			fx.Annotate(address.NewRepositoryImpl, fx.As(new(address.Repository))),
			address.NewHandler,
		),
		fx.Provide(server.NewServer),
		fx.Invoke(func(*server.Server) {}),
	)
}
