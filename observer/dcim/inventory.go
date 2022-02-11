package dcim

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

func StrToIP(ip string) (net.IP, error) {
	octets := strings.Split(ip, ".")
	var res [4]byte
	for index, octet := range octets {
		if index > 3 {
			return nil, errors.New("main:StrToIP bad ip")
		}
		b, err := strconv.Atoi(octet)
		if err != nil {
			// TODO: implement
			return nil, err
		}
		res[index] = byte(b)
	}
	return net.IP{res[0], res[1], res[2], res[3]}, nil
}

type DeviceRepository interface {
	Initialize() error
	Delete() error
	Add() error
	Get(DeviceQuery) ([]Device, error)
	GetAll() error
}

//
type searchIndex struct {
	count   uint64
	bits    []byte
	indexes []uint64
}

// InMemory storage of devices, just a JSON file
type RepoInMemory struct {
	sync.RWMutex
	Devices       []Device
	searchIndexes map[string]*searchIndex
}

func NewRepoInMemory() *RepoInMemory {
	return &RepoInMemory{
		searchIndexes: make(map[string]*searchIndex),
	}
}

func (r *RepoInMemory) Initialize() error {
	var (
		file *os.File
		err  error
		path string = os.Getenv("DCIM_DB")
	)
	if file, err = os.Open(path); err != nil {
		return err
	}
	defer file.Close()
	devices := make([]DeviceMarshal, 1024)
	if err := json.NewDecoder(file).Decode(&devices); err != nil {
		return err
	}
	if err := r.build(devices); err != nil {
		return err
	}
	r.buildIndex()
	return nil
}

func (r *RepoInMemory) Delete() error {
	r.Lock()
	defer r.Unlock()

	return nil
}
func (r *RepoInMemory) Add() error {
	r.Lock()
	defer r.Unlock()
	return nil
}
func (r *RepoInMemory) Get(q DeviceQuery) ([]Device, error) {
	r.RLock()
	defer r.RUnlock()

	// (1) collect searchIndexes
	var searchIndexes []*searchIndex
	// fmt.Println(searchIndexes)
	for _, query := range q.attributes {
		key := fmt.Sprintf("%v:%v", query.name, query.value)
		// fmt.Println(key)
		if searcher, ok := r.searchIndexes[key]; ok {
			searchIndexes = append(searchIndexes, searcher)
		}
	}
	// (1.1) If nothing
	if searchIndexes == nil {
		// fmt.Println("No results founded for this query.")
		return nil, nil
	}

	// fmt.Println("find select search indexes")
	// fmt.Println(searchIndexes)

	// (2) find shorter searhIndexes
	// TODO: create search-index-all for getting all devices
	var shorter *searchIndex = searchIndexes[0]

	for _, searcher := range searchIndexes {
		if searcher.count < shorter.count {
			shorter = searcher
		}
	}

	// (3) find value that match all query
	match := make([]Device, 0)
	for _, value := range shorter.indexes {
		isOk := byte(1)
		for _, searcher := range searchIndexes {
			isOk &= searcher.bits[value]
		}
		if isOk == 1 {
			match = append(match, r.Devices[value])
		}
	}

	// fmt.Println("--- match")
	// fmt.Println(match)

	return match, nil
}

func (r *RepoInMemory) GetAll() error {
	r.RLock()
	defer r.RUnlock()
	return nil
}

func (r *RepoInMemory) build(devices []DeviceMarshal) error {

	r.Devices = make([]Device, 0, len(devices))

	for index, device := range devices {
		ip, err := StrToIP(device.Ip)
		if err != nil {
			return err
		}
		newDev := NewDevice(uint64(index), device.Name, ip, nil)
		if device.Attr != nil {
			newDev.Attributes = make([]Attribute, 0, len(device.Attr))
			for key, value := range device.Attr {
				newDev.Attributes = append(
					newDev.Attributes,
					Attribute{name: key, value: value},
				)
			}
		}
		r.Devices = append(r.Devices, *newDev)
	}
	return nil
}

func (r *RepoInMemory) buildIndex() {

	total := len(r.Devices)

	allIndex := &searchIndex{
		count:   uint64(total),
		bits:    make([]byte, total),
		indexes: make([]uint64, total),
	}

	for index, dev := range r.Devices {
		for _, attr := range dev.Attributes {
			key := fmt.Sprintf("%v:%v", attr.name, attr.value)
			if _, ok := r.searchIndexes[key]; !ok {
				r.searchIndexes[key] = &searchIndex{
					count:   0,
					bits:    make([]byte, total),
					indexes: make([]uint64, 0),
				}
			}
			search := r.searchIndexes[key]
			search.count += 1
			search.bits[r.Devices[index].id] = 1
			search.indexes = append(search.indexes, uint64(index))
		}

		// mark index
		allIndex.bits[index] = 1
		allIndex.indexes[index] = uint64(index)
	}

	// add * searchIndex
	r.searchIndexes["*:*"] = allIndex

	// for index, value := range r.searchIndexes {
	// 	fmt.Println(index, "\t", value)
	// }
}

type DeviceQuery struct {
	attributes []Attribute
}

func NewDeviceQuery(raw map[string]string) *DeviceQuery {
	if raw == nil {
		// get all
		return &DeviceQuery{
			[]Attribute{
				NewAttribute("*", "*"),
			},
		}
	}
	query := &DeviceQuery{
		make([]Attribute, 0, len(raw)),
	}
	for key, value := range raw {
		query.attributes = append(query.attributes, NewAttribute(key, value))
	}
	return query
}

func NewDeviceAllQuery() DeviceQuery {
	query := DeviceQuery{
		make([]Attribute, 1),
	}
	query.attributes[0] = NewAttribute("*", "*")
	return query
}

// DeviceShot used for API-response
type DeviceShort struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
}

// DeviceMarshal used for unmarshalling from JSON
type DeviceMarshal struct {
	Name string            `json:"name"`
	Ip   string            `json:"ip"`
	Attr map[string]string `json:"attributes"`
}

// Device represent managed network element
type Device struct {
	id         uint64
	Name       string
	OamIP      net.IP
	Attributes []Attribute
}

// NewDevice create one new
func NewDevice(id uint64, name string, ip net.IP, attr []Attribute) *Device {
	return &Device{
		id:         id,
		Name:       name,
		OamIP:      ip,
		Attributes: attr,
	}
}

// Attribute in free format
type Attribute struct {
	name  string
	value string
}

func (a Attribute) String() string {
	return fmt.Sprintf("%v:%v", a.name, a.value)
}

func NewAttribute(name, value string) Attribute {
	return Attribute{name: name, value: value}
}

func createResponse(devices []Device) ([]byte, error) {
	var result []DeviceShort = make([]DeviceShort, 0, len(devices))
	for _, device := range devices {
		dev := DeviceShort{
			Name: device.Name,
			Ip:   device.OamIP.String(),
		}
		result = append(result, dev)
	}
	raw, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func parseQuery(r *http.Request) (*DeviceQuery, error) {
	var attr map[string]string

	if r.ContentLength == 0 {
		return NewDeviceQuery(attr), nil
	}

	if err := json.NewDecoder(r.Body).Decode(&attr); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return NewDeviceQuery(attr), nil
}
