package csp

import (
	"runtime"
)

func PeriodicGC[T any](period int, in chan T, out chan T) {
	counter := 0
	for {
		v := <-in
		counter += 1
		if counter == period {
			// It is possible to set a GOMEM limit that would force GCs
			// at some pre-defined level below the memory limits of (say)
			// a container. This saves OOM kills.
			runtime.GC()
			counter = 0
		}
		out <- v
	}
}
