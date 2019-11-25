package middle

import (
	"fmt"
	"testing"
	"time"
)

func TestPublisher_Publish(t *testing.T) {
	topic1 := "token"
	topic2 := "some"
	p1 := NewPublisher(topic1, 10*time.Second, 1000)
	defer p1.Close()
	p2 := NewPublisher(topic2, 10*time.Second, 1000)
	defer p2.Close()

	s1 := SubscribeTopic(topic1)
	s2 := SubscribeTopic(topic2)

	p1.Publish("hello, world")
	p1.Publish("hello, golang")
	p2.Publish("hello, p2")
	p2.Publish("hello, p22")
	go func() {
		for Message := range s1.Message {
			fmt.Println("token ", Message)
		}
	}()
	go func() {
		for Message := range s2.Message {
			fmt.Println("some ", Message)
		}
	}()

	time.Sleep(3 * time.Second)
}