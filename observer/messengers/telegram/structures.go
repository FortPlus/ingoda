package telegram
import (
    "regexp"
    "html"

    "fort.plus/fperror"
)

const (
    MAX_MSG_SIZE = 4000
    HTML_MSG        = "HTML"
    MARKDOWN_MSG    = "MarkdownV2"

)

type TelegramMessage struct {
	ChatId    int64  `json:"chat_id"`
	ParseMode string `json:"parse_mode"`
	Text      string `json:"text"`
}

func (m *TelegramMessage) SetTextHtml(text string) error {
    if len(text) == 0 {
        return fperror.Warning("Message text is empty", nil)
    }
    if len(text) > MAX_MSG_SIZE {
        text = text[:MAX_MSG_SIZE]
    }
    m.ParseMode = HTML_MSG
	m.Text = "<pre>" + html.EscapeString(text)+"</pre>"
	return nil
}

func (m *TelegramMessage) SetChatId(chatId int64) error {
    if chatId == 0 {
        return fperror.Warning("chat ID unspecified", nil)
    }
    m.ChatId = chatId
    return nil
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
    return regexp.MatchString(pattern, m.Text)
}
