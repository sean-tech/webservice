package middle

import "sync"


var (
	_catchOnce sync.Once
	_catch *pubsubCatch
)

func getCatch() *pubsubCatch {
	_catchOnce.Do(func() {
		_catch = &pubsubCatch{
			lock:       sync.RWMutex{},
			publishers: make(map[string]*Publisher),
			subscribers:make(map[string][]*Subscriber),
		}
	})
	return _catch
}


type pubsubCatch  struct {
	lock       	sync.RWMutex
	publishers 	map[string]*Publisher
	subscribers map[string][]*Subscriber //订阅者信息
}

func (this *pubsubCatch) addPublisher(topic string, publisher *Publisher)  {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.publishers[topic] = publisher
}

func (this *pubsubCatch) deletePublisher(topic string)  {
	publisher := this.getPublisher(topic)
	if publisher != nil {
		this.lock.Lock()
		defer this.lock.Unlock()
		delete(this.publishers, topic)
	}
}

func (this *pubsubCatch) getPublisher(topic string) *Publisher {
	this.lock.RLock()
	defer this.lock.RUnlock()
	publisher, ok := this.publishers[topic]
	if ok {
		return publisher
	}
	return nil
}

func (this *pubsubCatch) addSubscriber(topic string, subscriber *Subscriber)  {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.subscribers[topic] = append(this.subscribers[topic], subscriber)
}

func (this *pubsubCatch) getSubscribers(topic string) []*Subscriber {
	this.lock.RLock()
	defer this.lock.RUnlock()
	subscribers, ok := this.subscribers[topic]
	if ok {
		return subscribers
	}
	return []*Subscriber{}
}

func (this *pubsubCatch) deleteAllSubscribers(topic string)  {
	subscribers := this.getSubscribers(topic)
	if len(subscribers) > 0 {
		this.lock.Lock()
		defer this.lock.Unlock()
		delete(this.subscribers, topic)
	}
}

func (this *pubsubCatch) deleteSubscriber(subscriber *Subscriber) {
	subscribers := this.getSubscribers(subscriber.Topic)
	if len(subscribers) > 0 {
		i := len(subscribers) - 1
		for idx, sub := range subscribers {
			if subscriber == sub {
				i = idx
				break
			}
		}
		this.subscribers[subscriber.Topic] =  append(subscribers[:i], subscribers[i+1:]...)
	}
}