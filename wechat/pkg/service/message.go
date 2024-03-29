package service

import (
	v "github.com/go-ozzo/ozzo-validation"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/pkg/errors"
	"sdkeji/wechat/lib/applets"
	"sdkeji/wechat/pkg/app"
	"sdkeji/wechat/pkg/code"
	"sdkeji/wechat/pkg/entity"
)

var Message = &MessageService{
	MsgQueue: make(chan MessageWork, 10000),
}

type MessageWork struct {
	PersonID string `json:"person_id"`
	Message  map[string]interface{}
}

type MessageService struct {
	MsgQueue chan MessageWork
}

func (m *MessageService) Boot() error {
	go func() {
		for {
			msg := <-m.MsgQueue
			e := m.SendMsg(msg)
			if e != nil {
				app.Logger.Warn("MsgQueueSend err.", "error", e)
			}
		}
	}()
	return nil
}
func (m *MessageService) SendMsg(msg MessageWork) (err error) {
	var (
		mem   entity.Member
		exist bool
	)
	exist, err = app.DB.Where("person_id = ?", msg.PersonID).And("state = ?", entity.MemberStateBind).Get(&mem)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	if !exist {
		err = code.Error(code.UserNotBindWechat).WithMessage("user not bind wechat.")
		return
	}
	//获取form_id
	formID, err := app.Redis.Cmd("GET", "wx_id:"+msg.PersonID).Str()
	if err != nil {
		if err == redis.ErrRespNil {
			err = code.Error(code.FormIDNotExist)
			return err
		}
		err = errors.WithStack(err)
		return err
	}

	var sendReq applets.SendRequest
	sendReq.Touser = mem.OpenID
	sendReq.Data = msg.Message
	sendReq.FormID = formID
	err = app.Wechat.SendMsg(sendReq)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

func (m *MessageService) Sender(work MessageWork) error {
	err := v.ValidateStruct(&work,
		v.Field(&work.PersonID, v.Required),
		v.Field(&work.Message, v.Required),
	)
	if err != nil {
		return err
	}
	select {
	case m.MsgQueue <- work:
		return nil
	default:
		return errors.New("fail to add job to queue")
	}
}
