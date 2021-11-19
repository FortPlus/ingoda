package notify

/*
	Send some messages to telegram channel
*/
import(
	 "log"

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

func init() {
	repository.Register(".*BGP.*", notifyTelegram)
    //asa
	repository.Register("Lost Failover", notifyTelegram)

    repository.Register(".*EC.*BUNDLE.*", notifyTelegram)
    repository.Register(".*HSRP.*CHANGE.*", notifyTelegram)
    repository.Register(".*OSPF.*ADJ.*", notifyTelegram)
    repository.Register(".*OSPF.*ADJ.*", notifyTelegram)
    repository.Register(".*ower to module.*", notifyTelegram)
    repository.Register(".*FRU.*ower.*", notifyTelegram)
    repository.Register(".*ower.*upply.*", notifyTelegram)
}

var notifyTelegram = func(message repository.RegExComparator)  {
    var syslogMessage mbroker.SyslogMessage = message.(mbroker.SyslogMessage)
    log.Printf("notifyTelegram, message is:%s",syslogMessage.Text)

    if diff.IsThresholdExceeded(syslogMessage.Text, DIFF_FILTER_THRESHOLD) {
	    telegram.SendMessage(TELEGRAM_NOTIFICATION_GROUP, syslogMessage.Text)
	} else {
    	log.Printf("Message is diff filtered, skip sending:%s",syslogMessage.Text)
	}
}

