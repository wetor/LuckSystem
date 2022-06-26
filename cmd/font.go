/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// fontCmd represents the font command
var fontCmd = &cobra.Command{
	Use:   "font",
	Short: "LucaSystem 字体",
	Long: `LucaSystem 字体
文件头为'CZ0'、'CZ1'、'CZ2'、'CZ3'、'CZ4'等，同时需要配合与文件名字号一致的info文件
如'明朝24'和'info24'为一个字体`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("font called")
	},
}
var (
	FontCzSource   string // 输入
	FontInfoSource string // 输入
	FontOutput     string // 输出
	FontInfoOutput string // 输出info或info字符集
)

func init() {
	rootCmd.AddCommand(fontCmd)

	fontCmd.PersistentFlags().StringVarP(&FontCzSource, "source", "s", "", "原字体cz文件，用作输入")
	fontCmd.PersistentFlags().StringVarP(&FontInfoSource, "source_info", "S", "", "原字体对应字号info文件，用作输入")
	fontCmd.PersistentFlags().StringVarP(&FontOutput, "output", "o", "", "输出文件")

	fontCmd.MarkPersistentFlagRequired("source")
	fontCmd.MarkPersistentFlagRequired("source_info")
	fontCmd.MarkPersistentFlagRequired("output")
	fontCmd.MarkFlagsRequiredTogether("source", "source_info", "output")
}
