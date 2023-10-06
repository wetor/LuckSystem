package runtime

import (
	"fmt"
	"os"
	"strings"

	"lucksystem/charset"
	"lucksystem/game/engine"
	"lucksystem/game/enum"
	"lucksystem/script"
)

type Runtime struct {
	Script *script.Script

	*GlobalGoto

	OpcodeMap     map[uint8]string
	ExprCharset   charset.Charset
	TextCharset   charset.Charset
	DefaultExport bool

	// 引擎前端
	Engine *engine.Engine

	// 当前下标
	CIndex int

	// 等待阻塞
	ChanEIP chan int

	// 运行模式
	RunMode enum.VMRunMode
}

func NewRuntime(mode enum.VMRunMode) *Runtime {
	ctx := &Runtime{
		GlobalGoto: NewGlobalGoto(),
		Engine:     &engine.Engine{},
		ChanEIP:    make(chan int),
		RunMode:    mode,
	}
	return ctx
}

func (ctx *Runtime) SwitchScript(scr *script.Script) {
	ctx.Script = scr
}

func (ctx *Runtime) Init(exprCharset, textCharset charset.Charset, defaultExport bool) {
	ctx.ExprCharset = exprCharset
	ctx.TextCharset = textCharset
	ctx.DefaultExport = defaultExport
}

func (ctx *Runtime) LoadOpcode(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
	}
	strlines := strings.Split(string(data), "\n")
	ctx.OpcodeMap = make(map[uint8]string, len(strlines)+1)
	for i, line := range strlines {
		line = strings.Replace(line, "\r", "", -1)
		line = strings.Replace(line, "\n", "", -1)
		line = strings.Replace(line, " ", "", -1)
		ctx.OpcodeMap[uint8(i)] = line
	}
}

func (ctx *Runtime) Opcode(index uint8) string {
	if op, ok := ctx.OpcodeMap[index]; ok {
		return op
	}
	return fmt.Sprintf("0x%02X", index)
}

// Code 获取当前code
func (ctx *Runtime) Code() *script.CodeLine {
	return ctx.Script.Codes[ctx.CIndex]
}
