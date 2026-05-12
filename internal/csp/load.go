package csp

import (
	"context"
	"encoding/json"
	"throughput/internal/sqlite"
)

func Load(s *sqlite.SQLite, rowsPerBundle int, in <-chan *Package) {
	rows := make([]*Package, rowsPerBundle)
	ndx := 0
	bundles := 0
	for {
		bundle := <-in
		rows[ndx] = bundle
		ndx += 1
		if ndx == rowsPerBundle {
			bundles += 1
			// log.Println("bundles: ", bundles)
			s.BeginTx()
			for _, b := range rows {
				jsonData, _ := json.Marshal(b.Input)
				s.QTx.InsertBundle(context.Background(), string(jsonData))
			}
			s.EndTx()
			rows = make([]*Package, rowsPerBundle)
			ndx = 0
		}
	}
}
