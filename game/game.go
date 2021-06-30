package game

import (
	"fmt"
	"lucascript/operater"
	"lucascript/script"
	"lucascript/utils"
	"os"
	"reflect"
	"strings"
)

type Game struct {
	OpcodeMap map[uint8]string
	Operate   operater.Operater
	CodeLines []*script.CodeLine
	CodeIndex int // 当前指令序号
	EIP       int // 下一指令偏移
}

func NewGame(game string) *Game {
	switch game {
	case "LB_EN":
		return &Game{
			Operate: operater.GetLB_EN(),
		}
	case "SP":
		return &Game{
			Operate: operater.GetSP(),
		}
	}
	return &Game{}
}

// 在对EIP修改后调用，查找下一条具体指令，返回指令序号
func (g *Game) findCode(oldEIP int) int {
	index := g.CodeIndex
	if g.EIP > oldEIP {
		// 向下查找
		for index < len(g.CodeLines) && g.CodeLines[index].Pos < g.EIP {
			index++
		}
		if g.CodeLines[index].Pos == g.EIP {
			return index
		} else {
			panic(fmt.Sprintf("未找到跳转位置 [%d]%d -> %d", g.CodeIndex, oldEIP, g.EIP))
		}
	} else if g.EIP < oldEIP {
		// 向上查找
		for index >= 0 && g.CodeLines[index].Pos > g.EIP {
			index--
		}
		if g.CodeLines[index].Pos == g.EIP {
			return index
		} else {
			panic(fmt.Sprintf("未找到跳转位置 [%d]%d -> %d", g.CodeIndex, oldEIP, g.EIP))
		}
	} else {
		return index
	}

}

func (g *Game) getNextPos() int {
	if g.CodeIndex+1 >= len(g.CodeLines) {
		return 0
	} else {
		return g.CodeLines[g.CodeIndex+1].Pos
	}

}
func (g *Game) Run(codes []*script.CodeLine) {
	g.CodeLines = codes
	g.EIP = 0
	g.CodeIndex = 0

	var in []reflect.Value
	var code *script.CodeLine
	for {
		code = g.CodeLines[g.CodeIndex]
		opname, ok := g.OpcodeMap[code.Opcode]
		if !ok {
			utils.LogA(g.CodeIndex, "Opcode不存在", code.Opcode)
			g.CodeIndex++
			continue
		}
		operat := reflect.ValueOf(g.Operate).MethodByName(opname)
		if !operat.IsValid() {
			utils.LogA(g.CodeIndex, "Operation不存在", opname)
			// 方法未定义，调用UNDEFINE
			operat = reflect.ValueOf(g.Operate).MethodByName("UNDEFINE")
			in = make([]reflect.Value, 2)
			in[0] = reflect.ValueOf(code)
			in[1] = reflect.ValueOf(opname)
		} else {
			// 方法已定义，反射调用
			in = make([]reflect.Value, 1)
			in[0] = reflect.ValueOf(code)
		}
		fun := operat.Call(in) // 反射调用 operater，将CodeBytes解析为参数列表、并返回一个可执行func

		next := g.getNextPos() // 取得下一句位置

		if fun[0].Kind() == reflect.Func {
			eip := fun[0].Interface().(func() int)() // 调用，默认传递参数列表，取得跳转的位置
			if eip > 0 {                             // 为0则默认下一句
				next = eip
			}
		} else {
			utils.LogT(g.CodeIndex, fun[0].String())

		}
		utils.LogT("next:", next)

		if next == 0 || opname == "END" {
			break // 结束
		}
		g.EIP = next
		g.CodeIndex = g.findCode(code.Pos)
	}
}
func (g *Game) LoadOpcode(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		utils.Log("os.ReadFile", err.Error())
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
