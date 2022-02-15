/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"

	"github.com/connorjcantrell/dapper/config"
	"github.com/connorjcantrell/dapper/dapperfs"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

// Flags
var boilerplate string
var globalBytes uint
var globalInts uint
var localBytes uint
var localInts uint

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize dapper in current working directory",
	Long: `Initialize dapper in current working dirctory
	Creates project a /.dapper and /src directory. 

	A config.json file is generated inside /.dapper, which contains application
	details that will be used to interact with the Algorand Blockchain

	If boilerplate flag was used, code templates will be written inside /src`,
	Run: func(cmd *cobra.Command, args []string) {
		fs, err := dapperfs.New()
		if err != nil {
			log.Fatal(err)
		}

		// Create project directories
		fs.Mkdir(".dapper")
		fs.Mkdir("src")

		// Write boilerplate code if boilerplate flag was provided
		if cmd.Flags().Changed("boilerplate") {
			fs.CopyFromBoilerplateDir(boilerplate)
		}

		config := config.ApplicationDetails{
			GlobalStateSchema: config.GlobalStateSchema{
				NumByteSlice: globalBytes,
				NumUint:      globalInts,
			},
			LocalStateSchema: config.LocalStateSchema{
				NumByteSlice: localBytes,
				NumUint:      localInts,
			},
		}
		err = fs.WriteStructToJSON(config, ".dapper", "config.json")
		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	flag.StringVarP(&boilerplate, "boilerplate", "b", "", "Start project with specified boilerplate code")
	flag.UintVar(&globalBytes, "global-byteslices", 0, "Maximum number of byte slices that may be stored in the global key/value store. Immutable")
	flag.UintVar(&globalInts, "global-ints", 0, "Maximum number of integer values that may be stored in the global key/value store. Immutable.")
	flag.UintVar(&localBytes, "local-byteslices", 0, "Maximum number of byte slices that may be stored in local (per-account) key/value stores for this app. Immutable.")
	flag.UintVar(&localInts, "local-ints", 0, "Maximum number of integer values that may be stored in local (per-account) key/value stores for this app. Immutable.")
}
