/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "LucaSystem cz图像",
	Long: `LucaSystem cz图像
文件头为'CZ0'、'CZ1'、'CZ2'、'CZ3'、'CZ4'等
目前实现'CZ0'、'CZ1'、'CZ3'和CZ2的导出`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("image called")
	},
}
var (
	CzInput  string // 输入
	CzOutput string // 输出
	CzSource string // 原文件
)

func init() {
	rootCmd.AddCommand(imageCmd)

	imageCmd.PersistentFlags().StringVarP(&CzSource, "source", "s", "", "原cz文件名")
	imageCmd.PersistentFlags().StringVarP(&CzInput, "input", "i", "", "输入文件")
	imageCmd.PersistentFlags().StringVarP(&CzOutput, "output", "o", "", "输出文件")

	imageCmd.MarkPersistentFlagRequired("input")
	imageCmd.MarkPersistentFlagRequired("output")
	imageCmd.MarkFlagsRequiredTogether("output", "input")
}
