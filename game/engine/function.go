// function 模拟器相关
// 此包内的方法与模拟器相关，负责画面的输出、按键的输入等操作
package engine

import (
	"github.com/golang/glog"
)

func (Engine) FARCALL(params ...interface{}) int {
	if len(params) != 3 {
		panic("参数数量错误")
	}
	index := params[0].(uint16)
	fileStr := params[1].(string)
	jumpPos := params[2].(uint32)
	glog.V(3).Infof("Engine: FARCALL (%d) {goto \"%s\", %d}\n", index, fileStr, jumpPos)
	return 0 // 向下执行
}

func (Engine) JUMP(params ...interface{}) int {
	if len(params) != 2 {
		panic("参数数量错误")
	}

	fileStr := params[0].(string)
	jumpPos := params[1].(uint32)
	glog.V(3).Infof("Engine: JUMP {goto \"%s\", %d}\n", fileStr, jumpPos)
	return 0 // 向下执行
}
