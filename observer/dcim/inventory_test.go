package dcim

import (
	"fmt"
	"net"
	"testing"
)

func TestDevice(t *testing.T) {
	d1 := NewDevice(1, "sw-1", net.IP{10, 1, 1, 1}, nil)
	fmt.Println(d1)

	// d2 := NewDevice()
	ip := net.ParseIP("10.1.1.1")
	fmt.Println(ip)

}

func TestRepoInMemory(t *testing.T) {
	r := &RepoInMemory{
		searchIndexes: make(map[string]*searchIndex),
	}
	if err := r.Initialize(); err != nil {
		t.Fatal(err)
	}
	var query Attributes = make(Attributes)
	query["*"] = "*"
	r.Get(query)
}
