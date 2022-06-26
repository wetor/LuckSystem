/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"lucksystem/charset"
	"lucksystem/pak"
	"os"
)

// pakExtractCmd represents the pakExtract command
var pakExtractCmd = &cobra.Command{
	Use:   "extract",
	Short: "解包Pak文件",
	Long: `LucaSystem pak文件
无具体文件头，确定是LucaSystem引擎的游戏，文件名为大写的***.PAK`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pakExtract called")
		p := pak.LoadPak(PakInput, charset.Charset(Charset))
		out, err := os.Create(PakOutput)
		if err != nil {
			glog.Fatalln(err)
		}
		defer out.Close()
		if len(PakAll) > 0 {
			err := p.Export(out, "all", PakAll)
			if err != nil {
				glog.Fatalln(err)
			}
		} else if PakIndex >= 0 {
			err := p.Export(out, "index", PakIndex)
			if err != nil {
				glog.Fatalln(err)
			}
		} else if PakId >= 0 {
			err := p.Export(out, "id", PakId)
			if err != nil {
				glog.Fatalln(err)
			}
		} else if len(PakName) > 0 {
			err := p.Export(out, "name", PakName)
			if err != nil {
				glog.Fatalln(err)
			}
		} else {
			fmt.Println("Error: required flag(s) \"all\" or \"index\" or \"id\" or \"name\" not set")
		}

	},
}

var (
	PakAll   string
	PakIndex int
	PakId    int
	PakName  string
)

func init() {
	pakCmd.AddCommand(pakExtractCmd)

	pakExtractCmd.Flags().StringVarP(&PakAll, "all", "a", "", "提取所有文件保存到文件夹")
	pakExtractCmd.Flags().IntVarP(&PakIndex, "index", "x", -1, "提取指定位置文件，从0开始")
	pakExtractCmd.Flags().IntVarP(&PakId, "id", "d", -1, "提取指定ID文件")
	pakExtractCmd.Flags().StringVarP(&PakName, "name", "n", "", "提取指定文件名文件")

	pakExtractCmd.MarkFlagsMutuallyExclusive("all", "index", "id", "name")
}
