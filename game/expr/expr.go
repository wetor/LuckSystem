package expr

import (
	"container/list"
	"errors"
	"strconv"
)

const (
	TOperator = 0
	TNumber   = 1
	TVariable = 2
)

type Token struct {
	Data string
	Type int
}

func RunExpr(exprStr string, variable map[string]int) (bool, error) {
	tokens, err := Parser(exprStr)
	if err != nil {
		return false, err
	}
	result, err := Exec(tokens, variable)
	if err != nil {
		return false, err
	}
	if result != 0 {
		return true, nil
	}
	return false, nil
}

func Exec(tokens []Token, variable map[string]int) (int, error) {
	stack := list.New()
	for _, token := range tokens {
		if token.Type == TVariable {
			val, has := variable[token.Data]
			if !has {
				// 变量不存在则添加为0
				variable[token.Data] = 0
				//return 0, errors.New(token.Data + " 变量不存在")
			}
			stack.PushBack(val)
		} else if token.Type == TNumber {
			val, err := strconv.Atoi(token.Data)
			if err != nil {
				return 0, err
			}
			stack.PushBack(val)
		} else if token.Type == TOperator {
			if stack.Len() < 2 {
				return 0, errors.New("表达式错误")
			}
			B := stack.Back()
			stack.Remove(B)
			A := stack.Back()
			stack.Remove(A)
			valA := A.Value.(int)
			valB := B.Value.(int)

			result := Calc(valA, valB, token.Data)
			stack.PushBack(result)
		}
	}
	result := stack.Back()
	return result.Value.(int), nil
}

// Parser 将字符串表达式转换为逆序表达式
func Parser(exprStr string) (tokens []Token, err error) {
	tokens = make([]Token, 0, len(exprStr)/2)
	if exprStr[0] != '(' {
		exprStr = "(" + exprStr + ")"
	}
	stack := list.New()
	word := make([]byte, 0, 10)
	isNum := false
	for i := 0; i < len(exprStr); i++ {
		ch := exprStr[i]
		if ch == ' ' || IsOperator(ch) { // 单词读取结束（换行或读到操作符）
			sword := string(word)
			if len(word) > 0 {
				if isNum {
					tokens = append(tokens, Token{
						Data: sword,
						Type: TNumber,
					})
				} else {
					tokens = append(tokens, Token{
						Data: sword,
						Type: TVariable,
					})
				}
				isNum = false
				word = word[0:0]
			}
			if ch == ' ' {
				continue
			}
			word = append(word, ch)
			if i+1 < len(exprStr) && IsOperator2(ch, exprStr[i+1]) {
				word = append(word, exprStr[i+1])
				i++
			}
			sword = string(word)
			// 判断优先级
			if sword == ")" {
				top := stack.Back()
				for top != nil && top.Value.(string) != "(" {
					stack.Remove(top)
					tokens = append(tokens, Token{
						Data: top.Value.(string),
						Type: TOperator,
					})
					top = stack.Back()
				}
				stack.Remove(top)
			} else if sword == "(" {
				stack.PushBack(sword)
			} else {
				top := stack.Back()
				for top != nil && stack.Len() > 0 && (GetOperatorLevel(sword) <= GetOperatorLevel(top.Value.(string))) {
					stack.Remove(top)
					tokens = append(tokens, Token{
						Data: top.Value.(string),
						Type: TOperator,
					})
					top = stack.Back()
				}
				stack.PushBack(sword)
			}
			word = word[0:0]
			continue
		}
		if isNum {
			if (ch < '0' || ch > '9') && ch != '.' {
				return nil, errors.New("数字中包含非法字符")
			}
		}
		if len(word) == 0 { // 首个单词
			if ch >= '0' && ch <= '9' { //单词首字符为数字，则此单词为number型
				isNum = true
			} else {
				isNum = false
			}
		}
		word = append(word, ch)
	}

	top := stack.Back()
	for top != nil && stack.Len() > 0 {
		stack.Remove(top)

		tokens = append(tokens, Token{
			Data: top.Value.(string),
			Type: TOperator,
		})
		top = stack.Back()
	}
	return tokens, nil
}
