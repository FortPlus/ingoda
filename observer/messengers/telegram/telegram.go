package telegram

import (
	"fort.plus/config"
	"fort.plus/fperror"
	"fmt"
	"log"
	"time"
	"strconv"
	httpTransport "fort.plus/transport"
)

const (
	TELEGRAM_URL         = "https://api.telegram.org/bot%s/sendMessage"
	TELEGRAM_GET_MSG_URL = "https://api.telegram.org/bot%s/getUpdates?offset=%d"
	MESSAGE_EXPIRED_TIME = 180 //seconds
)


var lastUpdateId int64 = 0
var TELEGRAM_TOKEN string = config.GetCurrent().TelegramToken

func GetMessages() ([]Message, error) {

	var url = fmt.Sprintf(TELEGRAM_GET_MSG_URL, TELEGRAM_TOKEN, lastUpdateId)
	var err error
	var response []Message
	var res TelegramUpdates

	httpTransport.GetAndUnmarshall(url, &res)

	if !res.Ok {
		log.Panic(fperror.Warning("got response with false status", nil))
		return response, err
	}

	timeNow := time.Now()
	timeUnix := timeNow.Unix()
	for _, element := range res.Result {

		// prepare offset for next requests
		if element.UpdateId >= lastUpdateId {
			lastUpdateId = element.UpdateId + 1
		}

		if element.Message.Date > 0 {
			if (timeUnix - element.Message.Date) < MESSAGE_EXPIRED_TIME {
				response = append(response, Message{
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
				response = append(response, Message{
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

func SendTextMessage(chatId int64, message string) error {
	var url = fmt.Sprintf(TELEGRAM_URL, TELEGRAM_TOKEN)
	var msg TelegramMessage

	msg.SetChatId(chatId)
    msg.SetTextHtml(message)

	err := httpTransport.PostJson(url, msg)
    if err != nil {
        log.Println(fperror.Warning("got error in PostJson", err))
    }
	return err
}





