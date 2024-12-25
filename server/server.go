package server

import (
	"context"

	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/voice0726/todo-app-api/feats/address"
	"github.com/voice0726/todo-app-api/feats/todo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Server struct {
	todoHandler    *todo.Handler
	addressHandler *address.Handler
	e              *echo.Echo
	lg             *zap.Logger
}

func NewServer(lg *zap.Logger, todoHandler *todo.Handler, addressHandler *address.Handler, lc fx.Lifecycle) *Server {
	e := echo.New()
	e.Use(echozap.ZapLogger(lg))
	s := &Server{
		e:              e,
		todoHandler:    todoHandler,
		addressHandler: addressHandler,
		lg:             lg,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				lg.Info("server starting")
				s.Register()
				err := s.e.Start(":8000")
				if err != nil {
					lg.Fatal("failed to start server", zap.Error(err))
				}
				lg.Info("server started")
			}()
			return nil
		},
		// todo: graceful shutdown
		// read echo instance from channel
		OnStop: func(ctx context.Context) error {
			return s.e.Shutdown(ctx)
		},
	})
	return s
}

func (s *Server) Register() {
	tg := s.e.Group("/todos", s.Auth)
	tg.GET("", s.todoHandler.FindAll)
	tg.GET("/:id", s.todoHandler.Find)
	ag := s.e.Group("/addresses", s.Auth)
	ag.GET("", s.addressHandler.FindAllByUserID)
	ag.GET("/:id", s.addressHandler.Find)
}
