package model

// swagger:model build
type ScmAccount struct {
	ID           int64  `json:"id,omitempty"            meddler:"scm_id,pk"`
	Host         string `json:"host,omitempty"          meddler:"scm_host"`
	Login        string `json:"login,omitempty"         meddler:"scm_login"`
	Password     string `json:"password,omitempty"      meddler:"scm_password"`
	PrivateToken string `json:"private_token,omitempty" meddler:"scm_private_token"`
	Type         string `json:"scm_type,omitempty"      meddler:"scm_type"`
}

/*
Field Explanation:
*/
