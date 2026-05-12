package csp

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func Emit(dir string, req chan bool, resp chan *Package) {
	bufferedFiles := make(map[string][]byte)
	entries, _ := os.ReadDir(dir)
	if len(entries) < 1 {
		log.Println("no files found")
		log.Println(dir)
		os.Exit(-1)
	}
	index := 0
	for {
		for _, e := range entries {
			// Wait for a request
			<-req

			P := &Package{}
			P.StartTime = time.Now()
			fullPath := filepath.Join(dir, e.Name())
			// Buffering these shaves 2ms off each.
			// Still cannot see any benefits of concurrency
			if _, ok := bufferedFiles[string(fullPath)]; !ok {
				jsonFile, _ := os.Open(fullPath)
				byteValue, _ := io.ReadAll(jsonFile)
				bufferedFiles[string(fullPath)] = byteValue
			}
			var input interface{}
			if err := json.Unmarshal(bufferedFiles[fullPath], &input); err != nil {
				log.Fatalln(err)
			}

			P.Input = input
			P.Index = index
			index += 1
			P.Filename = e.Name()
			P.Filters = []GoJqFilter{
				{
					Name: "Remove amount from adjudication",
					GoJq: "del((.entry // [])[] | .resource | (.item // [])[] | (.adjudication // [])[] | .amount)",
				},
				{
					Name: "Remove amount from total",
					GoJq: "del((.entry // [])[] | .resource | (.total // [])[] | .amount)",
				},
				{
					Name: "Remove amount from payment",
					GoJq: "del((.entry // [])[] | .resource | .payment | .amount)",
				},
				{
					Name: "Remove usedMoney from benefitBalance",
					GoJq: "del((.entry // [])[] | .resource | (.benefitBalance // [])[]? | .usedMoney)",
				},
			}

			resp <- P
		}
	}
}
