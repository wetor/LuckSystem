package context

import (
	"lucksystem/game/engine"
	"lucksystem/game/enum"
	"lucksystem/game/variable"
	"lucksystem/script"
)

type Context struct {
	Scripts map[string]*script.Script
	// 运行时变量存储
	Variable *variable.VariableStore

	// 引擎前端
	Engine *engine.Engine
	// 当前脚本名
	CScriptName string
	// 当前下标
	CIndex int
	// 下一步执行下标
	CNext int

	// 等待按键
	KeyPress chan int
	// 等待阻塞
	ChanEIP chan int

	// 运行模式
	RunMode enum.VMRunMode
}

// Script 获取当前script
func (ctx *Context) Script() *script.Script {
	return ctx.Scripts[ctx.CScriptName]
}

// Code 获取当前code
func (ctx *Context) Code() *script.CodeLine {
	return ctx.Scripts[ctx.CScriptName].Codes[ctx.CIndex]
}
