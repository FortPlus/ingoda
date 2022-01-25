package banmgmt

import (
	"log"
	"regexp"
	"time"

	//	"sync"
	"fmt"
	"strconv"

	"fort.plus/banlist"
	"fort.plus/banlist/webapi"
	"fort.plus/messengers/telegram"

	"fort.plus/config"
	"fort.plus/repository"
	httpTransport "fort.plus/transport"
)

var (
	BAN_SERVER_URI string = config.GetCurrent().BanConnectionString
)

func init() {
	repository.Register("/ban ls", showBanList)
	repository.Register("/ban add [0-9hm]* .*", addBanRecord)
	repository.Register("/ban rm [0-9]*", deleteBanRecord)
	repository.Register("/ban help", helpBan)

}

var showBanList = func(message repository.RegExComparator) {
	var response string = ""
	var telegramMessage telegram.Message = message.(telegram.Message)
	log.Printf("banmgmt:showBanList(), message is:%s", telegramMessage.Text)
	from, _ := strconv.ParseInt(telegramMessage.From, 10, 64)

	var bannedList map[uint32]banlist.Item

	err := httpTransport.GetAndUnmarshall(BAN_SERVER_URI+webapi.GET_LIST_URI, &bannedList)

	if err != nil {
		telegram.SendTextMessage(from, fmt.Sprintf("can't get data from ban server, %s", err))
		return
	}

	for key, item := range bannedList {
		response += fmt.Sprintf("%d:%s:%s\n", key, item.Pattern, item.ExpiredAt)
	}

	telegram.SendTextMessage(from, response)
}

var addBanRecord = func(message repository.RegExComparator) {
	var telegramMessage telegram.Message = message.(telegram.Message)
	log.Printf("banmgmt:addBanRecord(), message is:%s", telegramMessage.Text)
	from, _ := strconv.ParseInt(telegramMessage.From, 10, 64)

	re := regexp.MustCompile("ban add ([0-9hm]*) (.*)")
	match := re.FindStringSubmatch(telegramMessage.Text)
	fmt.Println(match[0], "-", match[1], ":", match[2], "|")

	duration, err := time.ParseDuration(match[1])

	pattern := match[2]

	item := banlist.Item{ExpiredAt: time.Now().Add(duration), Pattern: pattern}

	err = httpTransport.PostJson(BAN_SERVER_URI+webapi.ADD_URI, &item)

	if err != nil {
		telegram.SendTextMessage(from, fmt.Sprintf("can't add record to ban server, %s", err))
		return
	}

	telegram.SendTextMessage(from, "ok")
}

var deleteBanRecord = func(message repository.RegExComparator) {
	var telegramMessage telegram.Message = message.(telegram.Message)
	log.Printf("banmgmt:addBanRecord(), message is:%s", telegramMessage.Text)
	from, _ := strconv.ParseInt(telegramMessage.From, 10, 64)

	re := regexp.MustCompile("ban rm ([0-9]*)")

	match := re.FindStringSubmatch(telegramMessage.Text)
	fmt.Println(match[0], "-", match[1], "|")

	patternId := match[1]

	err := httpTransport.Delete(BAN_SERVER_URI + webapi.BASE_URI + "/" + patternId)

	if err != nil {
		telegram.SendTextMessage(from, fmt.Sprintf("can't add record to ban server, %s", err))
		return
	}

	telegram.SendTextMessage(from, "ok")
}

var helpBan = func(message repository.RegExComparator) {
	var response string
	var telegramMessage telegram.Message = message.(telegram.Message)
	log.Printf("banmgmt:helpBan(), message is:%s", telegramMessage.Text)
	from, _ := strconv.ParseInt(telegramMessage.From, 10, 64)

	response = "/ban ls\n"
	response += "/ban add [0-9mhDYM]* .*\n"
	response += "/ban rm [0-9]*\n"

	telegram.SendTextMessage(from, response)
}
