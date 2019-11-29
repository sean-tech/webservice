package middle

import (
	"fmt"
	"testing"
	"time"
)

func TestPublisher_Publish(t *testing.T) {
	topic1 := "token"
	topic2 := "some"

	s1 := SubscribeTopic(topic1, 1000)
	s2 := SubscribeTopic(topic2, 1000)

	p1 := NewPublisher(topic1, 10*time.Second)
	defer p1.Close()
	p2 := NewPublisher(topic2, 10*time.Second)
	defer p2.Close()



	p1.Publish("hello, world")
	p1.Publish("hello, golang")
	p2.Publish("hello, p2")
	s2.Exit()
	p2.Publish("hello, p22")
	p1.Close()
	p2.Close()
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