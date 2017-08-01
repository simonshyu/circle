package gitlab

import (
	// "crypto/md5"
	// "encoding/hex"
	"fmt"
	// "net/url"
	"strconv"
	"strings"

	"github.com/simonshyu/circle/remote/gitlab/client"
)

const (
	gravatarBase = "https://www.gravatar.com/avatar"
)

// NewClient is a helper function that returns a new GitHub
// client using the provided OAuth token.
// func NewClient(url, accessToken string, skipVerify bool) *client.Client {
func NewClient(url, token string, skipVerify bool) *client.Client {
	client := client.New(url, "/api/v3", token, skipVerify)
	return client
}

func ns(owner, name string) string {
	return fmt.Sprintf("%s%%2F%s", owner, name)
}

func ExtractFromPath(str string) (string, string, error) {
	s := strings.Split(str, "/")
	if len(s) < 2 {
		return "", "", fmt.Errorf("Minimum match not found")
	}
	return s[0], s[1], nil
}

func GetProjectId(r *Gitlab, c *client.Client, owner, name string) (projectId string, err error) {
	if r.Search {
		_projectId, err := c.SearchProjectId(owner, name)
		if err != nil || _projectId == 0 {
			return "", err
		}
		projectId := strconv.Itoa(_projectId)
		return projectId, nil
	} else {
		projectId := ns(owner, name)
		return projectId, nil
	}
}
