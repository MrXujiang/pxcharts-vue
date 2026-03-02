package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	configFile string
	env        string
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "mvtable",
	Short: "Doc Space - 一个基于Go的文档协同编辑平台",
	Long: `
			Doc Space 是一个基于Go语言的文档协同编辑平台，
			集成了Gin、GORM、Zap、JWT、Cobra、Wire等常用框架和工具。
			
			主要功能：
			- 用户管理和认证
			- 文档管理和协同编辑
			- 权限控制
			- 日志记录
			- 数据库操作
			- API文档
		`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 全局前置处理
	},
}

// Execute 执行命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "执行命令失败: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// 全局标志
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "配置文件路径")
	rootCmd.PersistentFlags().StringVarP(&env, "env", "e", "dev", "环境 (dev|prod)")

	// 添加子命令
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(versionCmd)
}
