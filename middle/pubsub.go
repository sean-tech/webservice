package middle

import (
	"sync"
	"time"
)

type Publisher struct {
	topic       string        //发布者名称
	buffer      int           //订阅队列缓存大小
	timeout     time.Duration //发布超时
	lock        sync.RWMutex
	subscribers []*Subscriber //订阅者信息
}

//新建一个发布者对象，可以设置发布超时和缓存队列打长度
func NewPublisher(topic string, publishTimeout time.Duration, buffer int) *Publisher {
	if publisher := getCatch().getPublisher(topic); publisher != nil {
		return publisher
	}
	publisher := &Publisher{
		topic:   topic,
		buffer:  buffer,
		timeout: publishTimeout,
		lock:		sync.RWMutex{},
		subscribers:[]*Subscriber{},
	}
	getCatch().addPublisher(topic, publisher)
	return publisher
}

func DestoryPublisher(topic string) {
	if publisher := getCatch().getPublisher(topic); publisher != nil {
		publisher.Close()
		getCatch().deletePublisher(topic)
	}
}

//发送主题 可以容忍一定的超时
func (this *Publisher) sendTopic(sub *Subscriber, topic string, value interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case sub.Message <- value:
	case <-time.After(this.timeout):
	}
}

//发布主题
func (this *Publisher) Publish(v interface{}) {
	this.lock.RLock()
	defer this.lock.RUnlock()

	//通过waitgroup来等待collection中管道finished
	var wg sync.WaitGroup
	for _, sub := range this.subscribers { //每次发布一个消息就发送给所有的订阅者
		wg.Add(1)
		go this.sendTopic(sub, sub.Topic, v, &wg)
	}
	wg.Wait()
}

func (this *Publisher) Close() {
	this.lock.Lock()
	defer this.lock.Unlock()
	for _, sub := range this.subscribers {
		close(sub.Message)
	}
	this.subscribers = []*Subscriber{}
}

type Subscriber struct {
	Topic 	string
	Message chan interface{} //订阅者
}

// 添加一个新的订阅者，订阅过滤筛选后的主题
func SubscribeTopic(topic string) *Subscriber {
	if publisher := getCatch().getPublisher(topic); publisher != nil {
		subscriber := &Subscriber{ topic,make(chan interface{}, publisher.buffer)}
		getCatch().addSubscriber(topic, subscriber)
		return subscriber
	}
	return &Subscriber{topic, make(chan interface{}, 1)}
}

//退出订阅
func (this *Subscriber) Exit() {
	if publisher := getCatch().getPublisher(this.Topic); publisher != nil {
		i := len(publisher.subscribers) - 1
		for idx, sub := range publisher.subscribers {
			if this == sub {
				i = idx
				break
			}
		}
		publisher.lock.Lock()
		publisher.subscribers =  append(publisher.subscribers[:i], publisher.subscribers[i+1:]...)
		publisher.lock.Unlock()
	}
	close(this.Message)
}
