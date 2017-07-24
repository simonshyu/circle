package handler

import (
	"github.com/labstack/echo"
)

type CircleContext struct {
	echo.Context
}

func (c CircleContext) DefaultQueryParam(name string, defaultVal string) string {
	value := c.QueryParam(name)
	if value == "" {
		return defaultVal
	}
	return value
}
