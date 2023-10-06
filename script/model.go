package script

import (
	"lucksystem/charset"
	"lucksystem/pak"
)

type GlobalLabel struct {
	ScriptName string
	Index      int // 序号，从1开始
	Position   int
}

type JumpParam struct {
	GlobalIndex int
	ScriptName  string
	Position    int
}

type StringParam struct {
	Data   string
	Coding charset.Charset
	HasLen bool
}

type Info struct {
	FileName string
	Name     string
	CodeNum  int
}

type CodeInfo struct {
	Index      int // 序号
	Pos        int // 文件偏移，和Index绑定
	LabelIndex int // 跳转目标标记，从1开始
	GotoIndex  int // 跳转标记，为0则不使用

	GlobalLabelIndex int // 全局跳转目标标记，从1开始
	GlobalGotoIndex  int // 全局跳转标记，为0则不使用
}

type LoadOptions struct {
	Filename string
	Entry    *pak.Entry

	Name string // optional
}
