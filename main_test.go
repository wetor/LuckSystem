package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"testing"

	"lucksystem/charset"
	"lucksystem/game"
	"lucksystem/game/VM"
	"lucksystem/game/enum"
	"lucksystem/script"

	"github.com/go-restruct/restruct"
)

func TestMain(m *testing.M) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "log")
	flag.Set("v", "5")
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

	script := script.LoadScript(&script.LoadOptions{
		Filename: "data/LB_EN/SCRIPT/SEEN2005",
	})

	vm := VM.NewVM(&VM.Options{
		GameName: "LB_EN",
		Mode:     enum.VMRunExport,
	})
	vm.LoadScript(script, true)
	vm.LoadOpcode("data/LB_EN/OPCODE.txt")

	vm.Run()
	f, _ := os.Create("data/LB_EN/TXT/SEEN2005.txt")
	defer f.Close()
	script.Export(f)

}

func TestLoadLB_EN(t *testing.T) {
	restruct.EnableExprBeta()
	script := script.LoadScript(&script.LoadOptions{
		Filename: "data/LB_EN/SCRIPT/SEEN2005",
	})

	f, _ := os.Open("data/LB_EN/TXT/SEEN2005.txt")
	defer f.Close()
	err := script.Import(f)
	if err != nil {
		fmt.Println(err)
	}

	vm := VM.NewVM(&VM.Options{
		GameName: "LB_EN",
		Mode:     enum.VMRunImport,
	})
	vm.LoadScript(script, true)

	vm.LoadOpcode("data/LB_EN/OPCODE.txt")

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

	script := script.LoadScript(&script.LoadOptions{
		Filename: "C:/Users/wetor/Desktop/Prototype/SCRIPT.PAK_unpacked/10_日常0729",
	})

	// entry, err := pak.Get("10_日常0730")
	// if err != nil {
	// 	fmt.Println(err)
	// 	panic(err)
	// }
	// script.ReadByEntry(entry)
	vm := VM.NewVM(&VM.Options{
		GameName: "SP",
		Mode:     enum.VMRunExport,
	})
	vm.LoadScript(script, true)
	vm.LoadOpcode("data/SP/OPCODE.txt")
	// game := game.NewGame("SP")
	// err := game.LoadOpcode("data/SP/OPCODE.txt")

	vm.Run()
	//fmt.Println(vm.Runtime.Variable.ValueMap)
	f, _ := os.Create("C:/Users/wetor/Desktop/Prototype/SCRIPT.PAK_unpacked/TXT/10_日常0729.txt")
	defer f.Close()
	script.Export(f)
}

func TestPlugin(t *testing.T) {
	restruct.EnableExprBeta()
	script := script.LoadScript(&script.LoadOptions{
		Filename: "C:/Users/wetor/Desktop/Prototype/SCRIPT.PAK_unpacked/TXT/_称号_CS用処理",
	})

	vm := VM.NewVM(&VM.Options{
		GameName:   "SP",
		Mode:       enum.VMRunExport,
		PluginFile: "data/SP.py",
	})
	vm.LoadScript(script, true)
	vm.LoadOpcode("data/SP/OPCODE.txt")

	vm.Run()
	for i := 1; i < len(vm.GlobalLabelGoto)+1; i++ {
		fmt.Println(i, " ", vm.GlobalLabelGoto[i])
	}

	f, _ := os.Create("C:/Users/wetor/Desktop/Prototype/SCRIPT.PAK_unpacked/TXT/_称号_CS用処理.txt")
	defer f.Close()
	script.Export(f)
}

func TestPluginImport(t *testing.T) {
	restruct.EnableExprBeta()
	var err error
	script := script.LoadScript(&script.LoadOptions{
		Filename: "C:/Users/wetor/Desktop/Prototype/SCRIPT.PAK_unpacked/TXT/10_日常0729",
	})

	f, _ := os.Open("C:/Users/wetor/Desktop/Prototype/SCRIPT.PAK_unpacked/TXT/10_日常0729.txt")
	defer f.Close()
	err = script.Import(f)
	if err != nil {
		fmt.Println(err)
	}

	vm := VM.NewVM(&VM.Options{
		GameName:   "SP",
		Mode:       enum.VMRunImport,
		PluginFile: "data/SP.py",
	})
	vm.LoadScript(script, true)
	vm.LoadOpcode("data/SP/OPCODE.txt")

	vm.Run()
	sf, _ := os.Create(script.FileName + ".out")
	defer sf.Close()
	err = script.Write(sf)
	if err != nil {
		fmt.Println(err)
	}
}

func TestLoadSP(t *testing.T) {
	restruct.EnableExprBeta()
	script := script.LoadScript(&script.LoadOptions{
		Filename: "C:/Users/wetor/Desktop/Prototype/SCRIPT.PAK_unpacked/TXT/10_日常0729",
	})

	f, _ := os.Open("data/SP/TXT/10_日常0729.txt")
	defer f.Close()
	err := script.Import(f)
	if err != nil {
		fmt.Println(err)
	}

	vm := VM.NewVM(&VM.Options{
		GameName: "SP",
		Mode:     enum.VMRunImport,
	})
	vm.LoadScript(script, true)
	vm.LoadOpcode("data/SP/OPCODE.txt")

	vm.Run()
	sf, _ := os.Create(script.FileName + ".out")
	defer sf.Close()
	err = script.Write(sf)
	if err != nil {
		fmt.Println(err)
	}
}

func TestLoopersExportScript(t *testing.T) {
	// 反编译LOOPERS SCRIPT.PAK
	restruct.EnableExprBeta()
	g := game.NewGame(&game.GameOptions{
		GameName:   "LOOPERS",
		PluginFile: "C:/Users/wetor/GolandProjects/LuckSystem/data/LOOPERS.py",
		OpcodeFile: "C:/Users/wetor/GolandProjects/LuckSystem/data/LOOPERS.txt",
		Coding:     charset.UTF_8,
		Mode:       enum.VMRunExport,
	})
	g.LoadScriptResources("D:\\Game\\LOOPERS\\LOOPERS\\files\\src\\SCRIPT.PAK")
	g.RunScript()

	g.ExportScript("D:\\Game\\LOOPERS\\LOOPERS\\files\\Export")

}

func TestLoopersImportScript(t *testing.T) {
	// LOOPERS修改后导入 SCRIPT.PAK
	restruct.EnableExprBeta()
	g := game.NewGame(&game.GameOptions{
		GameName:   "LOOPERS",
		PluginFile: "C:/Users/wetor/GolandProjects/LuckSystem/data/LOOPERS.py",
		OpcodeFile: "C:/Users/wetor/GolandProjects/LuckSystem/data/LOOPERS.txt",
		Coding:     charset.UTF_8,
		Mode:       enum.VMRunImport,
	})
	g.LoadScriptResources("D:\\Game\\LOOPERS\\LOOPERS\\files\\src\\SCRIPT.PAK")
	g.ImportScript("D:\\Game\\LOOPERS\\LOOPERS\\files\\Export")
	g.RunScript()

	g.ImportScriptWrite("D:\\Game\\LOOPERS\\LOOPERS\\files\\Import\\SCRIPT.PAK")

	os.Rename("D:\\Game\\LOOPERS\\LOOPERS\\files\\Import\\SCRIPT.PAK", "D:\\Game\\LOOPERS\\LOOPERS\\files\\SCRIPT.PAK")
}
