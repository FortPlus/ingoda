package listmgmt

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"fort.plus/im"
	"fort.plus/listmanager"
	"fort.plus/repository"
	httpTransport "fort.plus/transport"
)

type listBotManager struct {
	serverUri string
	listName  string
	carrier   im.Carrier
}

func Register(carrier im.Carrier, serverUri string, listName string) {
	b := &listBotManager{
		carrier:   carrier,
		serverUri: serverUri,
		listName:  listName,
	}
	repository.Register("/list "+b.listName+" ls", b.getList)
	repository.Register("/list "+b.listName+" add [0-9hm]* .*", b.addRecord)
	repository.Register("/list "+b.listName+" rm [0-9]*", b.deleteRecord)
	repository.Register("/list "+b.listName+" help", b.showHelp)
}

func (b *listBotManager) getList(message repository.RegExComparator) {
	var response string = ""
	msg := im.Cast(message)
	log.Printf("banmgmt:showBanList(), message is:%s", msg.Text)

	var bannedList map[uint32]listmanager.Item

	err := httpTransport.GetAndUnmarshall(b.serverUri+"/api/v1/"+b.listName, &bannedList)

	if err != nil {
		b.carrier.Send(msg.From, fmt.Sprintf("can't get data from ban server, %s", err))
		return
	}

	for key, item := range bannedList {
		response += fmt.Sprintf("%d:%s:%s\n", key, item.Pattern, item.ExpiredAt)
	}

	b.carrier.Send(msg.From, response)
}

func (b *listBotManager) addRecord(message repository.RegExComparator) {
	msg := im.Cast(message)
	log.Printf("banmgmt:addBanRecord(), message is:%s", msg.Text)

	re := regexp.MustCompile("/list " + b.listName + " add ([0-9hm]*) (.*)")
	match := re.FindStringSubmatch(msg.Text)
	fmt.Println(match[0], "-", match[1], ":", match[2], "|")

	duration, err := time.ParseDuration(match[1])
	if err != nil {
		b.carrier.Send(msg.From, fmt.Sprintf("duration %s is incorrect", match[1]))
		return
	}

	item := listmanager.Item{ExpiredAt: time.Now().Add(duration), Pattern: match[2]}
	err = httpTransport.PostJson(b.serverUri+"/api/v1/"+b.listName, &item)
	if err != nil {
		b.carrier.Send(msg.From, fmt.Sprintf("can't add record to ban server, %s", err))
		return
	}
	b.carrier.Send(msg.From, "record added")
}

func (b *listBotManager) deleteRecord(message repository.RegExComparator) {
	msg := im.Cast(message)
	log.Printf("banmgmt:deleteRecord(), message is:%s", msg.Text)

	re := regexp.MustCompile("/list " + b.listName + " rm ([0-9]*)")

	match := re.FindStringSubmatch(msg.Text)
	fmt.Println(match[0], "-", match[1], "|")

	patternId := match[1]

	err := httpTransport.Delete(b.serverUri + "/api/v1/" + b.listName + "/" + patternId)

	if err != nil {
		b.carrier.Send(msg.From, fmt.Sprintf("can't add record to ban server, %s", err))
		return
	}

	b.carrier.Send(msg.From, "record removed")
}

func (b *listBotManager) showHelp(message repository.RegExComparator) {
	var response string
	msg := im.Cast(message)
	log.Printf("listBotManager:showHelp(), message is:%s", msg.Text)

	response = "/list " + b.listName + " ls\n"
	response += "/list " + b.listName + " add [0-9mh]* .*\n"
	response += "/list " + b.listName + " rm [0-9]*\n"

	b.carrier.Send(msg.From, response)
}
