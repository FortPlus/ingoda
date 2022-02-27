package whois

import (
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {
	client := &dcimBotClient{
		serverUri: "http://localhost:38000",
	}
	fmt.Println(client)
	q := query{Query: "999"}
	response := client.fromService(q)
	fmt.Println(response)
}
