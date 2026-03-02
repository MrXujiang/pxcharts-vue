package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = "1.0.0"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

// versionCmd 版本命令
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  `显示应用程序的版本信息`,
	Run:   runVersion,
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("Go Server Template\n")
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Build Time: %s\n", BuildTime)
	fmt.Printf("Git Commit: %s\n", GitCommit)
}
