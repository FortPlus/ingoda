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
	q := query{match[0]}
	response := b.fromService(q)
	b.carrier.Send(msg.From, response)
}

// fromService send http-request to whois-service and parse response
func (b *dcimBotClient) fromService(q query) string {
	// masrhal request to json-bytes
	raw, err := json.Marshal(q)
	if err != nil {
		return fmt.Sprintf("can' marshal JSON query, %s", err)
	}

	// send request to service
	resp, err := httpTransport.Request(http.MethodGet, b.serverUri+"/api/v1/whois", bytes.NewBuffer(raw))
	if err != nil {
		return fmt.Sprintf("failed then request dcim/whois service, %s", err)
	}

	if resp.StatusCode == http.StatusNoContent {
		return fmt.Sprintf("get status code: 204 (no content) from whois-service for query:%v", q.Query)
	}
	if resp.StatusCode == http.StatusBadRequest {
		return fmt.Sprintf("get status code: 404 (bad request) from whois-service for query:%v", q.Query)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		return fmt.Sprintf("get status code: 500 from whois-service for query:%v", q.Query)
	}

	// parse response
	rawResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("failed then read bytes from dcim/whois service response, %s", err)
	}

	var response []string
	if err := json.Unmarshal(rawResponse, &response); err != nil {
		return fmt.Sprintf("failed then parse JSON from dcim/whois service response, %s", err)
	}

	return strings.Join(response, "\n")
}
