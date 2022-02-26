package dcim

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWhoisService(t *testing.T) {

	// path := "/var/spool/whois.txt"
	path := "/home/ag/whois.txt"

	t.Run("bad file-path", func(t *testing.T) {
		whois := NewWhoisService("/bad/path/to-file")
		whois.Run()
	})

	t.Run("load text-file", func(t *testing.T) {
		whois := NewWhoisService(path)
		whois.Run()
		log.Println(len(whois.data))
	})

	t.Run("math #1", func(t *testing.T) {
		whois := NewWhoisService(path)
		whois.Run()
		result := whois.match("10.1.1.1")
		for _, line := range result {
			fmt.Println(line)
		}
		result = whois.match("Vlan111")
		for _, line := range result {
			fmt.Println(line)
		}
	})

	t.Run("test handler", func(t *testing.T) {
		whois := NewWhoisService(path)
		whois.Run()
		q := struct {
			Query string `json:"query"`
		}{
			Query: ".*Vlan111",
		}
		rawQuery, _ := json.Marshal(q)

		resp := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/whois", bytes.NewBuffer(rawQuery))
		whois.Get(resp, req)

		fmt.Println("response code:", resp.Code)
		fmt.Println("response body:", resp.Body)
	})

	t.Run("update goroutine", func(t *testing.T) {
		whois := NewWhoisService(path)
		whois.Run()

		time.Sleep(time.Second * 6)
		whois.resetHoldTimer <- true

		time.Sleep(time.Second * 4)
		whois.resetHoldTimer <- true

		time.Sleep(time.Second * 4)
		whois.resetHoldTimer <- true

		time.Sleep(time.Minute)
	})
}
