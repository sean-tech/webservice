package webservice

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestSome(t *testing.T) {

	var fillSpeed float64 = 1.0 / float64(500)
	fillInterval := time.Duration(fillSpeed * math.Pow(10, 9))
	fmt.Println(fillInterval)
}