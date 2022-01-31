package notify

/*
	Send some messages to telegram channel
*/
import (
	"encoding/json"
	"log"

	diff "fort.plus/filter"
	"fort.plus/im"
	"fort.plus/mbroker"
	"fort.plus/repository"
)

var (
	threshold  int
	notifyAddr string
	carrier    im.Carrier
)

//
//  prepare repository to store reference between callback function and patterns from config file
//
func Register(m im.Carrier, to string, th int) {
	carrier = m
	notifyAddr = to
	threshold = th
	//  notifyPatterns := []//config.GetCurrent().ExtSnifferNotifyTelegram
	// for _, pattern := range notifyPatterns {
	// 	repository.Register(pattern, notifyTelegram)
	// }
}

var notifyTelegram = func(message repository.RegExComparator) {
	var natsMessage mbroker.NatsMessage = message.(mbroker.NatsMessage)

	var syslogMessage mbroker.SyslogMessage

	json.Unmarshal([]byte(natsMessage.Text), &syslogMessage)

	log.Printf("notifyTelegram, message is:%s", syslogMessage.AsText())

	if diff.IsThresholdExceeded(natsMessage.Text, threshold) {
		log.Printf("Message look the same, skip sending:%s", syslogMessage.AsText())
	} else {
		carrier.Send(notifyAddr, syslogMessage.AsText())
	}
}
