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
	start                = 3
)

var avgDatas struct {
	sync.Mutex
	data map[string]*avgData
}

func (ag *agregator) addData(data *dataSource) {
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
	if ag.Start != nil {
		ag.Start <- start // allow the aggregation
		ag.Start = nil
	}
}

func (ag *agregator) newAg(wg *sync.WaitGroup) {

	ag.newSub()

	if avgDatas.data == nil {
		avgDatas.data = map[string]*avgData{}
	}

	wg.Add(1)
	go func(ag *agregator) {
		defer wg.Done()

		// Starts and aggregates data while input channel is "alive".
		<-ag.Start

		ticker := time.NewTicker(time.Second * time.Duration(ag.AgregatePeriodS))
		cWrite := make(chan avgData)
		cStop := make(chan bool)
		writeData(wg, cWrite, cStop)

		for {
			select {
			case <-ag.Event:
				cStop <- true
				return
			case <-ticker.C:
				if len(avgDatas.data) == 0 {
					continue
				}
				avgDatas.Lock()
				for _, v := range avgDatas.data {
					for _, subID := range ag.SubIds {
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
	}(ag)
}
