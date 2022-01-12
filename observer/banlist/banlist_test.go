package banlist_test

import (
	"hash/crc32"
	"testing"
	"time"

	. "fort.plus/banlist"
)

var records = NewBannedRecords("ban1")

func TestCleanup(t *testing.T) {
	var rec = NewBannedRecords("ban2")
	defer rec.Close()

	samePattern := "this is the pattern"
	err := rec.AddRecord(Item{Pattern: samePattern, ExpiredAt: time.Now()})
	time.Sleep(1 * time.Second)
	if err != nil {
		t.Errorf("Unexpected response when trying to store single pattern %s", err)
	}
}
func TestAddSingleRecords(t *testing.T) {
	records.Clear()
	samePattern := "this is the pattern"
	err := records.AddRecord(Item{Pattern: samePattern, ExpiredAt: time.Now()})
	if err != nil {
		t.Errorf("Unexpected response when trying to store single pattern %s", err)
	}
}

func TestAddSameRecords(t *testing.T) {
	records.Clear()
	samePattern := "this is the same pattern"
	records.AddRecord(Item{Pattern: samePattern, ExpiredAt: time.Now()})
	err := records.AddRecord(Item{Pattern: samePattern, ExpiredAt: time.Now()})
	if err == nil {
		t.Errorf("Unexpected response when trying to store duplicate patterns, %s", err)
	}
}

func TestDeleteRecord(t *testing.T) {
	records.Clear()
	p1 := "this is the pattern1"
	p2 := "this is the pattern2"
	records.AddRecord(Item{Pattern: p1, ExpiredAt: time.Now()})
	records.AddRecord(Item{Pattern: p2, ExpiredAt: time.Now()})
	err := records.Delete(crc32.ChecksumIEEE([]byte(p2)))
	if err != nil {
		t.Errorf("Unexpected response when trying to delete patterns, %s", err)
	}
}

func TestDeleteNonexistedRecord(t *testing.T) {
	records.Clear()
	p1 := "this is the pattern1"
	p2 := "this is the pattern2"
	p3 := "unexpected pattern"
	records.AddRecord(Item{Pattern: p1, ExpiredAt: time.Now()})
	records.AddRecord(Item{Pattern: p2, ExpiredAt: time.Now()})
	err := records.Delete(crc32.ChecksumIEEE([]byte(p3)))
	if err == nil {
		t.Errorf("Unexpected response when trying to delete non existed pattern, %s", err)
	}
}
