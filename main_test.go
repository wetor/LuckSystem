package main

import (
	"flag"
	"fmt"
	"lucksystem/game/VM"
	"lucksystem/game/enum"
	"lucksystem/script"
	"os"
	"strconv"
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
	n, err := strconv.Atoi("123a")
	fmt.Println(n, err)
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

	script := script.LoadScript(
		"data/LB_EN/SCRIPT/SEEN2005",
		"LB_EN",
		3,
	)

	script.Read()
	vm := VM.NewVM(script, enum.VMRunExport)
	err := vm.LoadOpcode("data/LB_EN/OPCODE.txt")
	if err != nil {
		fmt.Println(err)
	}
	vm.Run()
	f, _ := os.Create("data/LB_EN/TXT/SEEN2005.txt")
	defer f.Close()
	script.Export(f)

}

func TestLoadLB_EN(t *testing.T) {
	restruct.EnableExprBeta()
	script := script.LoadScript(
		"data/LB_EN/SCRIPT/SEEN2005",
		"LB_EN",
		3,
	)

	script.Read()
	f, _ := os.Open("data/LB_EN/TXT/SEEN2005.txt")
	defer f.Close()
	err := script.Import(f)
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
	sf, _ := os.Create(script.FileName + ".out")
	defer sf.Close()
	err = script.Write(sf)
	if err != nil {
		fmt.Println(err)
	}
}
func TestSP(t *testing.T) {
	restruct.EnableExprBeta()

	var err error

	script := script.LoadScript(
		"data/SP/SCRIPT/10_日常0729",
		"SP",
		3,
	)

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
	f, _ := os.Create("data/SP/TXT/10_日常0729.txt")
	defer f.Close()
	script.Export(f)
}

func TestLoadSP(t *testing.T) {
	restruct.EnableExprBeta()
	script := script.LoadScript(
		"data/SP/SCRIPT/10_日常0729",
		"SP",
		3,
	)

	script.Read()
	f, _ := os.Open("data/SP/TXT/10_日常0729.txt")
	defer f.Close()
	err := script.Import(f)
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
	sf, _ := os.Create(script.FileName + ".out")
	defer sf.Close()
	err = script.Write(sf)
	if err != nil {
		fmt.Println(err)
	}
}
