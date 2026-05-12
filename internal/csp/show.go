package csp

import (
	"sync"
	"time"

	"github.com/spf13/viper"
)

func Show(wg *sync.WaitGroup, in chan *Package) {
	// Average the last 100 values
	times := make([]time.Duration, 100)
	time_ndx := 0
	last_avg := time.Duration(0)
	bundlesProcessed := 0
	for {
		bundle := <-in
		// b, _ := json.MarshalIndent(bundle.Output, "", "  ")
		// fmt.Println(string(b))
		delta := bundle.EndTime.Sub(bundle.StartTime)
		times[time_ndx] = delta
		time_ndx += 1
		if time_ndx == len(times) {
			time_ndx = 0
		}
		// log.Printf("%s: %s\n", bundle.Filename, delta.String())

		var total time.Duration
		for _, d := range times {
			total += d
		}

		avg := total / time.Duration(len(times))
		if avg-last_avg > 0 {
			//log.Printf("avg: %s\n", avg.String())
			last_avg = avg
		}
		bundlesProcessed += 1
		if bundlesProcessed >= viper.GetInt("rows") {
			wg.Done()
			for {
				time.Sleep(1 * time.Second)
			}
		}
	}
}
