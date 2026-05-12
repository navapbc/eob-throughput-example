/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"strings"
	"sync"
	"throughput/internal/csp"
	sqlite "throughput/internal/sqlite"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const ROWS_PER_BUNDLE = 100

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		rows := viper.GetInt("rows")
		// Integer division
		if rows < ROWS_PER_BUNDLE {
			rows = ROWS_PER_BUNDLE
		}
		rows = (rows / ROWS_PER_BUNDLE) * ROWS_PER_BUNDLE
		log.Println("loading rows: ", rows)
		s := sqlite.NewSQLite(rows)

		req := make(chan bool)
		resp := make(chan *csp.Package)
		var wg sync.WaitGroup

		log.Println("running network")
		go csp.Emit(viper.GetString("directory"), req, resp)
		go csp.Load(s, ROWS_PER_BUNDLE, resp)

		wg.Go(func() {
			for range rows {
				req <- true
			}
		})

		wg.Wait()
		log.Println("done waiting")
		time.Sleep(5 * time.Second)
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)

	loadCmd.Flags().IntVarP(&rows, "rows", "r", ROWS_PER_BUNDLE, "rows to simulate")
	_ = viper.BindPFlag("rows", loadCmd.Flags().Lookup("rows"))

	loadCmd.Flags().StringVarP(&directory, "directory", "d", "bundles", "JSON bundle directory")
	_ = viper.BindPFlag("directory", loadCmd.Flags().Lookup("directory"))

	// Env: FLAGSAPP_PORT=9090
	viper.SetEnvPrefix("throughput")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}
