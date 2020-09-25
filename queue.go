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

var subs []*agregator

func (ag *agregator) newSub() {
	subs = append(subs, ag)
}

func sendEvent(event event) {
	for a := 0; a < len(subs); a++ {
		subs[a].Event <- event
	}
}

func sendData(source *dataSource, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < len(subs); i++ {
		for y := 0; y < len(subs[i].SubIds); y++ {
			if subs[i].SubIds[y] == source.ID {
				err := add(*source)
				if err != nil {
					log.Println(err)
				}
				subs[i].addData(source)
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
