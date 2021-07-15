package engine

import (
	"lucascript/utils"
	"strings"
)

func (Engine) MESSAGE(params ...interface{}) int {
	if len(params) != 2 {
		panic("参数数量错误")
	}

	voiceId := params[0].(uint16)
	str := params[1].(string)
	utils.Logf(`MESSAGE (%d, "%s")`, voiceId, str)
	return 0 // 向下执行
}

func (Engine) SELECT(params ...interface{}) int {
	if len(params) != 1 {
		panic("参数数量错误")
	}

	selectStr := strings.Split(params[0].(string), "$d")

	selectID := 1
	utils.Logf(`SELECT (%v) %d`, selectStr, selectID)

	return selectID // 向下执行
}
