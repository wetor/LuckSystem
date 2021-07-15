package context

import (
	"lucascript/game/engine"
	"lucascript/game/variable"
	"lucascript/script"
)

type Context struct {
	Script *script.ScriptFile
	// 运行时变量存储
	Variable *variable.VariableStore

	// 引擎前端
	Engine *engine.Engine
	// 当前下标
	CIndex int
	// 下一步执行下标
	CNext int

	// 等待按键
	KeyPress chan int
	// 等待阻塞
	ChanEIP chan int
}

// Code 获取当前code
func (ctx *Context) Code() *script.CodeLine {
	return ctx.Script.Codes[ctx.CIndex]
}
