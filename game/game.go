package game

import (
	"fmt"
	"lucascript/operation"
	"lucascript/script"
	"os"
	"reflect"
	"strings"
)

type Game struct {
	Debug     bool
	OpcodeMap map[uint8]string
	Operation operation.Operation
}

func NewGame(game string) *Game {
	if game == "LB_EN" {
		return &Game{
			Operation: operation.GetLB_EN(),
		}
	}
	return &Game{}
}

func (g *Game) Run(codes []*script.CodeLine) {
	var in []reflect.Value
	for i, code := range codes {
		opname, ok := g.OpcodeMap[code.Opcode]
		if !ok {
			if g.Debug {
				fmt.Println(i, "Opcode不存在", code.Opcode)
			}
			continue
		}
		fun := reflect.ValueOf(g.Operation).MethodByName(opname)
		if !fun.IsValid() {
			if g.Debug {
				fmt.Println(i, "Operation不存在", opname)
			}
			fun = reflect.ValueOf(g.Operation).MethodByName("UNDEFINE")
			in = make([]reflect.Value, 2)
			in[0] = reflect.ValueOf(code)
			in[1] = reflect.ValueOf(opname)
		} else {
			in = make([]reflect.Value, 1)
			in[0] = reflect.ValueOf(code)
		}
		result := fun.Call(in)
		fmt.Println(i, result[0].String())
	}
}
func (g *Game) LoadOpcode(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("os.ReadFile", err.Error())
		return err
	}
	strlines := strings.Split(string(data), "\n")
	g.OpcodeMap = make(map[uint8]string, len(strlines)+1)
	for i, line := range strlines {
		line = strings.Replace(line, "\r", "", -1)
		line = strings.Replace(line, "\n", "", -1)
		line = strings.Replace(line, " ", "", -1)
		g.OpcodeMap[uint8(i)] = line
	}
	return nil
}
