/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"lucksystem/charset"
)

// scriptCmd represents the script command
var scriptCmd = &cobra.Command{
	Use:   "script",
	Short: "LucaSystem script文件",
	Long: `LucaSystem script文件
无具体文件头，确定是LucaSystem引擎的游戏，SCRIPT.PAK中的文件即为script文件
其中'_VARNUM'、'_CGMODE'、'_SCR_LABEL'、'_VOICE_PARAM'、'_BUILD_COUNT'、'_TASK'等文件不支持单个解析`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("script called")
	},
}

var (
	ScriptOpcode       string
	ScriptPlugin       string
	ScriptSource       string
	ScriptExportDir    string
	ScriptImportDir    string
	ScriptImportOutput string
	ScriptNoSubDir     bool
	ScriptBlackList    string
	ScriptGameName     string
)

func init() {
	rootCmd.AddCommand(scriptCmd)

	scriptCmd.PersistentFlags().StringVarP(&ScriptSource, "source", "s", "SCRIPT.PAK", "SCRIPT.PAK文件")
	scriptCmd.PersistentFlags().StringVarP(&Charset, "charset", "c", string(charset.UTF_8), "PAK文件字符串编码")
	scriptCmd.PersistentFlags().StringVarP(&ScriptBlackList, "blacklist", "b", "", "额外的脚本黑名单")
	scriptCmd.PersistentFlags().StringVarP(&ScriptOpcode, "opcode", "O", "", "游戏的OPCODE文件")
	scriptCmd.PersistentFlags().StringVarP(&ScriptPlugin, "plugin", "p", "", "游戏OPCODE解析插件")
	scriptCmd.PersistentFlags().StringVarP(&ScriptGameName, "game", "g", "", "Game name (e.g. LB_EN, SP). Overrides auto-detection from OPCODE path")
	scriptCmd.PersistentFlags().BoolVarP(&ScriptNoSubDir, "no_subdir", "n", false, "输入和输出路径的不追加 '/SCRIPT.PAK/' 子目录")
	scriptCmd.MarkPersistentFlagRequired("source")
}
