package modelDB

import "time"

type NotificationQueue struct {
	ID             int64
	MessageQueueId int64
	Message        string // transfer prop., does not exists as column
	ChatId         int64
	ExternalChatId int64 // transfer prop., does not exists as column
	IsSent         bool
	IsActive       bool
	CreatedAt      time.Time
	ScheduledFor   time.Time
}

// NewNotificationQueue - creating a new instance of NotificationQueue
func NewNotificationQueue() *NotificationQueue {
	return &NotificationQueue{}
}

// NewNotificationQueueConstructor -constructor of NotificationQueue
func NewNotificationQueueConstructor(messageQueueId int64, chatId int64, scheduledFor time.Time) *NotificationQueue {
	return &NotificationQueue{
		MessageQueueId: messageQueueId,
		ChatId:         chatId,
		ScheduledFor:   scheduledFor,
	}
}

func (not *NotificationQueue) Inc(i int) int {
	return i + 1
}
