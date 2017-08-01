package middleware

import (
	"github.com/simonshyu/circle/store"
	"github.com/labstack/echo"
	// "github.com/labstack/echo/middleware"
)

func StoreWithConfig(s store.Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			store.ToContext(c, s)
			return next(c)
		}
	}
}
