package di

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/voice0726/todo-app-api/config"
	"github.com/voice0726/todo-app-api/feats/address"
	"github.com/voice0726/todo-app-api/feats/todo"
	"github.com/voice0726/todo-app-api/infra"
	"github.com/voice0726/todo-app-api/logger"
	"github.com/voice0726/todo-app-api/server"
)

func PrepareApp() *fx.App {
	return fx.New(
		config.Module,
		fx.Provide(logger.NewLogger),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
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
		fx.Invoke(func(base *infra.DataBase) {}),
		fx.Invoke(func(s *server.Server) {}),
	)
}
