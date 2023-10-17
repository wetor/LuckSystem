/*
Copyright © 2022 WeTor wetorx@qq.com
*/
package cmd

import (
	"flag"
	"os"
	"strconv"

	"github.com/go-restruct/restruct"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "LuckSystem",
	Version: "2.0.2",
	Short:   "LucaSystem引擎工具集",
	Long: `LucaSystem引擎工具集
https://github.com/wetor/LuckSystem
wetor(wetorx@qq.com)`,
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
	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		_ = rootCmd.Help()
		os.Exit(1)
	}

	restruct.EnableExprBeta()

	rootCmd.Flags().BoolVar(&Log, "log", true, "启用日志")
	rootCmd.Flags().IntVar(&LogLevel, "log_level", 5, "输出日志等级")
	rootCmd.Flags().StringVar(&LogDir, "log_dir", "log", "保存日志路径")
}
