package repository

import (
	"log"
	"sync"
)

type RegExComparator interface {
	IsRegExEqual(regex string) (bool, error)
}

type Callback func(RegExComparator)

//	Contain regexp as key, and function to call
var repoMap = make(map[string]Callback)
var lock sync.RWMutex = sync.RWMutex{}

//	Store regexp pattern and callback function in repoMap
func Register(pattern string, function Callback) {
	log.Println("rep:Register:" + pattern)
	lock.Lock()
	defer lock.Unlock()
	repoMap[pattern] = function
}

// Match message with repoMap keys and call corresponding function
func Call(message RegExComparator) {
	lock.RLock()
	defer lock.RUnlock()

	for key, function := range repoMap {
		if function != nil {
			matched, _ := message.IsRegExEqual(key)
			if matched {
				go function(message)
			}
		}
	}
}

func IsCallable(message RegExComparator) bool {
	response := false
	lock.RLock()
	defer lock.RUnlock()

	for key, function := range repoMap {
		if function != nil {
			matched, _ := message.IsRegExEqual(key)
			if matched {
				response = true
			}
		}
	}
	return response
}

func Show() {
	lock.RLock()
	defer lock.RUnlock()

	for key, element := range repoMap {
		log.Println("Key:", key, "=>", "Element:", element)
	}
}

func Clear() {
	lock.Lock()
	defer lock.Unlock()
	repoMap = make(map[string]Callback)
}
