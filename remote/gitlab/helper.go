package gitlab

import (
	// "crypto/md5"
	// "encoding/hex"
	// "fmt"
	// "net/url"
	// "strconv"
	// "strings"

	"github.com/SimonXming/circle/remote/gitlab/client"
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
