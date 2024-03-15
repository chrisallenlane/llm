package message

import (
	"time"

	"github.com/sashabaranov/go-openai"
)

// Unpack openai.ChatCompletionMessage into a native Message object
func Unpack(m openai.ChatCompletionMessage, sessID uint) Message {
	return Message{
		Content:   m.Content,
		Date:      time.Now(),
		Name:      "", // NB: this is unused
		Role:      openai.ChatMessageRoleAssistant,
		SessionID: sessID,
	}
}
