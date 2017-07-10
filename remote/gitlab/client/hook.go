package client

// import (
// 	"encoding/json"
// 	"net/url"
// )

const (
	projectUrlHooks = "/projects/:id/hooks"          // Get list of project hooks
	projectUrlHook  = "/projects/:id/hooks/:hook_id" // Get single project hook
)

type Hook struct {
	Id           int    `json:"id,omitempty"`
	Url          string `json:"url,omitempty"`
	CreatedAtRaw string `json:"created_at,omitempty"`
}

/*
Add new project hook.
    POST /projects/:id/hooks
Parameters:
    id                    The ID or NAMESPACE/PROJECT_NAME of a project
    hook_url              The hook URL
    push_events           Trigger hook on push events
    issues_events         Trigger hook on issues events
    merge_requests_events Trigger hook on merge_requests events
*/
// Get a list of projects owned by the authenticated user.
func (c *Client) AddProjectHook(id, link, token string, allow_push, ssl_verify bool) error {
	hookOptions := QMap{
		"url":   link,
		"token": token,
	}

	if allow_push {
		hookOptions["push_events"] = "true"
	} else {
		hookOptions["push_events"] = "false"
	}
	if ssl_verify {
		hookOptions["enable_ssl_verification"] = "true"
	} else {
		hookOptions["enable_ssl_verification"] = "false"
	}

	url, opaque := c.ResourceUrl(projectUrlHooks, QMap{":id": id}, hookOptions)

	_, err := c.Do("POST", url, opaque, nil)
	return err
}
