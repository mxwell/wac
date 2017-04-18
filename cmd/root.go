package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/mxwell/wac/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "wac",
	Short: "Contestant helper",
	Long: `WAC is a CLI tool that helps contestants of programming contests
to write, build and test code of solutions.

Version 0.1`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	log.SetFlags(0)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	util.CheckConfiguration()
	viper.SetConfigName("wac") // name of config file (without extension)
	viper.AddConfigPath(util.GetDefaultLocation())

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
