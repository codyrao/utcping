package utils

import (
	"fmt"
	"sync"
)

var (
	once          sync.Once
	globalRoutine *Routine
)

type Routine struct {
	queue         chan *Event
	callbackQueue chan *CallbackEvent
}

type Callback struct {
	callbackResQueue chan *CallbackRes
}

func GetCallback() *Callback {
	return &Callback{
		callbackResQueue: make(chan *CallbackRes, 1),
	}
}

type Handle func(interface{})

type CallbackHandle func(interface{}) (interface{}, error)

type CallbackRes struct {
	Data interface{}
	Err  error
}

type Event struct {
	F        Handle
	Arg      interface{}
}

type CallbackEvent struct {
	F        CallbackHandle
	Arg      interface{}
	callback *Callback
}

func GetRoutine() *Routine {
	once.Do(func() {
		globalRoutine = &Routine{
			queue:         make(chan *Event, 65535),
			callbackQueue: make(chan *CallbackEvent, 65535)}
	})
	return globalRoutine
}

func PutEvent(event *Event) {
	GetRoutine().queue <- event
}

func PutCallbackEvent(event *CallbackEvent) *Callback {
	if nil == event {
		return nil
	}
	routine := GetRoutine()
	callback := GetCallback()
	event.callback = callback
	routine.callbackQueue <- event
	return callback
}

func (ptr *Callback) Get() (interface{}, error) {
	res := <-ptr.callbackResQueue
	if nil == res {
		return nil, fmt.Errorf("callback error")
	}
	close(ptr.callbackResQueue)
	return res.Data, res.Err
}

func (ptr *Routine) Do() {

	for {
		event := <-ptr.queue
		if nil == event {
			continue
		}
		event.F(event.Arg)
	}
}

func (ptr *Routine) CallbackDo() {

	for {
		event := <-ptr.callbackQueue
		if nil == event {
			continue
		}

		res, err := event.F(event.Arg)
		event.callback.callbackResQueue <- &CallbackRes{
			Data: res,
			Err:  err,
		}

	}
}

func (ptr *Routine) Init() {

	go ptr.Do()

	go ptr.CallbackDo()
}
