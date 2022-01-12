package banlist

import (
	"context"
	"encoding/json"
	"hash/crc32"
	"log"
	"regexp"
	"sync"
	"time"

	"fort.plus/fperror"
)

const CLEANUP_PERIOD = 60

type Item struct {
	Pattern   string    `json:"pattern"`
	ExpiredAt time.Time `json:"expired_at"`
}

type ItemsMap map[uint32]Item

type BannedRecords struct {
	items  ItemsMap
	name   string
	lock   sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

// Create new banned records map
func NewBannedRecords(name string) *BannedRecords {
	var b BannedRecords
	b.items = make(ItemsMap)
	b.lock = sync.RWMutex{}
	b.name = name
	b.ctx, b.cancel = context.WithCancel(context.TODO())
	go cleanExpired(b.ctx, &b)
	return &b
}

func cleanExpired(ctx context.Context, b *BannedRecords) {
	log.Println("cleanExpired(ctx,", b.name, ")")
	go func() {
		for {
			log.Println(b.name, "time to start cleanup")
			b.cleanExpired()
			time.Sleep(CLEANUP_PERIOD * time.Second)
		}
	}()
	<-ctx.Done()
	log.Println(b.name, "banlist::cleanExpired() - done")
}

// Prepare banned records service to stop cleanup and storage update goroutines
func (b *BannedRecords) Close() {
	log.Println(b.name, "banlist::Close()")
	b.cancel()
}

func (b *BannedRecords) cleanExpired() {
	log.Println(b.name, "banlist::CleanExpired()")
	currentTime := time.Now()
	b.lock.Lock()
	defer b.lock.Unlock()
	for id, element := range b.items {
		if currentTime.Sub(element.ExpiredAt) > 0 {
			delete(b.items, id)
		}
	}
}

// Delete all patterns from banned records
func (b *BannedRecords) Clear() {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.items = make(ItemsMap)
}

// Get all records in JSON format
func (b *BannedRecords) GetRecords() ([]byte, error) {
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
func (b *BannedRecords) AddRecord(item Item) error {
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
func (b *BannedRecords) CheckIfExists(msg string) bool {
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
func (b *BannedRecords) Delete(id uint32) error {
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
