package internal

import (
	"github.com/go-ozzo/ozzo-routing"
	"sdkeji/wechat/pkg/app"
	"sdkeji/wechat/pkg/code"
)

func Register(router *routing.RouteGroup) {
	router.Use(internalTokenChecker)
	{
		handler := NewMessageHandler()
		router.Post("/message", handler.SendMessage) //推送模板消息
	}
	{
		handler := NewMemberHandler()
		router.Patch("/member", handler.ModifyMember) //存储formID
	}
}

func internalTokenChecker(c *routing.Context) error {
	token := c.Request.Header.Get("X-Internal-Token")
	if token == "" || token != app.Conf.InternalToken {
		return code.Error(code.InvalidInternalToken)
	}
	return c.Next()
}
