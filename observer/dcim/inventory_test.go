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
	ip, _ := StrToIP("10.1.1.1")
	fmt.Println(ip)

}

func TestAttribute(t *testing.T) {
	attr1 := NewAttribute("site", "msk")
	attr2 := NewAttribute("platform", "ios")
	attr3 := NewAttribute("vendor", "cisco")
	attr4 := NewAttribute("vendor", "access-sw")
	fmt.Println(attr1, attr2, attr3, attr4)
}

func TestRepoInMemory(t *testing.T) {
	r := &RepoInMemory{
		searchIndexes: make(map[string]*searchIndex),
	}
	if err := r.Initialize(); err != nil {
		t.Fatal(err)
	}

	q := DeviceQuery{
		attributes: []Attribute{
			{
				name:  "*",
				value: "*",
			},
			// {
			// 	name:  "role",
			// 	value: "access",
			// },
		},
	}
	r.Get(q)
}
