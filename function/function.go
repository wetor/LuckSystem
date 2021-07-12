// function 模拟器相关
// 此包内的方法与模拟器相关，负责画面的输出、按键的输入等操作
package function

import (
	"lucascript/utils"
	"reflect"
)

type HandlerFunc func()

type Function interface {
	Call(...interface{}) int
}

type EQU struct {
}

func (f *EQU) Call(params ...interface{}) int {
	if len(params) != 2 {
		return 0
	}
	name := reflect.TypeOf(f).String()[10:] // remove "*function."

	key := params[0].(uint16)
	value := params[1].(uint16)
	utils.Logf("%s #%d = %d", name, key, value)
	return 0 // 向下执行
}

type EQUN struct {
}

func (f *EQUN) Call(params ...interface{}) int {
	if len(params) != 2 {
		return 0
	}
	name := reflect.TypeOf(f).String()[10:] // remove "*function."

	key := params[0].(uint16)
	value := params[1].(uint16)
	utils.Logf("%s #%d = %d", name, key, value)
	return 0 // 向下执行
}

type IFN struct {
}

func (f *IFN) Call(params ...interface{}) int {
	if len(params) != 2 {
		return 0
	}
	name := reflect.TypeOf(f).String()[10:] // remove "*function."

	ifExprStr := params[0].(string)
	jumpPos := params[1].(uint32)
	utils.Logf("%s %s{goto %d}", name, ifExprStr, jumpPos)
	return 0 // 向下执行
}

type IFY struct {
}

func (f *IFY) Call(params ...interface{}) int {
	if len(params) != 2 {
		return 0
	}
	name := reflect.TypeOf(f).String()[10:] // remove "*function."

	ifExprStr := params[0].(string)
	jumpPos := params[1].(uint32)
	utils.Logf("%s %s{goto %d}", name, ifExprStr, jumpPos)
	return 0 // 向下执行
}

type FARCALL struct {
}

func (f *FARCALL) Call(params ...interface{}) int {
	if len(params) != 3 {
		return 0
	}
	name := reflect.TypeOf(f).String()[10:] // remove "*function."
	index := params[0].(uint16)
	fileStr := params[1].(string)
	jumpPos := params[2].(uint32)
	utils.Logf("%s (%d) {goto \"%s\", %d}", name, index, fileStr, jumpPos)
	return 0 // 向下执行
}

type GOTO struct {
}

func (f *GOTO) Call(params ...interface{}) int {
	if len(params) != 1 {
		return 0
	}
	name := reflect.TypeOf(f).String()[10:] // remove "*function."

	jumpPos := params[0].(uint32)
	utils.Logf("%s %d", name, jumpPos)

	return 0 // 向下执行
}

type JUMP struct {
}

func (f *JUMP) Call(params ...interface{}) int {
	if len(params) != 2 {
		return 0
	}
	name := reflect.TypeOf(f).String()[10:] // remove "*function."

	fileStr := params[0].(string)
	jumpPos := params[1].(uint32)
	utils.Logf("%s {goto \"%s\", %d}", name, fileStr, jumpPos)
	return 0 // 向下执行
}
