package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type MessageRoleType int

const (
	MessageRoleSystem MessageRoleType = iota
	MessageRoleUser
	MessageRoleGroup
)

type Message struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `sql:"index"`

	SenderID     uint `sql:"index"`
	SenderType   MessageRoleType
	ReceiverID   uint `sql:"index"`
	ReceiverType MessageRoleType

	ContentType string
	Content     []byte `gorm:"varchar(512)"`
}

func createMessage(e *gorm.DB, msg *Message) error {
	return e.Create(msg).Error
}

func CreateMessage(msg *Message) error {
	tx := engine.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	if msg.SenderType == MessageRoleUser {
		if exists, err := isUserExists(tx, int64(msg.SenderID)); err != nil {
			return err
		} else if !exists {
			return ErrMessageSenderNotExist{ID: msg.SenderID, Role: msg.SenderType}
		}
	} else {
		return fmt.Errorf("message sender role is not supported")
	}

	if msg.ReceiverType == MessageRoleUser {
		if exists, err := isUserExists(tx, int64(msg.ReceiverID)); err != nil {
			return err
		} else if !exists {
			return ErrMessageReceiverNotExist{ID: msg.ReceiverID, Role: msg.ReceiverType}
		}
	} else {
		return fmt.Errorf("message receiver role is not supported")
	}

	if err := createMessage(tx, msg); err != nil {
		return err
	}

	um := &UserMessage{
		UserID:    msg.ReceiverID,
		MessageID: msg.ID,
		Status:    UserMessageStatusUnread,
	}

	if err := createUserMessage(tx, um); err != nil {
		return err
	}

	return tx.Commit().Error
}
