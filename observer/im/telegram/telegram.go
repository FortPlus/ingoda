package telegram

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"fort.plus/fperror"
	"fort.plus/im"
	httpTransport "fort.plus/transport"
)

const (
	TELEGRAM_URL         = "https://api.telegram.org/bot%s/sendMessage"
	TELEGRAM_GET_MSG_URL = "https://api.telegram.org/bot%s/getUpdates?offset=%d"
	MESSAGE_EXPIRED_TIME = 180 //seconds
)

type bot struct {
	lastUpdateId  int64
	telegramToken string
}

var _ im.Carrier = &bot{}

func New(token string) im.Carrier {
	return &bot{lastUpdateId: 0, telegramToken: token}
}

func (t *bot) GetMessages() ([]im.Message, error) {

	var url = fmt.Sprintf(TELEGRAM_GET_MSG_URL, t.telegramToken, t.lastUpdateId)
	var err error
	var response []im.Message
	var res TelegramUpdates

	err = httpTransport.GetAndUnmarshall(url, &res)

	if !res.Ok {
		log.Println(fperror.Warning("got response with false status", err))
		return response, err
	}

	timeNow := time.Now()
	timeUnix := timeNow.Unix()
	for _, element := range res.Result {

		// prepare offset for next requests
		if element.UpdateId >= t.lastUpdateId {
			t.lastUpdateId = element.UpdateId + 1
		}

		if element.Message.Date > 0 {
			if (timeUnix - element.Message.Date) < MESSAGE_EXPIRED_TIME {
				response = append(response, im.Message{
					From: strconv.FormatInt(element.Message.Chat.ChatId, 10),
					Text: element.Message.Text,
				})
				continue
			} else {
				continue
			}
		}

		if element.ChannelPost.Date > 0 {
			if (timeUnix - element.ChannelPost.Date) < MESSAGE_EXPIRED_TIME {
				response = append(response, im.Message{
					From: strconv.FormatInt(element.ChannelPost.Chat.ChatId, 10),
					Text: element.ChannelPost.Text,
				})
				continue
			} else {
				continue
			}
		}

	}
	return response, err
}

func (t *bot) Send(to string, message string) error {
	var url = fmt.Sprintf(TELEGRAM_URL, t.telegramToken)
	var msg TelegramMessage
	chatId, _ := strconv.ParseInt(to, 10, 64)

	msg.SetChatId(chatId)
	msg.SetTextHtml(message)

	err := httpTransport.PostJson(url, msg)
	if err != nil {
		log.Println(fperror.Warning("error while SendTextMessage", err))
	}
	return err
}
