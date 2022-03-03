package game

import (
	"github.com/golang/glog"
	"lucksystem/charset"
	"lucksystem/game/context"
	"lucksystem/game/engine"
	"lucksystem/game/enum"
	"lucksystem/pak"
	"lucksystem/script"
	"path/filepath"
)

const (
	ResScript = "SCRIPT.PAK"
	ResImage  = "BG.PAK"
)

var ScriptBlackList = []string{"_VARNUM", "_CGMODE", "_SCR_LABEL", "_VOICE_PARAM", "_BUILD_COUNT", "_TASK"}

type GameOptions struct {
	GameName     string
	Version      int
	ResourcesDir string
	Coding       charset.Charset
	Mode         enum.VMRunMode
}

type Game struct {
	GameName     string
	Version      int
	ResourcesDir string
	Coding       charset.Charset
	Resources    map[string]*pak.Pak
	Context      *context.Context
}

func NewGame(opt *GameOptions) *Game {
	game := &Game{}
	game.GameName = opt.GameName
	game.Version = opt.Version
	game.ResourcesDir = opt.ResourcesDir
	if opt.Coding != "" {
		game.Coding = opt.Coding
	} else {
		game.Coding = charset.UTF_8
	}
	game.Context = &context.Context{
		Engine:   &engine.Engine{},
		Scripts:  make(map[string]*script.ScriptFile),
		KeyPress: make(chan int),
		ChanEIP:  make(chan int),
		RunMode:  opt.Mode,
	}
	game.Resources = make(map[string]*pak.Pak)
	game.Resources[ResScript] = pak.LoadPak(
		filepath.Join(game.ResourcesDir, ResScript),
		game.Coding,
	)
	return game
}

func (g *Game) LoadResources() {
	var err error

	for key, p := range g.Resources {
		p.Open()
		switch key {
		case ResScript:
			var entry *pak.Entry
			for i := 0; i < int(p.FileCount); i++ {
				entry, err = p.GetById(i)
				if err != nil {
					panic(err)
				}
				if !ScriptCanLoad(entry.Name) {
					glog.V(4).Infoln("Pass", entry.Name)
					continue
				}
				glog.V(4).Infof("%v %v\n", entry.Name, len(entry.Data))
				scr, err := script.LoadScriptEntry(entry)
				if err != nil {
					panic(err)
				}
				g.Context.Scripts[entry.Name] = scr
				glog.V(4).Infoln(scr.CodeNum)
			}
		}

	}
}
