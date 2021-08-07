package script

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

func ToStringCodeParams(code *CodeLine) string {
	paramStr := make([]string, 0, len(code.Params))
	for i := 0; i < len(code.Params); i++ {
		switch param := code.Params[i].(type) {
		case []uint16:
			for _, val := range param {
				paramStr = append(paramStr, strconv.FormatInt(int64(val), 10))
			}
		case byte:
			paramStr = append(paramStr, fmt.Sprintf("0x%X", param))
		case string:
			paramStr = append(paramStr, `"`+param+`"`)
		case *JumpParam:
			if code.GotoIndex > 0 {
				if len(param.ScriptName) > 0 {
					// paramStr = append(paramStr, fmt.Sprintf(`{goto "%s" label%d}`, param.ScriptName, code.GotoIndex))
				} else {
					paramStr = append(paramStr, fmt.Sprintf("{goto label%d}", code.GotoIndex))
				}
			} else {
				if len(param.ScriptName) > 0 {
					// paramStr = append(paramStr, fmt.Sprintf(`{goto "%s" %d}`, param.ScriptName, param.Position))
				} else {
					paramStr = append(paramStr, fmt.Sprintf("{goto %d}", param.Position))
				}
			}
		default:
			paramStr = append(paramStr, fmt.Sprintf("%v", param))

		}
	}
	str := strings.Join(paramStr, ", ")

	if code.LabelIndex > 0 {
		return fmt.Sprintf(`label%d: %s (%s)`, code.LabelIndex, code.OpStr, str)
	} else {
		return fmt.Sprintf(`%s (%s)`, code.OpStr, str)
	}
}

func ParseCodeParams(code *CodeLine, codeStr string) {
	word := make([]rune, 0, 32)
	params := make([]interface{}, 0, 8)
	opStr := ""
	labelIndex := 0
	gotoIndex := 0
	gotoFile := ""
	isString := false
	isSpecial := false

	for _, ch := range codeStr {
		if isString {
			if ch == '"' {
				isString = false
				continue
			}
			word = append(word, ch)
			continue
		}

		switch ch {
		case ' ', ',', '(', ')', '}', '\n':
			if len(word) > 0 {
				wordStr := string(word)
				if opStr == "" {
					opStr = wordStr
				} else if isSpecial {
					if len(word) > 5 && wordStr[0:5] == "label" {
						gotoIndex, _ = strconv.Atoi(wordStr[5:])
					} else if wordStr != "goto" {
						gotoFile = wordStr
					}
				} else {
					params = append(params, wordStr)
				}
				word = word[0:0]
			}
			if ch == '}' {
				isSpecial = false
			}
		case ':':
			labelIndex, _ = strconv.Atoi(string(word[5:]))
			word = word[0:0]
		case '"':
			isString = true
		case '{':
			isSpecial = true
		default:
			word = append(word, ch)
		}
	}
	code.OpStr = opStr
	code.LabelIndex = labelIndex
	code.GotoIndex = gotoIndex

	if gotoFile != "" || gotoIndex > 0 {
		params = append(params, &JumpParam{
			ScriptName: gotoFile,
			Position:   gotoIndex,
		})
	}

	code.Params = params
	// if labelIndex > 0 {
	// 	fmt.Printf("label%d: ", labelIndex)
	// }
	// fmt.Printf("%s %v", opStr, params)
	// if gotoIndex > 0 {
	// 	fmt.Printf(" {goto label%d}", gotoIndex)
	// }
	// fmt.Print("\n")
	// _, _, _ = labelIndex, gotoIndex, params
}
func GetOperateName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	name := f.Name()
	return name[strings.LastIndex(name, ".")+1:]
}
