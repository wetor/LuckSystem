package main

import (
	"fmt"
	vm "lucascript/game/VM"
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
	script.GameName = "LB_EN"
	script.Version = 3
	script.Read()
	//utils.Debug = utils.DebugNone
	vm := vm.NewVM(&script)
	err := vm.LoadOpcode("data/LB_EN/OPCODE.txt")
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Println(vm.Script.CodeNum)
	vm.Run()
	fmt.Println(vm.Context.Variable.ValueMap)
}

func TestSP(t *testing.T) {
	restruct.EnableExprBeta()
	var script script.ScriptFile
	script.FileName = "data/SP/SCRIPT/10_日常0729"
	script.GameName = "SP"
	script.Version = 3
	script.Read()

	vm := vm.NewVM(&script)
	err := vm.LoadOpcode("data/SP/OPCODE.txt")
	// game := game.NewGame("SP")
	// err := game.LoadOpcode("data/SP/OPCODE.txt")
	//game.Debug = true
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	vm.Run()
}
