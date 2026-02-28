/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/go-restruct/restruct"
	"lucksystem/charset"
	"lucksystem/game"
	"lucksystem/game/enum"
	"lucksystem/game/operator"

	"github.com/spf13/cobra"
)

// scriptImportCmd represents the scriptImportCmd command
var scriptImportCmd = &cobra.Command{
	Use:   "import",
	Short: "导入反编译的脚本",
	Run: func(cmd *cobra.Command, args []string) {
		restruct.EnableExprBeta()
		game.ScriptBlackList = append(game.ScriptBlackList, strings.Split(ScriptBlackList, ",")...)

		// PATCH YOREMI: Resolve game name from --game flag or auto-detect from OPCODE path.
		// Same logic as scriptDecompile.go — ensures MESSAGE/SELECT/BATTLE opcodes
		// are properly re-encoded during import (not treated as raw uint16).
		gameName := resolveGameName()

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

		// Print summary of undefined (non-text) opcodes that were skipped
		operator.PrintUndefinedOpcodeSummary()

		g.ImportScriptWrite(ScriptImportOutput)

	},
}

func init() {
	scriptCmd.AddCommand(scriptImportCmd)

	scriptImportCmd.Flags().StringVarP(&ScriptImportDir, "input", "i", "output", "输出的反编译脚本路径")
	scriptImportCmd.Flags().StringVarP(&ScriptImportOutput, "output", "o", "SCRIPT.PAK.out", "输出的SCRIPT.PAK文件")

	scriptImportCmd.MarkFlagsRequiredTogether("input", "output")
}
