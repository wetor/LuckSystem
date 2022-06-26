/*
Copyright © 2022 WeTor wetorx@qq.com

*/
package cmd

import (
	"flag"
	"github.com/go-restruct/restruct"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "LuckSystem",
	Version: "1.0.0",
	Short:   "LucaSystem引擎工具集",
	Long:    `LucaSystem引擎工具集`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if Log {
		flag.Set("alsologtostderr", "true")
		flag.Set("log_dir", LogDir)
		flag.Set("v", strconv.Itoa(LogLevel))
		flag.Parse()
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	Log      bool
	LogLevel int
	LogDir   string
)

func init() {
	restruct.EnableExprBeta()
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lucksystem.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolVar(&Log, "log", true, "启用日志")
	rootCmd.Flags().IntVar(&LogLevel, "log_level", 10, "输出日志等级")
	rootCmd.Flags().StringVar(&LogDir, "log_dir", "log", "保存日志路径")
}
