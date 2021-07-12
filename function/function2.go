package function

import (
	"lucascript/utils"
	"reflect"
	"strings"
)

type MESSAGE struct {
}

func (f *MESSAGE) Call(params ...interface{}) int {
	if len(params) != 2 {
		return 0
	}
	name := reflect.TypeOf(f).String()[10:] // remove "*function."

	voiceId := params[0].(uint16)
	str := params[1].(string)
	utils.Logf(`%s (%d, "%s")`, name, voiceId, str)
	return 0 // 向下执行
}

type SELECT struct {
}

func (f *SELECT) Call(params ...interface{}) int {
	if len(params) != 1 {
		return 0
	}
	name := reflect.TypeOf(f).String()[10:] // remove "*function."

	selectStr := strings.Split(params[0].(string), "$d")

	selectID := 1
	utils.Logf(`%s (%v) %d`, name, selectStr, selectID)

	return selectID // 向下执行
}
