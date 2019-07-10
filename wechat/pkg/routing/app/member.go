package app

import (
	"git.sdkeji.top/share/sdlib/api"
	routing "github.com/go-ozzo/ozzo-routing"
	"sdkeji/wechat/lib/applets"
	"sdkeji/wechat/pkg/app"
	"sdkeji/wechat/pkg/code"
	"sdkeji/wechat/pkg/form"
	"sdkeji/wechat/pkg/service"
)

func NewMemberHandler() MemberHandler {
	return MemberHandler{}
}

type MemberHandler struct{}

func (m MemberHandler) BindMember(c *routing.Context) error {
	var req form.BindMemberRequest
	err := c.Read(&req)
	if err != nil {
		return code.Error(api.InvalidData).WithDetails(err)
	}
	person := GetSessionPerson(c)
	req.Mobile = person.Mobile
	req.PersonID = person.ID
	openID, err := app.Wechat.GetOpenID(req.AuthCode)
	if err != nil {
		return code.Error(api.InvalidData).WithDetails(err)

	}
	req.OpenID = openID
	member, err := service.Member.BindMember(req)
	if err != nil {
		return err
	}
	return c.Write(member)
}

func (m MemberHandler) CheckToken(c *routing.Context) error {
	var req applets.CheckTokenRequest
	err := c.Read(&req)
	if err != nil {
		return code.Error(api.InvalidData).WithDetails(err)
	}
	err = app.Wechat.CheckSignature(req)
	if err != nil {
		return c.Write("failed.")
	}
	return c.Write(req.Echostr)
}
