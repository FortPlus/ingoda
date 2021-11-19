package telegram

import (
	"fort.plus/config"
	"fort.plus/fperror"
	"encoding/json"
	"fmt"
//	"strings"
	"html"
	"log"
	"time"
	"strconv"
	httpTransport "fort.plus/transport"
)

const (
	PARSE_MODE           = "MarkdownV2"
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
    //message = strings.Replace(message, "=","",-1)
    if len(message) > 4000 {
        message = message[:4000]
    }
	var url = fmt.Sprintf(TELEGRAM_URL, TELEGRAM_TOKEN)
	var msg TelegramMessage
	msg.ChatId = chatId
	msg.Text = html.EscapeString("<pre>"+message+"</pre>")
	msg.ParseMode = "HTML"

	err := httpTransport.PostJson(url, msg)
    if err != nil {
        log.Println(err)
    }
	return err
}

func SendMessage(chatId int64, message string) error {
	var (
	    url = fmt.Sprintf(TELEGRAM_URL, TELEGRAM_TOKEN)
	    msg TelegramMessage
	    err error
	)
	msg.ChatId = chatId
	msg.ParseMode = PARSE_MODE

	decodedMessage, _ := DecodeJson(message)
	msg.Text = decodedMessage

	httpTransport.PostJson(url, msg)
    //if err != nil {
    //    log.Println(fperror.Warning("can't post data", err))
    //}
	return err
}

func DecodeJson(s string) (string, error) {
	var js map[string]interface{}
	var returnValue string

	err := json.Unmarshal([]byte(s), &js)

	if err != nil {
		return s, err
	}
	returnValue = fmt.Sprintf("``` #syslog %s\n%s```", js["host"], js["msg"])

	return returnValue, nil
}
