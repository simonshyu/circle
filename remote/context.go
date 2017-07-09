package remote

import (
	"github.com/labstack/echo"
)

const key = "remote"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Remote associated with this context.
func FromContext(c echo.Context) Remote {
	return c.Get(key).(Remote)
}

// ToContext adds the Remote to this context if it supports
// the Setter interface.
func ToContext(c Setter, r Remote) {
	c.Set(key, r)
}

/*
account 需要调用 convertToRemote 方法
convertToRemote 需要调用 gitlab.New 方法
*/
