package csp

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func CopyMap(m map[string]interface{}) map[string]interface{} {
	cp := make(map[string]interface{})
	for k, v := range m {
		// Recursively copy nested maps
		if vm, ok := v.(map[string]interface{}); ok {
			cp[k] = CopyMap(vm)
		} else if vs, ok := v.([]interface{}); ok {
			// Slices also need deep copying
			newSlice := make([]interface{}, len(vs))
			for i, sv := range vs {
				if svm, ok := sv.(map[string]interface{}); ok {
					newSlice[i] = CopyMap(svm)
				} else {
					newSlice[i] = sv
				}
			}
			cp[k] = newSlice
		} else {
			cp[k] = v
		}
	}
	return cp
}

func SimuDB(dir string, req chan int, resp chan []*Package) {
	entries, _ := os.ReadDir(dir)
	if len(entries) < 1 {
		log.Println("no files found")
		log.Println(dir)
		os.Exit(-1)
	}
	// Pre-buffer the files into a dictionary.
	bufferedFiles := make([][]byte, 0)
	fileNames := make([]string, 0)
	jsonFiles := make([]any, 0)
	for _, e := range entries {
		fullPath := filepath.Join(dir, e.Name())
		// Buffering these shaves 2ms off each.
		jsonFile, err := os.Open(fullPath)
		if err != nil {
			panic(err)
		}
		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			panic(err)
		}
		var input interface{}
		if err := json.Unmarshal(byteValue, &input); err != nil {
			panic(err)
		}
		// log.Printf("loading %s index %d\n", e.Name(), ndx)
		bufferedFiles = append(bufferedFiles, byteValue)
		jsonFiles = append(jsonFiles, input)
		fileNames = append(fileNames, e.Name())
	}

	index := 0
	for {
		// Wait for a request
		howMany := <-req
		packages := make([]*Package, howMany)
		for ndx := range howMany {
			P := &Package{}

			// Consider
			// https://pkg.go.dev/github.com/tiendc/go-deepcopy
			P.Input = CopyMap(jsonFiles[ndx%len(entries)].(map[string]any))
			P.Index = index
			index += 1
			P.Filename = fileNames[ndx%len(entries)]
			P.StartTime = time.Now()

			packages[ndx] = P

		}
		resp <- packages
	}
}
