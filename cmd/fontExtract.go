/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/golang/glog"
	"lucksystem/font"
	"os"

	"github.com/spf13/cobra"
)

// fontExtractCmd represents the fontExtract command
var fontExtractCmd = &cobra.Command{
	Use:   "extract",
	Short: "提取字体图像和info",
	Long: `LucaSystem font文件
为FONT.PAK内容，字体为cz图像，并对应一个字号的info文件，如"明朝32"和"info32"为明朝32号字体
输出为png图像，info输出为txt字符集，对应png图像中的全字符`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fontExtract called")
		f := font.LoadLucaFontFile(FontInfoSource, FontCzSource)
		out, err := os.Create(FontOutput)
		if err != nil {
			glog.Fatalln(err)
		}
		defer out.Close()
		err = f.Export(out, FontInfoOutput)
		if err != nil {
			glog.Fatalln(err)
		}
	},
}

func init() {
	fontCmd.AddCommand(fontExtractCmd)
	fontExtractCmd.Flags().StringVarP(&FontInfoOutput, "output_info", "O", "", "提取info中的字符集txt保存路径")

}
