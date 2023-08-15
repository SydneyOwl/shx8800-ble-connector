package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sydneyowl/shx8800-ble-connector/config"
	"github.com/sydneyowl/shx8800-ble-connector/pkg/logger"
)

var (
	Verbose  = false
	Vverbose = false
)

var BaseCmd = &cobra.Command{
	Use:     "SHX8800_BLE",
	Short:   "SHX8800_BLE",
	Version: config.VER,
	Long:    `SHX8800_BLE - Easily transfer data to shx8800 on pc`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.InitLog(Verbose, Vverbose)
		StartAndRun()
	},
}

func init() {
	BaseCmd.PersistentFlags().BoolVar(&Verbose, "verbose", false, "Print Debug Level logs")
	BaseCmd.PersistentFlags().BoolVar(&Vverbose, "vverbose", false, "Print Debug/Trace Level logs")
}
