package game

import (
	"flag"
	"os"
	"testing"

	"lucksystem/charset"
	"lucksystem/game/enum"

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

func TestLoadPak(t *testing.T) {
	restruct.EnableExprBeta()
	game := NewGame(&GameOptions{
		GameName:     "SP",
		PluginFile:   "C:/Users/wetor/GolandProjects/LuckSystem/data/SP.py",
		OpcodeFile:   "C:/Users/wetor/GolandProjects/LuckSystem/data/SP/OPCODE.txt",
		ResourcesDir: "C:/Users/wetor/Desktop/Prototype",
		Coding:       charset.ShiftJIS,
		Mode:         enum.VMRunExport,
	})
	game.LoadResources()
	game.RunScript()

	game.ExportScript("C:/Users/wetor/Desktop/Prototype/Export")

}

func TestLoadPak2(t *testing.T) {
	restruct.EnableExprBeta()
	game := NewGame(&GameOptions{
		GameName:     "SP",
		PluginFile:   "C:/Users/wetor/GolandProjects/LuckSystem/data/SP.py",
		OpcodeFile:   "C:/Users/wetor/GolandProjects/LuckSystem/data/SP/OPCODE.txt",
		ResourcesDir: "C:/Users/wetor/Desktop/Prototype",
		Coding:       charset.ShiftJIS,
		Mode:         enum.VMRunImport,
	})
	game.LoadResources()
	game.ImportScript("C:/Users/wetor/Desktop/Prototype/Export")
	game.RunScript()

	game.ImportScriptWrite("C:/Users/wetor/Desktop/Prototype/Import/SCRIPT.PAK")
}
