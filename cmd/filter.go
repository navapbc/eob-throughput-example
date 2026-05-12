/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"throughput/internal/csp"
	sqlite "throughput/internal/sqlite"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func PrintMemUsage(m runtime.MemStats) {
	// var m runtime.MemStats
	// runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

var peak uint64

func peeker() {
	var m runtime.MemStats
	for {
		runtime.ReadMemStats(&m)
		if m.Sys > peak {
			peak = m.Sys
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// filterCmd represents the filter command
var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		filters := []csp.GoJqFilter{
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

		rows := viper.GetInt("rows")
		bundle := viper.GetInt("bundle")
		writing := viper.GetBool("writing")
		noisy := viper.GetBool("noisy")
		CONCURRENCY := viper.GetInt("concurrency")

		var startMem, endMem runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&startMem)

		// Integer division
		if rows < ROWS_PER_BUNDLE {
			rows = ROWS_PER_BUNDLE
		}
		rows = (rows / ROWS_PER_BUNDLE) * ROWS_PER_BUNDLE

		if noisy {
			log.Println("rows:", rows)
			log.Println("write bundle:", bundle)
			log.Println("writing data:", writing)
			log.Println("loading rows: ", rows)
			log.Println("running with concurrent filters:", CONCURRENCY)
		}

		s := sqlite.NewSQLite(rows)

		req := make(chan int)
		resp := make(chan []*csp.Package)

		demux := make(chan *csp.Package)
		gcchan := make(chan *csp.Package)
		showchan := make(chan *csp.Package)

		var wg sync.WaitGroup
		// emit -> mux -> filter -> demux -> store

		// A process network to do things stepwise
		// Emit the files as byte arrays
		// go csp.Emit(viper.GetString("directory"), req, resp)
		wg.Add(1)
		go csp.SimuDB(viper.GetString("directory"), req, resp)
		for range CONCURRENCY {
			go csp.Filter(filters, req, resp, demux)
		}

		go csp.Store(s, bundle, demux, gcchan)

		go csp.PeriodicGC(1000, gcchan, showchan)
		go csp.Show(&wg, showchan)

		start := time.Now()
		go peeker()

		wg.Wait()
		end := time.Now()

		runtime.GC()
		runtime.ReadMemStats(&endMem)

		dur := end.Sub(start)
		avg := dur / time.Duration(rows)
		// memDelta := endMem.Alloc - startMem.Alloc
		if noisy {
			fmt.Printf("concurrency,rows,bundle,writing,total time (ms),avg time (ms),alloc,totalalloc,sys,peak,numgc\n")
		}
		fmt.Printf("%d,%d,%d,%v,%d,%d,%d,%d,%d,%d,%d\n",
			CONCURRENCY,
			rows,
			bundle,
			writing,
			dur.Milliseconds(),
			avg.Microseconds(),
			bToMb(endMem.Alloc),
			bToMb(endMem.TotalAlloc),
			bToMb(endMem.Sys),
			bToMb(peak),
			endMem.NumGC,
		)
		// PrintMemUsage(startMem)
		// PrintMemUsage(endMem)
		time.Sleep(3 * time.Second)
	},
}

func init() {
	rootCmd.AddCommand(filterCmd)
	filterCmd.PersistentFlags().IntVarP(&rows, "rows", "r", 1000, "rows to simulate")
	err := viper.BindPFlag("rows", filterCmd.PersistentFlags().Lookup("rows"))
	if err != nil {
		panic(err)
	}

	filterCmd.PersistentFlags().IntVarP(&rowsPerBundle, "bundle", "b", 100, "rows to buffer for writing")
	err = viper.BindPFlag("bundle", filterCmd.PersistentFlags().Lookup("bundle"))

	filterCmd.PersistentFlags().IntVarP(&isWriting, "writing", "w", 0, "actually write data")
	err = viper.BindPFlag("writing", filterCmd.PersistentFlags().Lookup("writing"))
	if err != nil {
		panic(err)
	}

	filterCmd.PersistentFlags().IntVarP(&noisy, "noisy", "n", 0, "be noisy")
	err = viper.BindPFlag("noisy", filterCmd.PersistentFlags().Lookup("noisy"))
	if err != nil {
		panic(err)
	}

	filterCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "c", 1, "rows to simulate")
	_ = viper.BindPFlag("concurrency", filterCmd.PersistentFlags().Lookup("concurrency"))

	filterCmd.PersistentFlags().StringVarP(&directory, "directory", "d", "bundles", "JSON bundle directory")
	_ = viper.BindPFlag("directory", filterCmd.PersistentFlags().Lookup("directory"))
}
