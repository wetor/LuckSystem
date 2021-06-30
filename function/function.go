package function

import (
	"lucascript/paramter"
	"lucascript/utils"
)

type Function interface {
	String() string
	Call([]paramter.Paramter) []paramter.Paramter
}

type IfNot struct {
	Name string
}

func (f *IfNot) String() string {
	return f.Name
}
func (f *IfNot) Call(params []paramter.Paramter) int {
	if len(params) != 2 {
		return 0
	}
	ifExprStr := params[0].Value().(string)
	jumpPos := params[1].Value().(uint32)
	utils.Logf("run %s %s{goto %d}", f.Name, ifExprStr, jumpPos)
	if true {
		return int(jumpPos)
	}
	return 0 // 向下执行
}
