package expr

/*
判断单字符是否为操作符
*/
func IsOperator(ch byte) bool {
	op := "+-*/%()[]{}><=|&^!"
	for i := 0; i < len(op); i++ {
		if ch == op[i] {
			return true
		}
	}
	return false
}

/*
判断两个字符是否为组合操作符
*/
func IsOperator2(ch, ch2 byte) bool {
	if ch == '<' && ch2 == '<' {
		return true
	} else if ch == '>' && ch2 == '>' {
		return true
	} else if ch == '=' && ch2 == '=' {
		return true
	} else if ch == '!' && ch2 == '=' {
		return true
	} else if ch == '<' && ch2 == '=' {
		return true
	} else if ch == '>' && ch2 == '=' {
		return true
	} else if ch == '|' && ch2 == '|' {
		return true
	} else if ch == '&' && ch2 == '&' {
		return true
	}
	return false
}

func GetOperatorLevel(word string) int {

	level := [][]string{
		{"*", "/", "%", ""},
		{"+", "-", "", ""},
		{"<<", ">>", "", ""},
		{">", ">=", "<", "<="},
		{"==", "!=", "", ""},
		{"&", "", "", ""},
		{"^", "", "", ""},
		{"|", "", "", ""},
		{"&&", "", "", ""},
		{"||", "", "", ""},
		{"=", "", "", ""},
	}
	for i := 0; i < 11; i++ {
		for j := 0; j < 4; j++ {
			if len(level[i][j]) == 0 {
				break
			}
			if word == level[i][j] {
				return 10 - i
			}
		}
	}
	return -1

}

func Calc(A, B int, op string) int {
	switch op {
	case "+":
		return A + B
	case "-":
		return A - B
	case "*":
		return A * B
	case "/":
		return A / B
	case "%":
		return A % B
	case "&":
		return A & B
	case "^":
		return A ^ B
	case "|":
		return A | B
	case ">>":
		return A >> B
	case "<<":
		return A << B
	case "&&":
		if A != 0 && B != 0 {
			return 1
		}
	case "||":
		if A != 0 || B != 0 {
			return 1
		}
	case ">":
		if A > B {
			return 1
		}
	case "<":
		if A < B {
			return 1
		}
	case ">=":
		if A >= B {
			return 1
		}
	case "<=":
		if A <= B {
			return 1
		}
	case "==":
		if A == B {
			return 1
		}
	case "!=":
		if A != B {
			return 1
		}
	default:
		return 0
	}
	return 0
}
