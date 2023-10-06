/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/go-restruct/restruct"
	"lucksystem/charset"
	"lucksystem/game"
	"lucksystem/game/enum"

	"github.com/spf13/cobra"
)

// scriptDecompileCmd represents the scriptDecompileCmd command
var scriptDecompileCmd = &cobra.Command{
	Use:   "decompile",
	Short: "反编译脚本",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("scriptExtract called")
		restruct.EnableExprBeta()
		g := game.NewGame(&game.GameOptions{
			GameName:   "Custom",
			PluginFile: ScriptPlugin,
			OpcodeFile: ScriptOpcode,
			Coding:     charset.Charset(Charset),
			Mode:       enum.VMRunExport,
		})
		g.LoadScriptResources(ScriptSource)
		g.RunScript()

		g.ExportScript(ScriptExportDir)
	},
}

func init() {
	scriptCmd.AddCommand(scriptDecompileCmd)

	scriptDecompileCmd.Flags().StringVarP(&ScriptExportDir, "output", "o", "output", "反编译输出路径（目录不能有SCRIPT.PAK同名文件）")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imageExportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imageExportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
