package messager

import "github.com/czhj/ahfs/services/messager/proto"

type MessageService interface {
	SendMessage(*proto.SendMessageRequest) (*proto.SendMessageResponse, error)
	ListMessage(*proto.ListMessageRequest) (*proto.ListMessageResponse, error)
	UndoMessage(*proto.UndoMessageRequest) (*proto.UndoMessageResponse, error)
}
