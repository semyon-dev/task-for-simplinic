package main

import (
	"container/list"
	"errors"
	"log"
	"sync"
)

var messageQueue = list.New()
var messageQueueMaxSize int
var m sync.Mutex

var subs []agregator

func newSub(agregator *agregator) {
	cEvent := make(chan event)
	agregator.Event = cEvent
	subs = append(subs, *agregator)
}

func sendEvent(event event) {
	for _, v := range subs {
		v.Event <- event
	}
}

func sendData(source dataSource) {
	for _, v := range subs {
		for _, id := range v.SubIds {
			if id == source.ID {
				err := add(source)
				if err != nil {
					log.Println(err)
				}
				v.addData(source)
				remove()
				break
			}
		}
	}
}

// removing the first item from the queue
func remove() {
	m.Lock()
	messageQueue.Remove(messageQueue.Front())
	m.Unlock()
}

// adding an item to the end of the queue
func add(elem dataSource) error {
	m.Lock()
	if messageQueueMaxSize == messageQueue.Len() {
		m.Unlock()
		return errors.New("the queue is full")
	}
	messageQueue.PushBack(elem)
	m.Unlock()
	return nil
}
