package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"
)

var dataSources struct {
	sync.Mutex
	data map[string]*dataSource
}

var countGenerators int

func main() {

	configBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("Error reading the config: ", err)
	}
	var config config
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		log.Fatal(err)
	}

	if config.StorageType == 1 {
		isFileStorage = true
	}

	messageQueueMaxSize = config.Queue.Size

	dataSources.data = map[string]*dataSource{}

	// filling generators
	for y := 0; y < len(config.Generators); y++ {
		for i := 0; i < len(config.Generators[y].DataSources); i++ {
			dataSources.data[config.Generators[y].DataSources[i].ID] = &config.Generators[y].DataSources[i]
		}
	}

	wg := new(sync.WaitGroup)

	for i := 0; i < len(config.Agregators); i++ {
		cEvent := make(chan event)
		cStart := make(chan event, 1)
		config.Agregators[i].Event = cEvent
		config.Agregators[i].Start = cStart
		config.Agregators[i].newAg(wg)
	}

	var geratorsS []chan bool

	wg.Add(len(config.Generators))
	for i := 0; i < len(config.Generators); i++ {
		chanStop := make(chan bool, 1)
		geratorsS = append(geratorsS, chanStop)
		go config.Generators[i].new(wg, chanStop)
	}

	// react to SIGINT UNIX signal
	chanForceExit := make(chan os.Signal, 1)
	signal.Notify(chanForceExit, os.Interrupt, os.Kill)
	go func() {
		<-chanForceExit
		sendEvent(exit)
		for i := 0; i < len(geratorsS); i++ {
			geratorsS[i] <- true
			close(geratorsS[i])
		}
	}()
	wg.Wait()
}

func (generator generator) new(wg *sync.WaitGroup, stopForce chan bool) {

	countGenerators++

	stop := time.NewTimer(time.Second * time.Duration(generator.TimeoutS))
	tick := time.NewTicker(time.Second * time.Duration(generator.SendPeriodS))

	defer func() {
		wg.Done()
		tick.Stop()
		closeGenerator()
	}()

	for {
		select {
		case <-stop.C:
			// stop - Timer, which will give a signal to stop working after TimeoutS of seconds
			return
		case <-tick.C:
			// tick that sends a signal to complete the job every sendPeriodS seconds
			for dataID := range dataSources.data {
				increaseValue(dataID, wg)
			}
		case <-stopForce:
			return
		}
	}
}

func increaseValue(dataSourceID string, wg *sync.WaitGroup) {
	dataSources.Lock()
	dataSources.data[dataSourceID].Value += float64(rand.Intn(dataSources.data[dataSourceID].MaxChangeStep))
	dataSources.Unlock()
	wg.Add(1)
	go sendData(dataSources.data[dataSourceID], wg)
}

func closeGenerator() {
	var m sync.Mutex
	m.Lock()
	countGenerators--
	m.Unlock()
	if countGenerators == 0 {
		sendEvent(allGeneratorsStopped)
	}
}
