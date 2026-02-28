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
	"lucksystem/game/operator"

	"github.com/spf13/cobra"
)

// detectGameName extracts the game name from the OPCODE file path.
// It first checks the parent directory name (e.g. "data/LB_EN/OPCODE.txt"),
// then falls back to searching anywhere in the path for a known game name.
// Known game names: LB_EN, SP (matching the switch in vm.go NewVM)
// Returns "Custom" if no known game name is found.
func detectGameName(opcodePath string) string {
	if opcodePath == "" {
		return "Custom"
	}

	knownGames := []string{"LB_EN", "SP"}

	// Strategy 1: Check parent directory name (most precise)
	dir := filepath.Dir(opcodePath)
	name := filepath.Base(dir)
	for _, g := range knownGames {
		if strings.EqualFold(name, g) {
			return g
		}
	}

	// Strategy 2: Search anywhere in the full path (case-insensitive)
	normalizedPath := filepath.ToSlash(opcodePath)
	upperPath := strings.ToUpper(normalizedPath)
	for _, g := range knownGames {
		if strings.Contains(upperPath, strings.ToUpper(g)) {
			return g
		}
	}

	return "Custom"
}

// resolveGameName determines the game name from explicit flag, auto-detection, or default.
func resolveGameName() string {
	// Priority 1: Explicit --game flag
	if ScriptGameName != "" {
		fmt.Printf("[INFO] Using game: %s (from --game flag)\n", ScriptGameName)
		return ScriptGameName
	}

	// Priority 2: Auto-detect from OPCODE path
	if ScriptPlugin == "" && ScriptOpcode != "" {
		gameName := detectGameName(ScriptOpcode)
		if gameName != "Custom" {
			fmt.Printf("[INFO] Auto-detected game: %s (from OPCODE path)\n", gameName)
			return gameName
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

		// PATCH YOREMI: Resolve game name from --game flag or auto-detect from OPCODE path.
		// This ensures MESSAGE/SELECT/BATTLE opcodes are properly decoded as text
		// instead of raw uint16 codepoints via the generic fallback.
		gameName := resolveGameName()

		g := game.NewGame(&game.GameOptions{
			GameName:   gameName,
			PluginFile: ScriptPlugin,
			OpcodeFile: ScriptOpcode,
			Coding:     charset.Charset(Charset),
			Mode:       enum.VMRunExport,
		})
		g.LoadScriptResources(ScriptSource)
		g.RunScript()

		// Print summary of undefined (non-text) opcodes that were skipped
		operator.PrintUndefinedOpcodeSummary()

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
