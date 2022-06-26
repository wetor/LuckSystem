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

// imageExportCmd represents the imageExport command
var imageExportCmd = &cobra.Command{
	Use:   "export",
	Short: "提取cz文件到png图片",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("imageExport called")
		cz := czimage.LoadCzImageFile(CzInput)
		out, err := os.Create(CzOutput)
		if err != nil {
			glog.Fatalln(err)
		}
		defer out.Close()
		err = cz.Export(out)
		if err != nil {
			glog.Fatalln(err)
		}
	},
}

func init() {
	imageCmd.AddCommand(imageExportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imageExportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imageExportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
