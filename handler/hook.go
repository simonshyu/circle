package handler

import (
	"github.com/SimonXming/circle/remote"
	"github.com/SimonXming/circle/store"
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func PostHook(c echo.Context) error {
	scmId, err := strconv.ParseInt(c.QueryParam("scmID"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}
	_, err = store.SetupRemoteWithScmID(c, scmId)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}

	tmprepo, build, err := remote.Hook(c, c.Request())
	if err != nil {
		logrus.Errorf("failure to parse hook. %s", err)
		c.String(http.StatusBadRequest, err.Error())
		return err
	}
	if build == nil {
		logrus.Errorf("Failed generate build from hook request.")
		c.String(http.StatusBadRequest, "Failed generate build from hook request.")
		return err
	}
	if tmprepo == nil {
		logrus.Errorf("failure to ascertain repo from hook.")
		c.String(http.StatusBadRequest, "failure to ascertain repo from hook.")
		return err
	}
	return nil
}
