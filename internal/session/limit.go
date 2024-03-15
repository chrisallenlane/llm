// Package session manages session state
package session

import (
	"github.com/chrisallenlane/llm/internal/message"
)

// Limit constrains `num` messages to a session
func (s *Session) Limit(num int) []message.Message {
	msgs := s.Messages

	if num > len(msgs) {
		num = len(msgs)
	}

	return msgs[len(msgs)-num:]
}
