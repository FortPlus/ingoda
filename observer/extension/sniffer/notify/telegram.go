package notify

/*
	Send some messages to telegram channel
*/
import (
	"encoding/json"
	"log"
	"time"

	diff "fort.plus/filter"
	"fort.plus/im"
	"fort.plus/listmanager"
	"fort.plus/mbroker"
	"fort.plus/repository"
)

var (
	threshold  int
	notifyAddr string
	carrier    im.Carrier
)

//
// Set notification parameters
//
func SetCarrier(m im.Carrier, to string, th int) {
	carrier = m
	notifyAddr = to
	threshold = th
}

//
// periodically load patterns for notification
//
func LoadPatterns(p *listmanager.ListRecords, period int) {
	go func() {
		for {
			log.Println("LoadPatterns")
			patterns := p.GetPatterns()
			if len(patterns) != 0 {
				repository.Clear()
				for _, pattern := range patterns {
					log.Println("register pattern:", pattern)
					repository.Register(pattern, notifyTelegram)
				}
			}
			time.Sleep(time.Duration(period) * time.Second)
		}
	}()
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
