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
	//script.FileName = "data/LB_EN/SCRIPT/SEEN0513"
	script.FileName = "data/SP/SCRIPT/10_日常0729"
	script.Version = 3
	script.Read()

	game := game.NewGame("SP")
	err := game.LoadOpcode("data/SP/SP.txt")

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
