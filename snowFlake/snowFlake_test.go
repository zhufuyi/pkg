package snowFlake

import (
	"fmt"
	"testing"
)

func TestNewID(t *testing.T) {
	max := 100
	uidsChan := make(chan int64)
	uids := []int64{}

	InitSnowFlake(1)
	for i := 0; i < max; i++ {
		go func(uidsChan chan int64) {
			id := NewID()
			uidsChan <- id
		}(uidsChan)
	}

	for i := 0; i < max; i++ {
		uids = append(uids, <-uidsChan)
	}

	expect := max
	actual := RemoveDuplicate(uids)
	if actual != expect {
		t.Errorf("没有生成唯一的ID! len: %d", actual)
	}

	fmt.Printf("actual=%d, expect=%d\nuids=%v", actual, expect, uids)
}

func RemoveDuplicate(slice []int64) int {
	for i := 0; i < len(slice); i++ {
		for j := i + 1; j < len(slice); {
			if slice[i] == slice[j] {
				slice = append(slice[:j], slice[j+1:]...)
			} else {
				j++
			}
		}
	}

	return len(slice)
}
