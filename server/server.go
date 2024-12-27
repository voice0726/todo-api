package server

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"
	"runtime"
	"time"

	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/http/middleware"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/voice0726/todo-app-api/config"
	"github.com/voice0726/todo-app-api/feats/address"
	"github.com/voice0726/todo-app-api/feats/todo"
)

var domain = "localhost"
var key = "key.json"

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	dirname := filepath.Dir(filename)
	key = filepath.Join(dirname, key)
}

type Server struct {
	config         *config.Config
	todoHandler    *todo.Handler
	addressHandler *address.Handler
	e              *echo.Echo
	lg             *zap.Logger
}

func (s *Server) GetConfig() *config.Config {
	return s.config
}

func NewServer(lg *zap.Logger, config *config.Config, todoHandler *todo.Handler, addressHandler *address.Handler, lc fx.Lifecycle) (*Server, error) {
	e := echo.New()
	e.Use(echozap.ZapLogger(lg))
	s := &Server{
		e:              e,
		todoHandler:    todoHandler,
		addressHandler: addressHandler,
		lg:             lg,
		config:         config,
	}
	ctx := context.Background()
	authZ, err := authorization.New(ctx, zitadel.New(domain, zitadel.WithInsecure("8081")), oauth.DefaultAuthorization(key))
	if err != nil {
		s.lg.Error("zitadel sdk could not initialize", zap.Error(err))
		return nil, err
	}
	mw := middleware.New(authZ)
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lg.Info("starting server")
			go func() {
				s.Register(mw)
				if err := e.Start(":8000"); err != nil && !errors.Is(err, http.ErrServerClosed) {
					e.Logger.Fatal("error starting server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			lg.Info("server shutting down")
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			if err := e.Shutdown(ctx); err != nil {
				e.Logger.Fatal(err)
			}
			lg.Info("server shutdown")
			return nil
		},
	})
	return s, nil
}

func (s *Server) Register(mw *middleware.Interceptor[*oauth.IntrospectionContext]) {
	tg := s.e.Group("/todos", echo.WrapMiddleware(mw.RequireAuthorization()))
	tg.GET("", s.todoHandler.FindAll)
	tg.GET("/:id", s.todoHandler.Find)
	ag := s.e.Group("/addresses", echo.WrapMiddleware(mw.RequireAuthorization()))
	ag.GET("", s.addressHandler.FindAllByUserID)
	ag.GET("/:id", s.addressHandler.Find)
}
