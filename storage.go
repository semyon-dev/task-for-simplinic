package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

var isOpen bool
var file1 *os.File
var isFileStorage bool

func writeData(wg *sync.WaitGroup, c chan avgData, cStop chan bool) {

	if !isOpen && isFileStorage {
		var err error
		if file1, err = os.Create("storage.txt"); err != nil {
			log.Println(err)
		} else {
			isOpen = true
		}
	}

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		for {
			select {
			case <-cStop:
				return
			case data := <-c:
				res, err := json.Marshal(data)
				if err != nil {
					log.Println(err)
				}
				res = append(res, '\n')

				if isFileStorage {
					_, err = file1.Write(res)
					if err != nil {
						log.Println(err)
					}
				} else {
					fmt.Println(string(res))
				}
			}
		}
	}()
}
