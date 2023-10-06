package runtime

import (
	"strings"

	"github.com/golang/glog"
	"lucksystem/script"
)

type GlobalGoto struct {
	ScriptNames map[string]struct{}

	// 当前全局标签序号
	CLabelIndexNext int
	// 标签序号 -> 目标文件,目标位置
	GlobalLabelGoto map[int]*script.GlobalLabel
	// 标签序号 -> 目标位置。导入用
	IGlobalLabelMap map[int]int

	// 目标文件 -> 目标位置 -> 标签序号
	GlobalLabelMap map[string]map[int]int
	// 目标文件 -> 代码序号 -> 标签序号
	GlobalGotoMap map[string]map[int]int
}

func NewGlobalGoto() *GlobalGoto {
	return &GlobalGoto{
		ScriptNames:     make(map[string]struct{}),
		GlobalLabelGoto: make(map[int]*script.GlobalLabel),
		IGlobalLabelMap: make(map[int]int),
		CLabelIndexNext: 1,
		GlobalLabelMap:  make(map[string]map[int]int),
		GlobalGotoMap:   make(map[string]map[int]int),
	}
}

func (g *GlobalGoto) AddLabel(scriptName string, codeIndex, position int) (index int) {
	if _, ok := g.ScriptNames[scriptName]; !ok {
		scriptName = strings.ToUpper(scriptName)
		if _, ok := g.ScriptNames[scriptName]; !ok {
			scriptName = strings.ToLower(scriptName)
		}
	}

	index = g.CLabelIndexNext
	if _, ok := g.GlobalLabelMap[scriptName]; !ok {
		g.GlobalLabelMap[scriptName] = make(map[int]int)
	}
	if _, ok := g.GlobalGotoMap[scriptName]; !ok {
		g.GlobalGotoMap[scriptName] = make(map[int]int)
	}

	if i, ok := g.GlobalLabelMap[scriptName][position]; ok {
		return i
	}

	g.GlobalLabelMap[scriptName][position] = index
	g.GlobalGotoMap[scriptName][codeIndex] = index

	g.GlobalLabelGoto[index] = &script.GlobalLabel{
		ScriptName: scriptName,
		Index:      index,
		Position:   position,
	}
	g.CLabelIndexNext++
	return index
}

func (g *GlobalGoto) GetMaps(scriptName string) (labels map[int]int, gotos map[int]int) {
	return g.GlobalLabelMap[scriptName], g.GlobalGotoMap[scriptName]
}

func (g *GlobalGoto) AddGlobalLabelMap(labels map[int]int) {
	for index, pos := range labels {
		if _, ok := g.GlobalLabelGoto[index]; !ok {
			glog.Warningf("全局标签不存在：global%s", index)
			continue
		}
		if _, ok := g.IGlobalLabelMap[index]; ok {
			glog.Warningf("重复标签定义：global%s", index)
			continue
		}
		g.IGlobalLabelMap[index] = pos
	}
}
