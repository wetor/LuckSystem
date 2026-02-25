/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-restruct/restruct"
	"lucksystem/charset"
	"lucksystem/game"
	"lucksystem/game/enum"

	"github.com/spf13/cobra"
)

// detectGameName tries to infer the game name from the opcode file path.
// For example, "data/LB_EN/OPCODE.txt" → "LB_EN", "data/SP/OPCODE.txt" → "SP".
// Returns "Custom" if no known game is detected.
func detectGameName(opcodePath string) string {
	if opcodePath == "" {
		return "Custom"
	}
	// Get the parent directory name of the opcode file
	dir := filepath.Dir(opcodePath)
	dirName := strings.ToUpper(filepath.Base(dir))

	knownGames := []string{"LB_EN", "SP"}
	for _, g := range knownGames {
		if dirName == g {
			return g
		}
	}
	return "Custom"
}

// scriptDecompileCmd represents the scriptDecompileCmd command
var scriptDecompileCmd = &cobra.Command{
	Use:   "decompile",
	Short: "反编译脚本",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("scriptExtract called")
		restruct.EnableExprBeta()
		game.ScriptBlackList = append(game.ScriptBlackList, strings.Split(ScriptBlackList, ",")...)

		// PATCH YOREMI: Auto-detect game name from opcode file path when no plugin is specified.
		// This allows "lucksystem script decompile -O data/LB_EN/OPCODE.txt" to work
		// without requiring an explicit -p plugin flag.
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
			Mode:       enum.VMRunExport,
		})
		g.LoadScriptResources(ScriptSource)
		g.RunScript()

		g.ExportScript(ScriptExportDir, ScriptNoSubDir)
	},
}

func init() {
	scriptCmd.AddCommand(scriptDecompileCmd)

	scriptDecompileCmd.Flags().StringVarP(&ScriptExportDir, "output", "o", "output", "反编译输出路径")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imageExportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imageExportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
