package csp

import (
	"log"
	"time"

	"github.com/itchyny/gojq"
	"github.com/spf13/viper"
)

func Filter(filters []GoJqFilter, req chan int, unfiltered chan []*Package, filtered chan<- *Package) {
	for {
		req <- viper.GetInt("bundle")
		bundles := <-unfiltered

		for _, bundle := range bundles {
			bundle.Output = bundle.Input

			for _, gojqf := range filters {
				// Noisy; shows us every filter action name
				// log.Println("running: " + gojqf.Name)
				query, err := gojq.Parse(gojqf.GoJq)
				if err != nil {
					log.Fatalln(err)
				}
				iter := query.Run(bundle.Output)
				for {
					v, ok := iter.Next()
					if !ok {
						break
					}
					if err, ok := v.(error); ok {
						log.Fatalln(err)
					}
					//fmt.Printf("%s: %v\n\n", bundle.Filename, v) // Output: true
					bundle.Output = v
				}
			}

			bundle.EndTime = time.Now()
			// This ships them one-by-one...
			filtered <- bundle
		}

	}
}
