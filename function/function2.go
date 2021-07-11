package function

import (
	"lucascript/paramter"
	"lucascript/utils"
	"reflect"
)

type MESSAGE struct {
}

func (f *MESSAGE) Call(params []paramter.Paramter) int {
	if len(params) != 2 {
		return 0
	}
	name := reflect.TypeOf(f).String()[10:] // remove "*function."

	voiceId := params[0].Value().(uint16)
	str := params[1].Value().(string)
	utils.Logf(`%s (%d, "%s")`, name, voiceId, str)
	return 0 // 向下执行
}
