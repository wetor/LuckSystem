package main

import (
	"fmt"
	vm "lucascript/game/VM"
	"lucascript/game/context"
	"lucascript/script"
	"lucascript/utils"
	"strconv"
	"testing"

	"github.com/go-restruct/restruct"
)

func Test11(t *testing.T) {
	fmt.Println(strconv.ParseInt("c9", 16, 32))
}

func TestLB_EN(t *testing.T) {
	restruct.EnableExprBeta()

	script := script.NewScript(script.ScriptFileOptions{
		FileName: "data/LB_EN/SCRIPT/SEEN0514",
		GameName: "LB_EN",
		Version:  3,
	})

	script.Read()
	utils.Debug = utils.DebugNone
	vm := vm.NewVM(script, context.VMRunExport)
	err := vm.LoadOpcode("data/LB_EN/OPCODE.txt")
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	vm.Run()
	script.Export("data/LB_EN/SCRIPT/SEEN0514.txt")

}

func TestSP(t *testing.T) {
	restruct.EnableExprBeta()
	script := script.NewScript(script.ScriptFileOptions{
		FileName: "data/SP/SCRIPT/10_日常0729",
		GameName: "SP",
		Version:  3,
	})

	script.Read()
	utils.Debug = utils.DebugNone
	vm := vm.NewVM(script, context.VMRunExport)
	err := vm.LoadOpcode("data/SP/OPCODE.txt")
	// game := game.NewGame("SP")
	// err := game.LoadOpcode("data/SP/OPCODE.txt")

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	vm.Run()
	//fmt.Println(vm.Context.Variable.ValueMap)
	script.Export("data/SP/SCRIPT/10_日常0729.txt")
}
