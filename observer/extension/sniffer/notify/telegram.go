package notify

/*
	Send some messages to telegram channel
*/
import(
	 "log"
     "encoding/json"

	 "fort.plus/repository"
	 "fort.plus/messengers/telegram"
	 "fort.plus/mbroker"
	 "fort.plus/config"
	 "fort.plus/filter"
)

var (
    DIFF_FILTER_THRESHOLD int = config.GetCurrent().TelegramNotificationDiffThreshold
    TELEGRAM_NOTIFICATION_GROUP int64 = config.GetCurrent().TelegramNotificationGroup
)

//
//  prepare repository to store reference between callback function and patterns from config file
//
func init() {
    notifyPatterns := config.GetCurrent().ExtSnifferNotifyTelegram
    for _, pattern := range notifyPatterns {
	    repository.Register(pattern, notifyTelegram)
    }
}

var notifyTelegram = func(message repository.RegExComparator)  {
    var natsMessage mbroker.NatsMessage = message.(mbroker.NatsMessage)

    var syslogMessage mbroker.SyslogMessage

    json.Unmarshal([]byte(natsMessage.Text), &syslogMessage)

    log.Printf("notifyTelegram, message is:%s",syslogMessage.AsText())

    if diff.IsThresholdExceeded(natsMessage.Text, DIFF_FILTER_THRESHOLD) {
    	log.Printf("Message look the same, skip sending:%s",syslogMessage.AsText())
    } else {
	telegram.SendTextMessage(TELEGRAM_NOTIFICATION_GROUP, syslogMessage.AsText())
    }
}


