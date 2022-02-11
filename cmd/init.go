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
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fs, err := dapperfs.NewFs()
		if err != nil {
			log.Fatal(err)
		}

		// Create project directories
		fs.Mkdir(".dapper")
		fs.Mkdir("src")
		fs.Mkdir("public")

		// Write boilerplate code if boilerplate flag was provided
		if cmd.Flags().Changed("boilerplate") {
			fs.CopyFromTemplates(boilerplate)
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
		err = fs.WriteConfig(config)
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

/*
TODO:
`dappman init` will perform the following:
    - Retrieve Current Working Directory (i.e. ~/algorand/projects/my-app)
    - Create directories inside Current Working Directory
		- ./.dappman
        - ./src
    	- ./public
    - Create files in ./src
        - approval_program.py
    	- clear_state_program.py
	- Create a file `config.json` inside `./.dappman`
	    - `config.json` contains application specific information like Application ID, global/local state schemas, etc.
*/
