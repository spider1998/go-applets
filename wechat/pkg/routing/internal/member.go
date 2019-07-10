package internal

import (
	"git.sdkeji.top/share/sdlib/api"
	routing "github.com/go-ozzo/ozzo-routing"
	"sdkeji/wechat/pkg/code"
	"sdkeji/wechat/pkg/entity"
	"sdkeji/wechat/pkg/service"
)

func NewMemberHandler() MemberHandler {
	return MemberHandler{}
}

type MemberHandler struct{}

func (m MemberHandler) ModifyMember(c *routing.Context) error {
	var req entity.ModifyMemberRequest
	err := c.Read(&req)
	if err != nil {
		return code.Error(api.InvalidData).WithDetails(err)
	}
	member, err := service.Member.ModifyMember(req)
	if err != nil {
		return err
	}
	return c.Write(member)
}
