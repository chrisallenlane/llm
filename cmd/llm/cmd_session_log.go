package main

import (
	"github.com/chrisallenlane/llm/internal/session"
	"gorm.io/gorm"
)

func cmdSessionLog(_ map[string]interface{}, sess session.Session, _ *gorm.DB) {
	// display the matching messages
	for _, msg := range sess.Messages {
		msg.Display()
	}
}
