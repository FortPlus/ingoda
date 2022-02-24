package whois

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"fort.plus/im"
	"fort.plus/repository"
	httpTransport "fort.plus/transport"
)

var (
	carrier im.Carrier
)

type query struct {
	Query string `json:"query"`
}

type dcimBotClient struct {
	serverUri string
	carrier   im.Carrier
}

func Register(carrier im.Carrier, serverUri string) {
	b := &dcimBotClient{
		serverUri: serverUri,
		carrier:   carrier,
	}
	repository.Register("/whois .*", b.get)
}

func (b *dcimBotClient) get(message repository.RegExComparator) {
	msg := im.Cast(message)
	re := regexp.MustCompile("/whois (.*)")
	match := re.FindStringSubmatch(msg.Text)

	// build query
	q := query{match[0]}

	// send request to service
	raw, err := json.Marshal(q)
	if err != nil {
		b.carrier.Send(msg.From, fmt.Sprintf("can' marshal JSON query, %s", err))
		return
	}

	resp, err := httpTransport.Request(http.MethodGet, b.serverUri+"/api/v1/whois", bytes.NewBuffer(raw))
	if err != nil {
		b.carrier.Send(msg.From, fmt.Sprintf("failed then request dcim/whois service, %s", err))
		return
	}

	// parse response
	rawResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		b.carrier.Send(msg.From, fmt.Sprintf("failed then read bytes from dcim/whois service response, %s", err))
		return
	}

	var response []string
	if err := json.Unmarshal(rawResponse, &response); err != nil {
		b.carrier.Send(msg.From, fmt.Sprintf("failed then parse JSON from dcim/whois service response, %s", err))
		return
	}

	b.carrier.Send(msg.From, strings.Join(response, "\n"))
}
