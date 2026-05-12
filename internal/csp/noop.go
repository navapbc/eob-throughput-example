package csp

import (
	"throughput/internal/sqlite"
	"time"
)

func NoOp(s *sqlite.SQLite, in chan *Package, out chan *Package) {
	for {
		val := <-in
		val.EndTime = time.Now()
		out <- val
	}
}
