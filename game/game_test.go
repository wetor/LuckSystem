package game

import (
	"fmt"
	"lucascript/charset"
	"lucascript/game/enum"
	"lucascript/game/operater"
	"testing"

	"github.com/go-restruct/restruct"
)

func TestGame(t *testing.T) {
	by := []byte{0x60, 00, 0x52, 00, 0x69, 00, 0x6B, 00, 0x69, 0x00, 0x40, 00, 00, 0x00, 0x60, 00, 0x52, 00, 0x69, 00, 0x6B, 00, 0x69, 0x00, 0x40, 00}
	str, n := operater.DecodeString(by, 0, 0, charset.Unicode)
	fmt.Println(str, n)
	str, n = operater.DecodeString(by, n, 0, charset.Unicode)
	fmt.Println(str, n)

}
func TestGameJIS(t *testing.T) {
	// (@エサを所持している(20))
	by := []byte{0x28, 0x40, 0x83, 0x47, 0x83, 0x54, 0x82, 0xf0, 0x8f, 0x8a, 0x8e, 0x9d, 0x82, 0xb5, 0x82, 0xc4, 0x82, 0xa2, 0x82, 0xe9, 0x28, 0x32, 0x30, 0x29, 0x29, 00}

	str, n := operater.DecodeString(by, 0, 0, charset.ShiftJIS)
	fmt.Println(str, n)

}
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
