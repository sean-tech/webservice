package service

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestGenerateId(t *testing.T) {
	fmt.Println(time.Now().Unix())

	var testData = []struct{
		workerId int64
	} {
		{0},
		{1},
		{2},
		{3},
		{4},
	}
	for _, data := range testData {
		id, err := GenerateId(data.workerId)
		if err != nil {
			t.Errorf(err.Error())
		}
		fmt.Println(id)
	}
}

func TestGenerateId2(t *testing.T) {
	var ids []int64 = []int64{}
	var lock sync.Mutex
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id, err := GenerateId(1)
			if err != nil {
				t.Error(err)
			}
			lock.Lock()
			defer lock.Unlock()
			ids = append(ids, id)
		}()
	}
	wg.Wait()
	fmt.Printf("%v", ids)
}
