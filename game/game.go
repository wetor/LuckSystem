package game

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/golang/glog"
	"lucksystem/charset"
	"lucksystem/game/VM"
	"lucksystem/game/enum"
	"lucksystem/pak"
	"lucksystem/script"
)

const (
	ResScript    = "SCRIPT.PAK"
	ResScriptExt = ".txt"
)

var ScriptBlackList = []string{"TEST", "_VOICEOTHER", "_VARNAME", "_VARNUM", "_CGMODE", "_SCR_LABEL", "_VOICE_PARAM", "_BUILD_COUNT", "_TASK", "_BUILD_TIME", "_VARSTRNAME"}

type GameOptions struct {
	GameName     string
	PluginFile   string
	OpcodeFile   string
	Version      int
	ResourcesDir string
	Coding       charset.Charset
	Mode         enum.VMRunMode
}

type Game struct {
	Version      int
	ResourcesDir string
	Coding       charset.Charset
	Resources    map[string]*pak.Pak

	VM         *VM.VM
	ScriptList []string
}

func NewGame(opt *GameOptions) *Game {
	if opt.Coding == "" {
		opt.Coding = charset.UTF_8
	}
	game := &Game{
		Version:      opt.Version,
		ResourcesDir: opt.ResourcesDir,
		Coding:       opt.Coding,
		Resources:    make(map[string]*pak.Pak),
		VM: VM.NewVM(&VM.Options{
			GameName:   opt.GameName,
			Mode:       opt.Mode,
			PluginFile: opt.PluginFile,
		}),
	}
	if len(opt.OpcodeFile) > 0 {
		game.VM.LoadOpcode(opt.OpcodeFile)
	}
	return game
}

func (g *Game) LoadResources() {
	g.LoadScriptResources(filepath.Join(g.ResourcesDir, ResScript))
	g.load()
}

func (g *Game) LoadScriptResources(file string) {
	g.Resources[ResScript] = pak.LoadPak(file, g.Coding)
	g.load()
}

// PATCH YOREMI: isValidScript checks if entry data looks like a valid LucaSystem script.
// Invalid entries (data tables, mini-game data, etc.) have firstLen=0 or firstLen < 4,
// which would cause restruct.Unpack to panic with size underflow.
// A valid CodeLine needs at least 4 bytes: 2 (Len) + 1 (Opcode) + 1 (FixedFlag).
func isValidScript(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	firstLen := binary.LittleEndian.Uint16(data[0:2])
	// Minimum valid CodeLine length is 4 (header only, no params)
	// Also reject suspiciously large values that would exceed the data
	if firstLen < 4 {
		return false
	}
	return true
}

// PATCH YOREMI: safeLoadScript wraps script.LoadScript with panic recovery
// to gracefully skip scripts that cause parsing errors.
func safeLoadScript(opts *script.LoadOptions) (scr *script.Script, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to parse script '%s': %v", opts.Entry.Name, r)
			scr = nil
		}
	}()
	scr = script.LoadScript(opts)
	return scr, nil
}

func (g *Game) load() {
	var err error
	for key, p := range g.Resources {
		switch key {
		case ResScript:
			var entry *pak.Entry
			for i := 1; i <= int(p.FileCount); i++ {
				entry, err = p.GetById(i)
				if err != nil {
					panic(err)
				}
				if !ScriptCanLoad(entry.Name) {
					glog.V(6).Infoln("Pass", entry.Name)
					continue
				}

				// PATCH YOREMI: Validate script data before attempting to parse.
				// Some PAK entries (e.g. SEEN8500, SEEN8501 in Little Busters)
				// are data tables, not scripts. Their firstLen=0 causes
				// restruct.Unpack to panic with size underflow.
				if !isValidScript(entry.Data) {
					glog.Warningf("Skipping invalid script entry '%s' (firstLen < 4, likely data table)\n", entry.Name)
					continue
				}

				glog.V(6).Infof("%v %v\n", entry.Name, len(entry.Data))

				// PATCH YOREMI: Use safe loading with panic recovery
				scr, loadErr := safeLoadScript(&script.LoadOptions{
					Entry: entry,
				})
				if loadErr != nil {
					glog.Warningf("Skipping script '%s': %v\n", entry.Name, loadErr)
					continue
				}

				g.VM.LoadScript(scr, false)
				g.ScriptList = append(g.ScriptList, scr.Name)
				g.VM.ScriptNames[scr.Name] = struct{}{}
				glog.V(6).Infoln(scr.CodeNum)
			}
		}
	}
}

func (g *Game) RunScript() {
	for _, name := range g.ScriptList {
		g.VM.SwitchScript(name)
		g.VM.Run()
	}
	for _, name := range g.ScriptList {
		labels, gotos := g.VM.GetMaps(name)
		g.VM.Scripts[name].SetGlobalLabel(labels)
		g.VM.Scripts[name].SetGlobalGoto(gotos)
	}

}

func (g *Game) ExportScript(dir string, noSubDir bool) {
	if !noSubDir {
		dir = path.Join(dir, ResScript)
	}
	exist, isDir := IsExistDir(dir)
	if exist && !isDir {
		panic("已存在同名文件")
	}
	if !exist {
		os.MkdirAll(dir, os.ModePerm)
	}
	for _, name := range g.ScriptList {
		f, err := os.Create(path.Join(dir, name+ResScriptExt))
		if err != nil {
			panic(err)
		}
		err = g.VM.Scripts[name].Export(f)
		if err != nil {
			panic(err)
		}
		f.Close()
	}

	//for i := 1; i < len(g.VM.GlobalLabelGoto)+1; i++ {
	//	fmt.Println(i, " ", g.VM.GlobalLabelGoto[i])
	//}
}

func (g *Game) ImportScript(dir string, noSubDir bool) {
	if !noSubDir {
		dir = path.Join(dir, ResScript)
	}
	for _, name := range g.ScriptList {
		f, err := os.Open(path.Join(dir, name+ResScriptExt))
		if err != nil {
			panic(err)
		}
		err = g.VM.Scripts[name].Import(f)
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}

func (g *Game) ImportScriptWrite(out string) {
	for _, name := range g.ScriptList {
		g.VM.AddGlobalLabelMap(g.VM.Scripts[name].IGlobalLabelMap)
	}
	var err error
	for _, name := range g.ScriptList {
		g.VM.Scripts[name].SetImportGlobalLabel(g.VM.IGlobalLabelMap)
		w := bytes.NewBuffer(nil)
		err = g.VM.Scripts[name].Write(w)
		if err != nil {
			panic(err)
		}
		err = g.Resources[ResScript].Set(name, w)
		if err != nil {
			panic(err)
		}
		w.Reset()
	}

	f, _ := os.Create(out)
	g.Resources[ResScript].Write(f)
	f.Close()
}
