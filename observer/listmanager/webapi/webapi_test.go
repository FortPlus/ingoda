package webapi

import (
	"fmt"
	"testing"
	"time"

	banlist "fort.plus/listmanager"
)

func TestSerialize(t *testing.T) {

	lm := NewListManager()
	fmt.Println(lm)

	lm.Storage["list1"] = banlist.New("list1")
	list1 := lm.Storage["list1"]
	list1.AddRecord(banlist.Item{Pattern: "pattern1", ExpiredAt: time.Now().Add(time.Minute * 2)})
	list1.AddRecord(banlist.Item{Pattern: "pattern2", ExpiredAt: time.Now().Add(time.Minute * 3)})

	list2 := banlist.New("list2")
	lm.Storage["list2"] = list2
	list2.AddRecord(banlist.Item{Pattern: "pattern1", ExpiredAt: time.Now().Add(time.Minute * 4)})
	list2.AddRecord(banlist.Item{Pattern: "pattern2", ExpiredAt: time.Now().Add(time.Minute * 5)})

	t.Run("Serialze", func(t *testing.T) {

		data, err := lm.Serialize()
		fmt.Println(err)
		fmt.Println(string(data))
	})

	t.Run("Deserialize", func(t *testing.T) {
		data, err := lm.Serialize()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("\n\n---Unmarshal")
		lm2, err := Deserialize(data)
		if err != nil {
			t.Fatal(err)
		}
		// json.Unmarshal(data, &lm2)
		fmt.Println(lm2)

		for _, value := range lm2.Storage {
			// value.GetPatterns()
			fmt.Println("pattern", value.GetPatterns())
		}
	})
}
