package main

import (
	"flag"
	"fmt"
	"lucksystem/game/VM"
	"lucksystem/game/enum"
	"lucksystem/script"
	"os"
	"testing"

	"github.com/go-restruct/restruct"
)

func TestMain(m *testing.M) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "log")
	flag.Set("v", "10")
	flag.Parse()

	ret := m.Run()
	os.Exit(ret)
}
func Test11(t *testing.T) {
	var offset uint32 = 33
	var BlockSize uint32 = 32
	if offset/BlockSize*BlockSize != offset {
		offset = (offset/BlockSize + 1) * BlockSize
	}
	fmt.Println(offset)
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

	script := script.NewScriptFile(script.ScriptFileOptions{
		FileName: "data/LB_EN/SCRIPT/SEEN2005",
		GameName: "LB_EN",
		Version:  3,
	})

	script.Read()
	vm := VM.NewVM(script, enum.VMRunExport)
	err := vm.LoadOpcode("data/LB_EN/OPCODE.txt")
	if err != nil {
		fmt.Println(err)
	}
	vm.Run()
	script.Export("data/LB_EN/TXT/SEEN2005.txt")

}

func TestLoadLB_EN(t *testing.T) {
	restruct.EnableExprBeta()
	script := script.NewScriptFile(script.ScriptFileOptions{
		FileName: "data/LB_EN/SCRIPT/SEEN2005",
		GameName: "LB_EN",
		Version:  3,
	})

	script.Read()
	err := script.Import("data/LB_EN/TXT/SEEN2005.txt")
	if err != nil {
		fmt.Println(err)
	}

	vm := VM.NewVM(script, enum.VMRunImport)
	err = vm.LoadOpcode("data/LB_EN/OPCODE.txt")

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	vm.Run()
	err = script.Write()
	if err != nil {
		fmt.Println(err)
	}
}
func TestSP(t *testing.T) {
	restruct.EnableExprBeta()

	var err error
	// pak := pak.NewPak(&pak.PakFileOptions{
	// 	FileName: "data/SP/SCRIPT.PAK",
	// 	Coding:   charset.ShiftJIS,
	// })
	// err = pak.Open()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	script := script.NewScriptFile(script.ScriptFileOptions{
		FileName: "data/SP/SCRIPT/10_日常0729",
		GameName: "SP",
		Version:  3,
	})

	// entry, err := pak.Get("10_日常0730")
	// if err != nil {
	// 	fmt.Println(err)
	// 	panic(err)
	// }
	// script.ReadByEntry(entry)
	script.Read()
	vm := VM.NewVM(script, enum.VMRunExport)
	err = vm.LoadOpcode("data/SP/OPCODE.txt")
	// game := game.NewGame("SP")
	// err := game.LoadOpcode("data/SP/OPCODE.txt")

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	vm.Run()
	//fmt.Println(vm.Context.Variable.ValueMap)
	script.Export("data/SP/TXT/10_日常0729.txt")
}

func TestLoadSP(t *testing.T) {
	restruct.EnableExprBeta()
	script := script.NewScriptFile(script.ScriptFileOptions{
		FileName: "data/SP/SCRIPT/10_日常0729",
		GameName: "SP",
		Version:  3,
	})

	script.Read()
	err := script.Import("data/SP/TXT/10_日常0729.txt")
	if err != nil {
		fmt.Println(err)
	}

	vm := VM.NewVM(script, enum.VMRunImport)
	err = vm.LoadOpcode("data/SP/OPCODE.txt")

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	vm.Run()
	err = script.Write()
	if err != nil {
		fmt.Println(err)
	}
}
