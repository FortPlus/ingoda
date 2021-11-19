package telegram
import ("regexp")

type TelegramMessage struct {
	ChatId    int64  `json:"chat_id"`
	ParseMode string `json:"parse_mode"`
	Text      string `json:"text"`
}

type teleChat struct {
	ChatId int64  `json:"id"`
	Type   string `json:"type"`
}

type teleMessage struct {
	MessageId int      `json:"message_id"`
	Date      int64    `json:"date"`
	Text      string   `json:"text"`
	Chat      teleChat `json:"chat"`
}

type TelegramUpdate struct {
	UpdateId    int64       `json:"update_id"`
	ChannelPost teleMessage `json:"channel_post"`
	Message     teleMessage `json:"message"`
}

type TelegramUpdates struct {
	Ok     bool             `json:"ok"`
	Result []TelegramUpdate `json:"result"`
}

type Message struct {
	From string
	Text string
}
func (m Message) IsRegExEqual(pattern string) (bool, error) {
//     fmt.Printf("a1 function, message is:%s\n",pattern)
    return regexp.MatchString(pattern, m.Text)
}
