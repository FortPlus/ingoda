package dcim

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWhoisService(t *testing.T) {

	path := "/var/spool/whois.txt"

	t.Run("bad file-path", func(t *testing.T) {
		whois := NewWhoisService()
		if err := whois.Load("/bad/path/to-file"); err == nil {
			t.Fatal(err)
		}
		t.Log("test fails as expected")
	})

	t.Run("load text-file", func(t *testing.T) {
		whois := NewWhoisService()
		if err := whois.Load(path); err != nil {
			t.Fatal(err)
		}
		log.Println(len(whois.data))
	})

	t.Run("math #1", func(t *testing.T) {
		whois := NewWhoisService()
		if err := whois.Load(path); err != nil {
			t.Fatal(err)
		}
		result := whois.match("10.64.16.12")
		for _, line := range result {
			fmt.Println(line)
		}
		result = whois.match("Vlan1200")
		for _, line := range result {
			fmt.Println(line)
		}
	})

	t.Run("test handler", func(t *testing.T) {
		whois := NewWhoisService()
		if err := whois.Load(path); err != nil {
			t.Fatal(err)
		}
		q := struct {
			Query string `json:"query"`
		}{
			Query: ".*Vlan1200",
		}
		rawQuery, _ := json.Marshal(q)

		resp := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/whois", bytes.NewBuffer(rawQuery))
		whois.Get(resp, req)

		fmt.Println("response code:", resp.Code)
		fmt.Println("response body:", resp.Body)
	})
}
