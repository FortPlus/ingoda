package chuck

import (
	"log"

	"fort.plus/im"
	"fort.plus/repository"
	httpTransport "fort.plus/transport"
)

const (
	CHUCK_URI = "https://api.chucknorris.io/jokes/random"
)

type chuckResponse struct {
	TextOfJoke string `json:"value"`
}

var m im.Carrier

func Register(t im.Carrier) {
	m = t
	repository.Register(".*[Cc]huck.*", a)
}

var a = func(message repository.RegExComparator) {
	telegramMessage := im.Cast(message)
	log.Printf("any:subscriber function, message is:%s", telegramMessage.Text)
	joke := getRandomChuck()
	m.Send(telegramMessage.From, " #humor "+joke)

}

func getRandomChuck() string {
	var chuckMessage chuckResponse
	httpTransport.GetAndUnmarshall(CHUCK_URI, &chuckMessage)
	return chuckMessage.TextOfJoke
}
