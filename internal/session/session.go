package session

import (
	"github.com/chrisallenlane/llm/internal/message"
	"gorm.io/gorm"
)

// Session models a session.
type Session struct {
	gorm.Model
	Name     string `gorm:"unique"`
	Hint     string
	Messages []message.Message
}
