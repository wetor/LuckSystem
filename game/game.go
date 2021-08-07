package game

import (
	"fmt"
	"lucascript/charset"
	"lucascript/game/context"
	"lucascript/game/engine"
	"lucascript/game/enum"
	"lucascript/pak"
	"lucascript/script"
	"path/filepath"
)

const (
	ResScript = "SCRIPT.PAK"
	ResImage  = "BG.PAK"
)

type GameOptions struct {
	GameName     string
	Version      uint8
	ResourcesDir string
	Coding       charset.Charset
	Mode         enum.VMRunMode
}

type Game struct {
	GameName     string
	Version      uint8
	ResourcesDir string
	Coding       charset.Charset
	Resources    map[string]*pak.PakFile
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
	game.Resources = make(map[string]*pak.PakFile)
	game.Resources[ResScript] = pak.NewPak(&pak.PakFileOptions{
		FileName: filepath.Join(game.ResourcesDir, ResScript),
		Coding:   game.Coding,
	})
	return game
}

func (g *Game) LoadResources() {
	var err error
	for key, pak := range g.Resources {
		pak.Open()
		switch key {
		case ResScript:
			for _, entry := range pak.Files {
				fmt.Println(entry.Name)
				g.Context.Scripts[entry.Name], err = script.OpenScriptFile(entry)
				if err != nil {
					panic(err)
				}
			}
		}

	}
}
