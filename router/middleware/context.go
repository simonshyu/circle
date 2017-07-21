package middleware

import (
	"github.com/SimonXming/circle/handler"
	"github.com/labstack/echo"
)

func CircleContext() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &handler.CircleContext{c}
			return h(cc)
		}
	}
}
