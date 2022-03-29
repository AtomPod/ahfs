package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type UserMessageStatus int

const (
	MaxUserMessageCount = 32
)

const (
	UserMessageStatusNone UserMessageStatus = iota
	UserMessageStatusRead
	UserMessageStatusUnread
	UserMessageStatusUndo
)

type TimestampOrientation int

const (
	TimestampOrientationForward TimestampOrientation = iota // 向后查询
	TimestampOrientationBack                                // 向前查询
)

type UserMessage struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	UserID    uint       `gorm:"index"`
	MessageID uint
	Status    UserMessageStatus
}

func createUserMessage(e *gorm.DB, um *UserMessage) error {
	return e.Create(um).Error
}

type FindUserMessageOption struct {
	UserID      uint
	Status      UserMessageStatus
	Orientation TimestampOrientation
	Timestamp   time.Time
	Count       uint
}

func FindUserMessages(opts FindUserMessageOption) ([]*UserMessage, error) {
	queryEngine := engine
	if opts.UserID != 0 {
		queryEngine = queryEngine.Where("user_id=?", opts.UserID)
	}
	if opts.Status != UserMessageStatusNone {
		queryEngine = queryEngine.Where("status=?", opts.Status)
	}
	if !opts.Timestamp.IsZero() {
		if opts.Orientation == TimestampOrientationBack {
			queryEngine = queryEngine.Where("created_at > ?", opts.Timestamp)
		} else {
			queryEngine = queryEngine.Where("created_at < ?", opts.Timestamp)
		}
	}

	count := opts.Count
	if count == 0 || count > MaxUserMessageCount {
		count = MaxUserMessageCount
	}
	queryEngine = queryEngine.Limit(count)

	userMessages := make([]*UserMessage, 0)
	if err := queryEngine.Find(&userMessages).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return userMessages, nil
}
