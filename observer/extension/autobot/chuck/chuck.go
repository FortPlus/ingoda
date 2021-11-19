package chuck

import (
	"log"
	"strconv"

	"fort.plus/messengers/telegram"
	"fort.plus/repository"
	httpTransport "fort.plus/transport"
)

const (
	CHUCK_URI      = "https://api.chucknorris.io/jokes/random"
	EMPTY_RESPONSE = ""
)

type chuckResponse struct {
	TextOfJoke string `json:"value"`
}

func init() {
	repository.Register(".*[Cc]huck.*", a)
}

var a = func(message repository.RegExComparator) {
    var telegramMessage telegram.Message = message.(telegram.Message)
	log.Printf("any:subscriber function, message is:%s", telegramMessage.Text)
 	from, _ := strconv.ParseInt(telegramMessage.From, 10, 64)
	joke := getRandomChuck()
	telegram.SendTextMessage(from, " #humor "+joke)
}

func getRandomChuck() string {
	var chuckMessage chuckResponse
	httpTransport.GetAndUnmarshall(CHUCK_URI, &chuckMessage)
	return chuckMessage.TextOfJoke
}
