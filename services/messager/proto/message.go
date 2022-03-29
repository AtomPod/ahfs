package proto

import (
	"time"

	"github.com/czhj/ahfs/models"
)

type MessageAddress struct {
	ID       uint
	Nickname string
	Role     models.MessageRoleType
}

type Message struct {
	ID          uint
	CreatedAt   time.Time
	Sender      MessageAddress
	Receiver    MessageAddress
	ContentType string
	Content     []byte
	Status      models.UserMessageStatus
}

type SendMessageRequest struct {
	Sender      MessageAddress
	Receiver    MessageAddress
	ContentType string
	Content     []byte
}

type SendMessageResponse struct {
	MessageID uint
}

type ListMessageRequest struct {
	Sender             *MessageAddress
	Receiver           *MessageAddress
	UserID             uint
	Timestamp          uint64
	TimestampDirection models.TimestampOrientation
	Offset             int
	Count              int
	Status             models.UserMessageStatus
}

type ListMessageResponse struct {
	Messages []*Message
}

type UndoMessageRequest struct {
	MessageID uint
}

type UndoMessageResponse struct {
}
