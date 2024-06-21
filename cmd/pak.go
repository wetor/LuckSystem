/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"lucksystem/charset"
	"lucksystem/pak"

	"github.com/spf13/cobra"
)

// pakCmd represents the pak command
var pakCmd = &cobra.Command{
	Use:   "pak",
	Short: "LucaSystem pak文件",
	Long: `LucaSystem pak文件
无具体文件头，确定是LucaSystem引擎的游戏，文件名为大写的***.PAK`,
	Run: func(cmd *cobra.Command, args []string) {
		if List {
			if len(PakSource) == 0 {
				fmt.Println("Error: required flag(s) \"source\" not set")
				return
			}
			fmt.Println("index,id,offset,size,name")
			p := pak.LoadPak(PakSource, charset.Charset(Charset))
			for i, f := range p.Files {
				fmt.Printf("%d,%d,%d,%d,%s\n", i, p.NameMap[f.Name], f.Offset, f.Length, f.Name)
			}
		}
	},
}
var (
	List      bool   // 列表模式
	Charset   string // 编码
	PakInput  string // 输入
	PakOutput string // 输出
	PakSource string // 原文件
)

func init() {
	rootCmd.AddCommand(pakCmd)

	pakCmd.Flags().BoolVarP(&List, "list", "L", false, "查看文件列表，需要source")
	pakCmd.PersistentFlags().StringVarP(&Charset, "charset", "c", string(charset.UTF_8), "字符串编码")

	pakCmd.PersistentFlags().StringVarP(&PakSource, "source", "s", "", "原Pak文件名")
	pakCmd.PersistentFlags().StringVarP(&PakInput, "input", "i", "", "输入文件或文件夹")
	pakCmd.PersistentFlags().StringVarP(&PakOutput, "output", "o", "", "输出文件或文件夹")

	pakCmd.MarkFlagsRequiredTogether("output", "input")
}
