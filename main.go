package main

import (
	"fmt"
	"lucascript/game"
	"lucascript/script"

	"github.com/go-restruct/restruct"
)

func main() {
	restruct.EnableExprBeta()
	var script script.ScriptFile
	script.FileName = "data/LB_EN/SCRIPT/SEEN0513"
	script.Version = 3
	script.Read()

	game := game.NewGame("LB_EN")
	err := game.LoadOpcode("data/LB_EN/LB_EN.txt")
	//game.Debug = true
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	game.Run(script.Code)

	// fmt.Println(script.FileName, script.CodeNum)
	// for _, code := range script.Code {
	// 	fmt.Println(code.Opcode, code.InfoFlag, code.Info, code.CodeBytes)
	// }
}
