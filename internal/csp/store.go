package csp

import (
	"context"
	"encoding/json"
	"throughput/internal/sqlite"

	"github.com/spf13/viper"
)

func Store(s *sqlite.SQLite, rowsPerTx int, in <-chan *Package, out chan<- *Package) {
	rows := make([]*Package, rowsPerTx)
	ndx := 0
	bundles := 0
	for {
		bundle := <-in
		rows[ndx] = bundle
		ndx += 1
		if ndx == rowsPerTx {
			bundles += 1
			// log.Println("bundles: ", bundles)
			if viper.GetBool("writing") {
				s.BeginTx()
				for _, b := range rows {
					jsonData, _ := json.Marshal(b.Input)
					s.QTx.InsertBundle(context.Background(), string(jsonData))
				}
				s.EndTx()
			}
			rows = make([]*Package, rowsPerTx)
			ndx = 0
		}
		out <- bundle
	}
}
