package messager

import (
	"time"

	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/services/messager/proto"
)

type messager struct {
}

func (m *messager) SendMessage(request *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	message := &models.Message{
		SenderID:     request.Sender.ID,
		SenderType:   request.Sender.Role,
		ReceiverID:   request.Receiver.ID,
		ReceiverType: request.Receiver.Role,
		ContentType:  request.ContentType,
		Content:      request.Content,
	}
	err := models.CreateMessage(message)
	if err != nil {
		return nil, err
	}
	return &proto.SendMessageResponse{
		MessageID: message.ID,
	}, nil
}

func (m *messager) ListMessage(request *proto.ListMessageRequest) (*proto.ListMessageResponse, error) {
	messages , err := models.FindUserMessages(models.FindUserMessageOption{
		UserID:      request.UserID,
		Orientation: request.TimestampDirection,
		Timestamp:   time.Unix(int64(request.Timestamp)/1000, int64(request.Timestamp)*1e6),
		Count:       uint(request.Count),
		Status:      request.Status,
	})
	if err != nil {
		return nil, err
	}

	
}

func (m *messager) UndoMessage(*proto.UndoMessageRequest) (*proto.UndoMessageResponse, error) {

}
