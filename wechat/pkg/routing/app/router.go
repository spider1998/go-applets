package app

import (
	routing "github.com/go-ozzo/ozzo-routing"
)

// Register .
func Register(router *routing.RouteGroup) {
	router.Use(sessionChecker)
	{
		handler := NewMemberHandler()
		router.Post("/bind", handler.BindMember)       //绑定用户
		router.Get("/check-token", handler.CheckToken) //验证token（微信服务器使用）
	}

}
