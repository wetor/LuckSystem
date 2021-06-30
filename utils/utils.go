package utils

import "fmt"

type DebugMode int

const (
	DebugNone DebugMode = 1 // 无额外输出
	DebugTest DebugMode = 2 // 输出info
	DebugAll  DebugMode = 3 // 输出所有信息
)

var Debug DebugMode = DebugTest

func Log(a ...interface{}) {
	if Debug >= DebugNone {
		fmt.Println(a...)
	}

}
func Logf(format string, a ...interface{}) {
	if Debug >= DebugNone {
		fmt.Printf(format+"\n", a...)
	}
}

func LogT(a ...interface{}) {
	if Debug >= DebugTest {
		fmt.Println(a...)
	}

}
func LogTf(format string, a ...interface{}) {
	if Debug >= DebugTest {
		fmt.Printf(format+"\n", a...)
	}
}

func LogA(a ...interface{}) {
	if Debug >= DebugAll {
		fmt.Println(a...)
	}

}
func LogAf(format string, a ...interface{}) {
	if Debug >= DebugAll {
		fmt.Printf(format+"\n", a...)
	}
}
