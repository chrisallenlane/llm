// Package message encapsulates message methods and data
package message

import (
	"time"

	"gorm.io/gorm"
)

// Message encapsulates a chat message
type Message struct {
	gorm.Model
	Content   string
	Date      time.Time
	Name      string
	Role      string
	SessionID uint
}
