package operation

import (
	"lucascript/charset"
	"lucascript/script"
)

type Operation interface {
	UNDEFINE(code *script.CodeLine, opname string) string
	EQU(code *script.CodeLine) string
	ADD(code *script.CodeLine) string
	MESSAGE(code *script.CodeLine) string
	IFN(code *script.CodeLine) string
	GOTO(code *script.CodeLine) string
}

// ReadUnicode 从指定位置读取一个指定编码字符串
// 以"0x00 0x00"结尾
//   1.bytes 要读取的字节数据
//   2.start 开始位置
//   3.coding 源编码
// return
//   1.string 转码后（utf8）的字符串
//   2.uint16 读取完毕后的字节位置，结尾已跳过
func ReadString(bytes []byte, start int, coding charset.Charset) (string, int) {
	end := start
	eofLen := 0 //
	charLen := 0
	switch coding {
	case charset.UTF_8:
		eofLen = 1
		charLen = 1
		for end < len(bytes) && !(bytes[end] == 0) {
			end += charLen
		}
	case charset.Unicode:
		fallthrough
	default:
		eofLen = 2
		charLen = 2
		for end+1 < len(bytes) && !(bytes[end] == 0 && bytes[end+1] == 0) {
			end += charLen
		}
	}

	str, _ := charset.ToUTF8(coding, bytes[start:end])
	return str, end + eofLen
}
