package repository

import (
	"log"
)

type RegExComparator interface {
    IsRegExEqual(regex string) (bool, error)
}

type Callback func(RegExComparator)

//	Contain regexp as key, and function to call
var repoMap = make(map[string]Callback)

//	Store regexp pattern and callback function in repoMap
func Register(pattern string, function Callback) {
	log.Println("rep:Register:" + pattern)
	repoMap[pattern] = function
}

// Match message with repoMap keys and call corresponding function
func Call(message RegExComparator) {
	for key, function := range repoMap {
		if function != nil {
            matched, _ := message.IsRegExEqual(key)
            if matched {
                go function(message)
            }
		}
	}
}

func IsCallable(message RegExComparator) bool{
response := false
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
	for key, element := range repoMap {
		log.Println("Key:", key, "=>", "Element:", element)
	}
}

// func CallByPattern(pattern string, message parameter) {
// 	function := repoMap[pattern]
// 	if function != nil {
//  		function(message)
// 	}
// }
