package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sydneyowl/shx8800-ble-connector/config"
	"github.com/sydneyowl/shx8800-ble-connector/pkg/logger"
)

var (
	Verbose      = false
	Vverbose     = false
	PrintVersion = false
)

func printVer() {
	fmt.Println("SHX8800 dat editor")
	fmt.Printf("Version: %s\n", config.VER)
	fmt.Printf("Commit: %s\n", config.COMMIT)
	fmt.Printf("Build Time: %s\n", config.BUILDTIME)
}

var BaseCmd = &cobra.Command{
	Use:   "SHX8800_BLE",
	Short: "SHX8800_BLE",
	Long:  `SHX8800_BLE - Easily transfer data to shx8800 on pc`,
	Run: func(cmd *cobra.Command, args []string) {
		if PrintVersion {
			printVer()
			return
		}
		logger.InitLog(Verbose, Vverbose)
		StartAndRun()
	},
}

func init() {
	BaseCmd.PersistentFlags().BoolVar(&PrintVersion, "version", false, "Print Version")
	BaseCmd.PersistentFlags().BoolVar(&Verbose, "verbose", false, "Print Debug Level logs")
	BaseCmd.PersistentFlags().BoolVar(&Vverbose, "vverbose", false, "Print Debug/Trace Level logs")
}
