package dcim

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"sync"

	"fort.plus/fperror"
)

type Storing interface {
	Initialize() error
	Delete() error
	Add() error
	Get(Attributes) ([]Device, error)
}

type searchIndex struct {
	count   uint64
	bits    []byte
	indexes []uint64
}

// RepoInMemory storage of devices, just a JSON file
type RepoInMemory struct {
	sync.RWMutex
	Devices       []Device
	searchIndexes map[string]*searchIndex
	path          string
}

func NewRepoInMemory(path string) *RepoInMemory {
	return &RepoInMemory{
		searchIndexes: make(map[string]*searchIndex),
		path:          path,
	}
}

func (r *RepoInMemory) Initialize() error {
	r.Lock()
	defer r.Unlock()

	var (
		err error
		b   []byte
	)

	if b, err = ioutil.ReadFile(r.path); err != nil {
		return fperror.Critical("RepoInMemory::Initialize Cannot read JSON DCIM-db file", err)
	}

	if err = json.Unmarshal(b, &r.Devices); err != nil {
		return fperror.Critical("RepoInMemory::Initialize Cannot parse JSON DCIM-db file", err)
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
func (r *RepoInMemory) Get(q Attributes) ([]Device, error) {
	r.RLock()
	defer r.RUnlock()

	// (1) collect searchIndexes
	var searchIndexes []*searchIndex
	for attrName, attrValue := range q {
		key := attrName.Key(attrValue)
		if _, ok := r.searchIndexes[key]; !ok {
			return nil, fperror.Warning("RepoInMemory::Get no result for query: "+key, nil)
		}

		searchIndexes = append(searchIndexes, r.searchIndexes[key])
	}

	// (1.1) case when client send query like '{}'
	if searchIndexes == nil {
		return nil, fperror.Warning("RepoInMemory::Get no result founded", nil)
	}

	// (2) find shorter searhIndexes
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

	return match, nil
}

func (r *RepoInMemory) build(devices []DeviceMarshal) {

	// 	r.Devices = make([]Device, 0, len(devices))

	// 	for index, device := range devices {
	// 		ip := net.ParseIP(device.Ip)
	// 		newDev := NewDevice(uint64(index), device.Name, ip, nil)
	// 		if device.Attr != nil {
	// 			newDev.Attributes = make([]Attribute, 0, len(device.Attr))
	// 			for key, value := range device.Attr {
	// 				newDev.Attributes = append(
	// 					newDev.Attributes,
	// 					Attribute{Name: key, Value: value},
	// 				)
	// 			}
	// 		}
	// 		r.Devices = append(r.Devices, *newDev)
	// 	}
}

func (r *RepoInMemory) buildIndex() {

	total := len(r.Devices)

	allIndex := &searchIndex{
		count:   uint64(total),
		bits:    make([]byte, total),
		indexes: make([]uint64, total),
	}

	for index, dev := range r.Devices {
		r.Devices[index].id = uint64(index)
		for attrName, attrValue := range dev.Attrs {
			key := attrName.Key(attrValue)
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

type AttributeName string
type Attr string
type Attributes map[AttributeName]Attr

func (a AttributeName) Key(value Attr) string {
	return string(a) + ":" + string(value)
}

// Device represent managed network element
type Device struct {
	id    uint64     `json:"omitempty"`
	Name  string     `json:"name"`
	OamIP net.IP     `json:"ip"`
	Attrs Attributes `json:"attributes"`
}

// NewDevice create one new
func NewDevice(id uint64, name string, ip net.IP, attr Attributes) *Device {
	return &Device{
		id:    id,
		Name:  name,
		OamIP: ip,
		Attrs: attr,
	}
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
		return nil, fperror.Warning("createResponse:: cannot marshal to Json", err)
	}
	return raw, nil
}

func parseQuery(r *http.Request) (Attributes, error) {
	var attr Attributes

	if r.ContentLength == 0 {
		attr = make(Attributes)
		attr["*"] = "*"
		return attr, nil
	}

	if err := json.NewDecoder(r.Body).Decode(&attr); err != nil {
		return nil, fperror.Warning("parseQuery:: cannot parse query JSON in request", err)
	}
	defer r.Body.Close()

	return attr, nil
}
