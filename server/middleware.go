package server

import (
	"path/filepath"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/http/middleware"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
	"go.uber.org/zap"
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

func (s *Server) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		s.lg.Info("auth middleware", zap.Any("key", key))
		authZ, err := authorization.New(c.Request().Context(), zitadel.New(domain, zitadel.WithInsecure("8081")), oauth.DefaultAuthorization(key))
		if err != nil {
			s.lg.Error("failed to create authorization", zap.Error(err))
		}
		m := middleware.New(authZ)
		m.Context(c.Request().Context())

		// todo: this doesn't work. need look into the echo middleware and zitadel auth middleware
		c.Echo().Use(echo.WrapMiddleware(m.RequireAuthorization()))
		return next(c)
	}
}
