package message

import "github.com/sashabaranov/go-openai"

// Pack packages a native Message object into an openai.ChatCompletionMessage
func (m *Message) Pack() openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Content: m.Content,
		Name:    m.Name,
		Role:    m.Role, // TODO: use an `enum` for this?
	}
}
