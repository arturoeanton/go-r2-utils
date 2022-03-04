package observer

import "sync"

type Observable interface {
	AddObserver(observer Observer)
	RemoveObserver(observer Observer)
	Notify()
	ChangeState(interface{})
	GetState() interface{}
}

type BasicObservable struct {
	observers map[string]Observer
	state     interface{}
}

func (e *BasicObservable) AddObserver(observer Observer) {
	if len(e.observers) == 0 {
		e.observers = make(map[string]Observer)
	}
	e.observers[observer.GetID()] = observer
}
func (e *BasicObservable) RemoveObserver(observer Observer) {
	_, ok := e.observers[observer.GetID()]
	if !ok {
		return
	}
	delete(e.observers, observer.GetID())
}
func (e *BasicObservable) Notify() {
	wg := sync.WaitGroup{}
	for _, v := range e.observers {
		wg.Add(1)
		go func(wg *sync.WaitGroup, v Observer) {
			v.Notify(e.state)
			wg.Done()
		}(&wg, v)
	}
	wg.Wait()
}

func (e *BasicObservable) ChangeState(data interface{}) {
	e.state = data
	e.Notify()
}

func (e *BasicObservable) GetState() interface{} {
	return e.state
}
