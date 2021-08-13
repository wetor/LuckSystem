package game

import (
	"lucksystem/charset"
	"lucksystem/game/enum"
	"testing"

	"github.com/go-restruct/restruct"
)

func TestLoadPak(t *testing.T) {
	restruct.EnableExprBeta()
	game := NewGame(&GameOptions{
		GameName:     "SP",
		Version:      3,
		ResourcesDir: "../data/SP",
		Coding:       charset.ShiftJIS,
		Mode:         enum.VMRun,
	})
	game.LoadResources()
}
