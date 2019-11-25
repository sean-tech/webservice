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
		}
	})
	return _catch
}


type pubsubCatch  struct {
	lock       sync.RWMutex
	publishers map[string]*Publisher
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
	if publisher := this.getPublisher(topic); publisher != nil {
		publisher.lock.Lock()
		publisher.subscribers = append(publisher.subscribers, subscriber)
		publisher.lock.Unlock()
	}
}

func (this *pubsubCatch) getSubscriber(topic string) []*Subscriber {
	if publisher := this.getPublisher(topic); publisher != nil {
		publisher.lock.RLock()
		defer publisher.lock.RUnlock()
		return publisher.subscribers
	}
	return []*Subscriber{}
}