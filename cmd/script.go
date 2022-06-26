/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// scriptCmd represents the script command
var scriptCmd = &cobra.Command{
	Use:   "script",
	Short: "LucaSystem script文件",
	Long: `LucaSystem script文件
无具体文件头，确定是LucaSystem引擎的游戏，SCRIPT.PAK中的文件即为script文件
其中'_VARNUM'、'_CGMODE'、'_SCR_LABEL'、'_VOICE_PARAM'、'_BUILD_COUNT'、'_TASK'等文件不支持单个解析`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("script called")
	},
}

func init() {
	rootCmd.AddCommand(scriptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scriptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scriptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
