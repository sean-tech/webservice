package middle

import (
	"sync"
	"time"
)

type Publisher struct {
	topic       string        //发布者名称
	timeout     time.Duration //发布超时
	lock        sync.RWMutex
}

//新建一个发布者对象，可以设置发布超时和缓存队列打长度
func NewPublisher(topic string, publishTimeout time.Duration) *Publisher {
	if publisher := getCatch().getPublisher(topic); publisher != nil {
		return publisher
	}
	publisher := &Publisher{
		topic:   topic,
		timeout: publishTimeout,
		lock:		sync.RWMutex{},
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
	subscribers := getCatch().getSubscribers(this.topic)
	for _, sub := range subscribers { //每次发布一个消息就发送给所有的订阅者
		wg.Add(1)
		go this.sendTopic(sub, sub.Topic, v, &wg)
	}
	wg.Wait()
}

func (this *Publisher) Close() {
	this.lock.Lock()
	defer this.lock.Unlock()
	subscribers := getCatch().getSubscribers(this.topic)
	for _, sub := range subscribers {
		close(sub.Message)
	}
	getCatch().deleteAllSubscribers(this.topic)
}

type Subscriber struct {
	Topic 	string
	Message chan interface{} //订阅者
}

// 添加一个新的订阅者，订阅过滤筛选后的主题
func SubscribeTopic(topic string, buffer int) *Subscriber {
	subscriber := &Subscriber{ topic,make(chan interface{}, buffer)}
	getCatch().addSubscriber(topic, subscriber)
	return subscriber
}

//退出订阅
func (this *Subscriber) Exit() {
	subscribers := getCatch().getSubscribers(this.Topic)
	for _, sub := range subscribers {
		if sub == this {
			close(this.Message)
			getCatch().deleteSubscriber(this)
		}
	}
}
