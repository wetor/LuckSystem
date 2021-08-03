package main

import (
	"fmt"
	"lucascript/charset"
	vm "lucascript/game/VM"
	"lucascript/game/enum"
	"lucascript/pak"
	"lucascript/script"
	"lucascript/utils"
	"strconv"
	"testing"

	"github.com/go-restruct/restruct"
)

func Test11(t *testing.T) {
	fmt.Println(strconv.ParseInt("c9", 16, 32))
}

type lenString string

func Test22(t *testing.T) {
	var val1 lenString = "test"
	val2 := "test2"
	iface := []interface{}{val1, val2}
	for _, i := range iface {
		switch val := i.(type) {
		case string:
			fmt.Println("string", val)
		case lenString:
			fmt.Println("lenString", val)
		}
	}
}
func TestLB_EN(t *testing.T) {
	restruct.EnableExprBeta()

	script := script.NewScript(script.ScriptFileOptions{
		FileName: "data/LB_EN/SCRIPT/SEEN0513",
		GameName: "LB_EN",
		Version:  3,
	})

	script.Read()
	utils.Debug = utils.DebugNone
	vm := vm.NewVM(script, enum.VMRunExport)
	err := vm.LoadOpcode("data/LB_EN/OPCODE.txt")
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	vm.Run()
	script.Export("data/LB_EN/TXT/SEEN0513.txt")

}

func TestSP(t *testing.T) {
	restruct.EnableExprBeta()

	pak := pak.NewPak(&pak.PakFileOptions{
		FileName: "data/SP/SCRIPT.PAK",
		Coding:   charset.ShiftJIS,
	})
	err := pak.Open()
	if err != nil {
		fmt.Println(err)
	}

	script := script.NewScript(script.ScriptFileOptions{
		FileName: "data/SP/SCRIPT/10_日常0729",
		GameName: "SP",
		Version:  3,
	})

	entry, err := pak.Get("10_日常0730")
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	script.ReadByEntry(entry)
	utils.Debug = utils.DebugNone
	vm := vm.NewVM(script, enum.VMRunExport)
	err = vm.LoadOpcode("data/SP/OPCODE.txt")
	// game := game.NewGame("SP")
	// err := game.LoadOpcode("data/SP/OPCODE.txt")

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	vm.Run()
	//fmt.Println(vm.Context.Variable.ValueMap)
	script.Export("data/SP/TXT/10_日常0729.txt")
}

func TestLoadSP(t *testing.T) {
	restruct.EnableExprBeta()
	script := script.NewScript(script.ScriptFileOptions{
		FileName: "data/SP/SCRIPT/10_日常0729",
		GameName: "SP",
		Version:  3,
	})

	script.Read()
	utils.Debug = utils.DebugNone
	err := script.Import("data/SP/TXT/10_日常0729.txt")
	if err != nil {
		fmt.Println(err)
	}

	vm := vm.NewVM(script, enum.VMRunImport)
	err = vm.LoadOpcode("data/SP/OPCODE.txt")

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	vm.Run()
}

func foo(data []byte) {

	data = data[0:0]
	data = append(data, 1)
	data = append(data, 2)
	data = append(data, 3)
	fmt.Printf("%p %d\n", data, len(data))
}
func TestFuncName(t *testing.T) {
	data := make([]byte, 0, 10)
	data = append(data, 5)
	data = append(data, 6)
	foo(data)
	fmt.Printf("%p %d\n", data, len(data))
}
