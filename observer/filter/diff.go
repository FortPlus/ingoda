package diff

import (
	"log"
	"sync"
	"time"

	"github.com/agnivade/levenshtein"
)

const (
	EMPTY_STRING = ""
)

type notificationMessage struct {
	InsertTime time.Time
	Text       string
}

// settings
var sampleLifetimeSeconds float64 = 5 * 60
var numberOfSamples = 20

var samplesMutex sync.Mutex
var samples = make([]notificationMessage, numberOfSamples)
var currentIndex int = 0
var foundSamplesToCompare bool = false

func Reset() {
	currentIndex = 0
	foundSamplesToCompare = false
	samples = make([]notificationMessage, numberOfSamples)

}

//
//	Check how much message is different from previous, collected to samples
//
//	return True -	if no samples to compare with
//					if difference betwee message and a sample is less then threshold
//
func IsThresholdExceeded(message string, threshold int) bool {
	ld := GetLevenshteinDistance(message)
	log.Printf("LD:%d", ld)
	if !foundSamplesToCompare {
		return true
	}

	//return true if no samples to compare with, or threshold not exceeded
	if ld < threshold {
		return true
	} else {
		log.Printf("threshold exceeded, LD:%d", ld)
		return false
	}
}

func GetLevenshteinDistance(message string) int {
	samplesMutex.Lock()
	defer samplesMutex.Unlock()

	cleanExpiredSamples()
	minDistance := calcMinimumDistance(message)
	setSample(currentIndex, message) //TODO: rewrite magic with currentIndex
	adjustCurrentIndex()
	return minDistance
}

func calcMinimumDistance(message string) int {
	var minDistance int = -1
	foundSamplesToCompare = false

	for _, element := range samples {
		if len(element.Text) > 0 {
			foundSamplesToCompare = true
			distance := levenshtein.ComputeDistance(element.Text, message)
			//in first comparsion any resut override default value
			if minDistance == -1 {
				minDistance = distance
				continue
			}
			if distance < minDistance {
				minDistance = distance
			}
		}
	}
	return minDistance
}

func adjustCurrentIndex() {
	if currentIndex < len(samples)-1 {
		currentIndex++
	} else {
		currentIndex = 0
	}
}

func cleanExpiredSamples() {
	timeNow := time.Now()
	for index, element := range samples {
		dur := timeNow.Sub(element.InsertTime)
		if len(element.Text) > 0 && dur.Seconds() > sampleLifetimeSeconds {
			setSample(index, EMPTY_STRING)
		}
	}
}

func setSample(index int, text string) {
	samples[index].Text = text
	samples[index].InsertTime = time.Now()
}
