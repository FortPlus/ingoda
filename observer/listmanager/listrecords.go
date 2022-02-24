package listmanager

import (
	"context"
	"encoding/json"
	"hash/crc32"
	"log"
	"regexp"
	"sync"
	"time"

	"fort.plus/fperror"
	httpTransport "fort.plus/transport"
)

const CLEANUP_PERIOD = 60

type Item struct {
	Pattern   string    `json:"pattern"`
	ExpiredAt time.Time `json:"expired_at"`
}

type ItemsMap map[uint32]Item

func (i *ItemsMap) UnmarshallJSON(data []byte) error {
	err := json.Unmarshal(data, i)
	if err != nil {
		err = fperror.Warning("can't unmarshall JSON", err)
	}
	return err
}
func (i *ItemsMap) MarshallJSON(data []byte) ([]byte, error) {
	response, err := json.Marshal(i)
	if err != nil {
		err = fperror.Warning("can't marshall JSON", err)
	}
	return response, err
}

type ListRecords struct {
	items  ItemsMap
	name   string
	lock   sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

// Create new banned records map
func New(name string) *ListRecords {
	b := &ListRecords{
		items: make(ItemsMap),
		lock:  sync.RWMutex{},
		name:  name,
	}

	b.ctx, b.cancel = context.WithCancel(context.TODO())
	go b.cleanJob()
	return b
}

func (b *ListRecords) cleanJob() {
	println("run cleanJob")
	for {
		select {
		case <-b.ctx.Done():
			log.Println(b.name, "banlist::cleanExpired() - context is done")
			return
		case <-time.After(CLEANUP_PERIOD * time.Second):
			b.cleanExpired()
		}
	}
}

//
// Periodically fetch list from the server and replace the local with it
//
func (b *ListRecords) PeriodicImportFromServer(listManagerUri string, period int) {
	for {
		listFromServer := make(ItemsMap)
		select {
		case <-b.ctx.Done():
			log.Println(b.name, "banlist::PeriodicImportFromServer() - context is done")
			return
		default:
			err := httpTransport.Get(listManagerUri+b.name, &listFromServer)
			if err == nil {
				b.lock.Lock()
				b.items = listFromServer
				b.lock.Unlock()
			} else {
				log.Println(b.name, ", error while import list from server", err)
			}
			time.Sleep(time.Duration(period) * time.Second)
		}
	}
}

// Prepare banned records service to stop cleanup and storage update goroutines
func (b *ListRecords) IsEmpty() bool {
	b.lock.Lock()
	defer b.lock.Unlock()
	return len(b.items) == 0
}

// Prepare banned records service to stop cleanup and storage update goroutines
func (b *ListRecords) Close() {
	log.Println(b.name, "banlist::Close()")
	b.cancel()
}

func (b *ListRecords) cleanExpired() {
	log.Println(b.name, "banlist::CleanExpired()")
	currentTime := time.Now()
	b.lock.Lock()
	defer b.lock.Unlock()
	for id, element := range b.items {
		if !element.ExpiredAt.IsZero() && currentTime.Sub(element.ExpiredAt) > 0 {
			delete(b.items, id)
		}
	}
}

// Delete all patterns from banned records
func (b *ListRecords) Clear() {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.items = make(ItemsMap)
}

// Get all records in JSON format
func (b *ListRecords) GetRecords() ([]byte, error) {
	var err error

	b.lock.RLock()
	defer b.lock.RUnlock()

	json_data, err := json.Marshal(b.items)
	if err != nil {
		err = fperror.Warning("can't encode records to json", err)
	}
	return json_data, err
}

// Add item to banned records
func (b *ListRecords) AddRecord(item Item) error {
	var result error
	id := crc32.ChecksumIEEE([]byte(item.Pattern))
	b.lock.Lock()
	defer b.lock.Unlock()

	if _, ok := b.items[id]; ok {
		result = fperror.Warning("duplicate pattern", nil)
	} else {
		b.items[id] = item
	}
	return result
}

// Check if message is banned with regexp pattern matching
func (b *ListRecords) CheckIfExists(msg string) bool {
	b.lock.RLock()
	defer b.lock.RUnlock()
	for _, element := range b.items {
		res, _ := regexp.MatchString(element.Pattern, msg)
		if res {
			return true
		}
	}
	return false
}

// Delete pattern from banned records using crc32 pattern ID
func (b *ListRecords) Delete(id uint32) error {
	var result error

	b.lock.Lock()
	defer b.lock.Unlock()

	if _, ok := b.items[id]; ok {
		delete(b.items, id)
	} else {
		result = fperror.Warning("Can't find pattern to delete", nil)
	}
	return result
}

// Get slice with patterns
func (b *ListRecords) GetPatterns() []string {
	var result []string = []string{}
	b.lock.Lock()
	defer b.lock.Unlock()
	for _, element := range b.items {
		result = append(result, element.Pattern)
	}
	return result
}

func (b *ListRecords) MarshalJSON() ([]byte, error) {
	b.lock.Lock()
	defer b.lock.Unlock()
	var a struct {
		Items []Item `json:"items`
		Name  string `json:"name"`
	}
	for _, value := range b.items {
		a.Items = append(a.Items, value)
	}
	a.Name = b.name
	raw, err := json.Marshal(a)
	return raw, err
}

func (b *ListRecords) UnmarshalJSON(data []byte) error {

	// Create new instace with initialize and start goroutine
	*b = *New("")
	b.lock.Lock()
	b.lock.Unlock()

	var a struct {
		Items []Item `json:"items`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(data, &a); err != nil {
		return nil
	}
	for _, value := range a.Items {
		b.AddRecord(value)
	}
	return nil
}
