package form

type BindMemberRequest struct {
	AuthCode string `json:"auth_code"`
	PersonID string `json:"person_id"`
	Mobile   string `json:"mobile"`
	OpenID   string `json:"open_id"`
}

type CollectFormIDRequest struct {
	FormID string `json:"form_id"`
}
