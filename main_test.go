package main

import (
	"fmt"
	"lucascript/game"
	"lucascript/script"
	"testing"

	"github.com/go-restruct/restruct"
)

func Test11(t *testing.T) {
	fmt.Println(2&1, 3&1, 4&1, 7&1)
}

func TestLB_EN(t *testing.T) {
	restruct.EnableExprBeta()
	var script script.ScriptFile
	script.FileName = "data/LB_EN/SCRIPT/SEEN0514"
	script.Version = 3
	script.Read()

	game := game.NewGame("LB_EN")
	err := game.LoadOpcode("data/LB_EN/OPCODE.txt")
	//game.Debug = true
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	game.Run(script.Code)
}

func TestSP(t *testing.T) {
	restruct.EnableExprBeta()
	var script script.ScriptFile
	script.FileName = "data/SP/SCRIPT/10_日常0729"
	script.Version = 3
	script.Read()

	game := game.NewGame("SP")
	err := game.LoadOpcode("data/SP/OPCODE.txt")
	//game.Debug = true
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	game.Run(script.Code)
}
