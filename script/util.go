package script

import (
	"fmt"
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
			if param.LabelIndex > 0 {
				paramStr = append(paramStr, fmt.Sprintf("{goto label%d}", param.LabelIndex))
			} else if code.GotoIndex > 0 {
				paramStr = append(paramStr, fmt.Sprintf("{goto label%d}", code.GotoIndex))
			} else if param.GlobalIndex > 0 {
				paramStr = append(paramStr, fmt.Sprintf(`{goto "%s" global%d}`, param.ScriptName, param.GlobalIndex))
			} else {
				paramStr = append(paramStr, fmt.Sprintf("{goto %d}", param.Position))
			}
		default:
			paramStr = append(paramStr, fmt.Sprintf("%v", param))

		}
	}
	str := strings.Join(paramStr, ", ")
	str = fmt.Sprintf(`%s (%s)`, code.OpStr, str)
	if code.LabelIndex > 0 {
		str = fmt.Sprintf(`label%d: %s`, code.LabelIndex, str)
	}
	if code.GlobalLabelIndex > 0 {
		str = fmt.Sprintf(`global%d: %s`, code.GlobalLabelIndex, str)
	}
	return str
}

func ParseCodeParams(code *CodeLine, codeStr string) {
	word := make([]rune, 0, 32)
	params := make([]interface{}, 0, 8)
	opStr := ""
	labelIndex := 0
	gotoIndex := 0
	globalLabelIndex := 0
	globalGotoIndex := 0
	gotoFile := ""
	specialGotoIndex := 0
	specialGlobalGotoIndex := 0
	specialGotoFile := ""
	isString := false
	isSpecial := false

	for _, ch := range codeStr {
		if isString {
			if ch == '"' {
				if len(word) == 0 { // 空字符串
					word = append(word, '\x00')
				}
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
				if word[0] == '\x00' {
					wordStr = ""
				}
				if opStr == "" {
					opStr = wordStr
				} else if isSpecial {
					if len(word) > 5 && wordStr[0:5] == "label" {
						specialGotoIndex, _ = strconv.Atoi(wordStr[5:])
					} else if len(word) > 6 && wordStr[0:6] == "global" {
						specialGlobalGotoIndex, _ = strconv.Atoi(wordStr[6:])
					} else if wordStr != "goto" {
						specialGotoFile = wordStr
					}
				} else {
					params = append(params, wordStr)
				}
				word = word[0:0]
			}
			if ch == '}' {
				if specialGotoFile != "" || specialGotoIndex > 0 || specialGlobalGotoIndex > 0 {
					params = append(params, &JumpParam{
						GlobalIndex: specialGlobalGotoIndex,
						LabelIndex:  specialGotoIndex,
						ScriptName:  specialGotoFile,
						Position:    specialGotoIndex + specialGlobalGotoIndex, // 填充使用
					})
					if gotoIndex == 0 && specialGotoIndex > 0 {
						gotoIndex = specialGotoIndex
					}
					if globalGotoIndex == 0 && specialGlobalGotoIndex > 0 {
						globalGotoIndex = specialGlobalGotoIndex
					}
					if gotoFile == "" && specialGotoFile != "" {
						gotoFile = specialGotoFile
					}
				}
				specialGotoIndex = 0
				specialGlobalGotoIndex = 0
				specialGotoFile = ""
				isSpecial = false
			}
		case ':':
			if len(word) > 5 && string(word[0:5]) == "label" {
				labelIndex, _ = strconv.Atoi(string(word[5:]))
			} else if len(word) > 6 && string(word[0:6]) == "global" {
				globalLabelIndex, _ = strconv.Atoi(string(word[6:]))
			}
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

	code.GlobalGotoIndex = globalGotoIndex
	code.GlobalLabelIndex = globalLabelIndex

	if isSpecial && (specialGotoFile != "" || specialGotoIndex > 0 || specialGlobalGotoIndex > 0) {
		params = append(params, &JumpParam{
			GlobalIndex: specialGlobalGotoIndex,
			LabelIndex:  specialGotoIndex,
			ScriptName:  specialGotoFile,
			Position:    specialGotoIndex + specialGlobalGotoIndex, // 填充使用
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
