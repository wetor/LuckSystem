package variable

import "lucascript/game/expr"

// VariableStore 储存运行时变量的结构
type VariableStore struct {
	ValueMap map[string]int
}

func (v *VariableStore) Init() {
	v.ValueMap = make(map[string]int)
}

func (v *VariableStore) TestExpr(exprStr string) (bool, error) {
	return expr.RunExpr(exprStr, v.ValueMap)
}

func (v *VariableStore) Set(key string, value int) (create bool) {
	_, has := v.ValueMap[key]
	v.ValueMap[key] = value
	return !has
}

func (v *VariableStore) Get(key string) (int, bool) {
	value, has := v.ValueMap[key]
	return value, has
}
