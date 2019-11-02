package service

import (
	"fmt"
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
