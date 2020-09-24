package main

import (
	"math"
	"sync"
	"time"
)

type event uint8

const (
	allGeneratorsStopped = 1
	exit                 = 2
)

var avgDatas struct {
	sync.Mutex
	data map[string]*avgData
}

func (agregator agregator) addData(data dataSource) {
	avgDatas.Lock()
	if _, inMap := avgDatas.data[data.ID]; inMap {
		avgDatas.data[data.ID].Lenght++
		avgDatas.data[data.ID].Value += data.Value
	} else {
		var a = avgData{
			ID:     data.ID,
			Value:  data.Value,
			Lenght: 1,
		}
		avgDatas.data[data.ID] = &a
	}
	avgDatas.Unlock()
}

func (agregator agregator) new(wg *sync.WaitGroup) {

	if avgDatas.data == nil {
		avgDatas.data = map[string]*avgData{}
	}

	ticker := time.NewTicker(time.Second * time.Duration(agregator.AgregatePeriodS))
	newSub(&agregator)
	cWrite := make(chan avgData)
	cStop := make(chan bool)
	writeData(wg, cWrite, cStop)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-agregator.Event:
				cStop <- true
				return
			case <-ticker.C:
				if len(avgDatas.data) == 0 {
					continue
				}
				avgDatas.Lock()
				for _, v := range avgDatas.data {
					for _, subID := range agregator.SubIds {
						if subID == v.ID {
							v.Value /= float64(v.Lenght)            // finding the average
							v.Value = math.Round(v.Value*100) / 100 // rounding to .**
							cWrite <- *v
							delete(avgDatas.data, subID)
						}
					}
				}
				avgDatas.Unlock()
			}
		}
	}()
}
