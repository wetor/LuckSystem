/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/golang/glog"
	"lucksystem/czimage"
	"os"

	"github.com/spf13/cobra"
)

// imageImportCmd represents the imageImport command
var imageImportCmd = &cobra.Command{
	Use:   "import",
	Short: "导入png图像到cz文件中",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("imageImport called")
		cz := czimage.LoadCzImageFile(CzSource)
		out, err := os.Create(CzOutput)
		if err != nil {
			glog.Fatalln(err)
		}
		defer out.Close()

		f, err := os.Open(CzInput)
		if err != nil {
			return
		}
		defer f.Close()
		err = cz.Import(f, Fill)
		if err != nil {
			glog.Fatalln(err)
		}
		err = cz.Write(out)
		if err != nil {
			glog.Fatalln(err)
		}
	},
}

var (
	Fill bool // 填充为原大小，仅cz1支持
)

func init() {
	imageCmd.AddCommand(imageImportCmd)
	imageImportCmd.Flags().BoolVarP(&Fill, "fill", "f", false, "图像尺寸填充为与source一致，仅支持cz1")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imageImportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imageImportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
