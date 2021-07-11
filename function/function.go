package function

import (
	"lucascript/paramter"
	"lucascript/utils"
	"reflect"
)

type HandlerFunc func() int

type Function interface {
	Call([]paramter.Paramter) int
}

type EQU struct {
}

func (f *EQU) Call(params []paramter.Paramter) int {
	if len(params) != 2 {
		return 0
	}
	name := reflect.TypeOf(f).String()[10:] // remove "*function."

	key := params[0].Value().(uint16)
	value := params[1].Value().(uint16)
	utils.Logf("%s #%d = %d", name, key, value)
	return 0 // 向下执行
}

type IFN struct {
}

func (f *IFN) Call(params []paramter.Paramter) int {
	if len(params) != 2 {
		return 0
	}
	name := reflect.TypeOf(f).String()[10:] // remove "*function."

	ifExprStr := params[0].Value().(string)
	jumpPos := params[1].Value().(uint32)
	utils.Logf("%s %s{goto %d}", name, ifExprStr, jumpPos)
	return 0 // 向下执行
}

type FARCALL struct {
}

func (f *FARCALL) Call(params []paramter.Paramter) int {
	if len(params) != 3 {
		return 0
	}
	name := reflect.TypeOf(f).String()[10:] // remove "*function."
	index := params[0].Value().(uint16)
	fileStr := params[1].Value().(string)
	jumpPos := params[2].Value().(uint32)
	utils.Logf("%s (%d) {goto \"%s\", %d}", name, index, fileStr, jumpPos)
	if false {
		return int(jumpPos)
	}
	return 0 // 向下执行
}

type GOTO struct {
}

func (f *GOTO) Call(params []paramter.Paramter) int {
	if len(params) != 1 {
		return 0
	}
	name := reflect.TypeOf(f).String()[10:] // remove "*function."

	jumpPos := params[0].Value().(uint32)
	utils.Logf("%s %d", name, jumpPos)
	if false {
		return int(jumpPos)
	}
	return 0 // 向下执行
}

type JUMP struct {
}

func (f *JUMP) Call(params []paramter.Paramter) int {
	if len(params) != 2 {
		return 0
	}
	name := reflect.TypeOf(f).String()[10:] // remove "*function."

	fileStr := params[0].Value().(string)
	jumpPos := params[1].Value().(uint32)
	utils.Logf("%s {goto \"%s\", %d}", name, fileStr, jumpPos)
	if false {
		return int(jumpPos)
	}
	return 0 // 向下执行
}
