package service

import (
	v "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"sdkeji/wechat/pkg/app"
	"sdkeji/wechat/pkg/code"
	"sdkeji/wechat/pkg/entity"
	"sdkeji/wechat/pkg/form"
)

var Member MemberService

type MemberService struct{}

func (m *MemberService) ModifyMember(req entity.ModifyMemberRequest) (member entity.Member, err error) {
	err = v.ValidateStruct(&req,
		v.Field(&req.PersonID, v.Required),
	)
	if err != nil {
		return
	}
	exist, err := app.DB.Where("person_id = ?", req.PersonID).Get(&member)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	if !exist {
		err = code.Error(code.PersonNotExist).WithDetails(err)
		return
	}
	var cols []string
	if req.FormID != "" {
		member.FormID = req.FormID
		cols = append(cols, "form_id")
	}
	_, err = app.DB.ID(member.ID).Cols(cols...).Update(&member)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//绑定用户（一对一）
func (m *MemberService) BindMember(req form.BindMemberRequest) (member entity.Member, err error) {
	err = v.ValidateStruct(&req,
		v.Field(&req.OpenID, v.Required),
		v.Field(&req.Mobile, v.Required),
		v.Field(&req.PersonID, v.Required),
	)
	if err != nil {
		return
	}
	var (
		mem   entity.Member
		exist bool
	)
	exist, err = app.DB.Where("open_id = ?", req.OpenID).And("state = ?", entity.MemberStateBind).Get(&mem)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	if exist {
		if mem.PersonID == req.PersonID {
			app.Logger.Info("member and open id already exist.", "member:", mem.ID)
			member = mem
			return
		} else {
			mem.State = entity.MemberStateUnBind
			_, err = app.DB.ID(mem.ID).Cols("state").Update(&mem)
			if err != nil {
				err = errors.WithStack(err)
				return
			}
			app.Logger.Info("open id already occupied ,unbind member success.", "member:", mem.ID)
		}
	}
	member.ID = xid.New().String()
	member.PersonID = req.PersonID
	member.Mobile = req.Mobile
	member.OpenID = req.OpenID
	member.State = entity.MemberStateBind
	_, err = app.DB.Insert(&member)
	if err != nil {
		return
	}
	app.Logger.Info("bind member success.", "member:", member.ID)
	return
}
