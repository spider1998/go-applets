package entity

import (
	"sdkeji/wechat/pkg/util"
)

type MemberState int

const (
	MemberStateBind   MemberState = 1
	MemberStateUnBind MemberState = 2
)

type Member struct {
	ID         string      `json:"id" xorm:"pk"`
	PersonID   string      `json:"person_id"`
	Mobile     string      `json:"mobile"`
	OpenID     string      `json:"open_id"`
	State      MemberState `json:"state"`
	CreateTime util.Time   `json:"create_time" xorm:"created"`
	UpdateTime util.Time   `json:"update_time" xorm:"updated"`
}
