/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// fontCreateCmd represents the fontCreate command
var fontCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "暂不支持",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fontCreate called")
	},
}

func init() {
	fontCmd.AddCommand(fontCreateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fontCreateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fontCreateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
