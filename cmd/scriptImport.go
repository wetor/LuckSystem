/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/go-restruct/restruct"
	"lucksystem/charset"
	"lucksystem/game"
	"lucksystem/game/enum"

	"github.com/spf13/cobra"
)

// scriptImportCmd represents the scriptImportCmd command
var scriptImportCmd = &cobra.Command{
	Use:   "import",
	Short: "导入反编译的脚本",
	Run: func(cmd *cobra.Command, args []string) {
		restruct.EnableExprBeta()
		game.ScriptBlackList = append(game.ScriptBlackList, strings.Split(ScriptBlackList, ",")...)

		// PATCH YOREMI: Auto-detect game name from opcode file path (same as decompile)
		gameName := "Custom"
		if ScriptPlugin == "" && ScriptOpcode != "" {
			gameName = detectGameName(ScriptOpcode)
			if gameName != "Custom" {
				fmt.Printf("Auto-detected game: %s (from opcode path)\n", gameName)
			}
		}

		g := game.NewGame(&game.GameOptions{
			GameName:   gameName,
			PluginFile: ScriptPlugin,
			OpcodeFile: ScriptOpcode,
			Coding:     charset.Charset(Charset),
			Mode:       enum.VMRunImport,
		})
		g.LoadScriptResources(ScriptSource)
		g.ImportScript(ScriptImportDir, ScriptNoSubDir)
		g.RunScript()
		g.ImportScriptWrite(ScriptImportOutput)

	},
}

func init() {
	scriptCmd.AddCommand(scriptImportCmd)

	scriptImportCmd.Flags().StringVarP(&ScriptImportDir, "input", "i", "output", "输出的反编译脚本路径")
	scriptImportCmd.Flags().StringVarP(&ScriptImportOutput, "output", "o", "SCRIPT.PAK.out", "输出的SCRIPT.PAK文件")

	scriptImportCmd.MarkFlagsRequiredTogether("input", "output")
}
