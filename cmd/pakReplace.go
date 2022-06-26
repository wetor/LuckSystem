/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/golang/glog"
	"lucksystem/charset"
	"lucksystem/pak"
	"os"

	"github.com/spf13/cobra"
)

// pakReplaceCmd represents the pakReplace command
var pakReplaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "替换Pak子文件",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pakReplace called")
		if len(PakSource) == 0 {
			fmt.Println("Error: required flag(s) \"source\" not set")
			return
		}
		p := pak.LoadPak(PakSource, charset.Charset(Charset))
		out, err := os.Create(PakOutput)
		if err != nil {
			glog.Fatalln(err)
		}
		defer out.Close()
		if fi, err := os.Stat(PakInput); err != nil {
			glog.Fatalln(err)
		} else if fi.IsDir() {
			err := p.Import(nil, "dir", PakInput)
			if err != nil {
				glog.Fatalln(err)
			}
		} else {
			f, err := os.Open(PakInput)
			if err != nil {
				glog.Fatalln(err)
			}
			defer f.Close()
			if PakList {
				err := p.Import(f, "list", nil)
				if err != nil {
					glog.Fatalln(err)
				}
			} else if PakId >= 0 {
				err := p.Import(f, "file", PakId)
				if err != nil {
					glog.Fatalln(err)
				}
			} else if len(PakName) > 0 {
				err := p.Import(f, "file", PakName)
				if err != nil {
					glog.Fatalln(err)
				}
			}

		}
		err = p.Write(out)
		if err != nil {
			glog.Fatalln(err)
		}
	},
}

var (
	PakList bool
)

func init() {
	pakCmd.AddCommand(pakReplaceCmd)

	pakReplaceCmd.Flags().BoolVarP(&PakList, "list", "l", false, "input是否为列表文件")
	pakReplaceCmd.Flags().IntVarP(&PakId, "id", "d", -1, "替换指定ID文件")
	pakReplaceCmd.Flags().StringVarP(&PakName, "name", "n", "", "替换指定文件名")

	pakReplaceCmd.MarkFlagsMutuallyExclusive("list", "id", "name")
}
