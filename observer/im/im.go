package im

import (
	"regexp"

	"fort.plus/repository"
)

type Message struct {
	From string
	Text string
}

type Encoder interface {
	Encode() Message
}

type Carrier interface {
	GetMessages() ([]Message, error)
	Send(to string, text string) error
}

var _ repository.RegExComparator = Message{}

func (m Message) IsRegExEqual(pattern string) (bool, error) {
	return regexp.MatchString(pattern, m.Text)
}

//TODO: add error handling
func Cast(msg repository.RegExComparator) Message {
	var telegramMessage Message = msg.(Message)
	return telegramMessage
}
