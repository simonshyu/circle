package handler

import (
	// "crypto/sha256"
	// "encoding/base32"
	"fmt"
	// "github.com/gorilla/securecookie"

	// "github.com/SimonXming/circle/model"
	// "github.com/SimonXming/circle/remote"
	"github.com/SimonXming/circle/store"
	// "github.com/SimonXming/circle/utils"
	// "github.com/SimonXming/circle/utils/httputil"
	// "github.com/SimonXming/circle/utils/token"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func PostBuild(c echo.Context) error {
	repoID, err := strconv.ParseInt(c.Param("repoID"), 10, 64)
	if err != nil {
		c.Error(err)
		return err
	}

	repo, err := store.RepoLoad(c, repoID)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}

	conf, err := store.ConfigFind(c, repo)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return err
	}
	fmt.Printf("%v", conf.Data)

	return nil
}
