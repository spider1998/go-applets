package app

import (
	"git.sdkeji.top/share/sdlib/api"
	routing "github.com/go-ozzo/ozzo-routing"
	"sdkeji/wechat/pkg/app"
	"sdkeji/wechat/pkg/code"
	"strings"
)

//
const (
	sessionTokenHeaderKey = "X-Access-Token"
	sessionKey            = "person.operator"
)

func sessionChecker(c *routing.Context) error {
	token := c.Request.Header.Get(sessionTokenHeaderKey)
	if token == "" {
		token = c.Query(strings.ToLower(sessionTokenHeaderKey))
		if token == "" {
			return code.Error(code.InvalidPersonToken).WithMessage("ACCESS_TOKEN_REQUIRED")
		}
	}
	person, err := app.API.Platform().CheckOpenToken(token)
	if err != nil {
		return err
	}
	c.Set(sessionKey, person)
	c.Set(sessionTokenHeaderKey, token)
	return c.Next()
}

func GetSessionPerson(c *routing.Context) api.Person {
	return c.Get(sessionKey).(api.Person)
}
