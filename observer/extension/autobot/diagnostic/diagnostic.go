package diagnostic

import (
	"log"
//	"sync"
 	"strconv"
    "fmt"

	"fort.plus/messengers/telegram"
    "fort.plus/mbroker"

	"fort.plus/repository"
// 	httpTransport "fort.plus/transport"
)

const (
	EMPTY_RESPONSE = ""
	SUBJECT_PREFIX = "diagnostic"
)

type chuckResponse struct {
	TextOfJoke string `json:"value"`
}

func init() {
	repository.Register("/check .+", check)
}

var check = func(message repository.RegExComparator) {
    var telegramMessage telegram.Message = message.(telegram.Message)

	log.Printf("diagnostic:check(), message is:%s", telegramMessage.Text)

	batchId := mbroker.GetBatchId(telegramMessage.Text)
	subject := SUBJECT_PREFIX + "."+batchId

  	from, _ := strconv.ParseInt(telegramMessage.From, 10, 64)

 	//telegram.SendTextMessage(from, fmt.Sprintf(" placed job request in queue:%s, waiting for response", subject))
   
    response, err := mbroker.SendRequest(subject, telegramMessage.Text, 20)

    if err != nil {
        telegram.SendTextMessage(from, fmt.Sprintf(" can't get response for request:%s", subject))
    } else {
        telegram.SendTextMessage(from, fmt.Sprintf("response is:%s", response))
    }

}

