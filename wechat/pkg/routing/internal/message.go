package internal

import (
	"git.sdkeji.top/share/sdlib/api"
	routing "github.com/go-ozzo/ozzo-routing"
	"sdkeji/wechat/pkg/code"
	"sdkeji/wechat/pkg/service"
)

func NewMessageHandler() MessageHandler {
	return MessageHandler{}
}

type MessageHandler struct{}

func (m MessageHandler) SendMessage(c *routing.Context) error {
	var req service.MessageWork
	err := c.Read(&req)
	if err != nil {
		return code.Error(api.InvalidData).WithDetails(err)
	}
	err = service.Message.Sender(req)
	if err != nil {
		return err
	}
	return nil
}
