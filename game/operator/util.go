package operator

import (
	"encoding/binary"
	"fmt"

	"lucksystem/charset"
)

type lstring string // len + string

// DecodeString 从指定位置读取一个指定编码字符串
// 以"0x00 0x00"结尾
//   1.bytes 要读取的字节数据
//   2.start 开始位置
//   3.slen 不包含EOF的字符串字节长度，为0则读取到EOF
//   4.coding 源编码
// return
//   1.string 转码后（utf8）的字符串
//   2.uint16 读取完毕后的字节位置，结尾已跳过
// charset
//   1.UTF-8 1~3byte一字，EOF为0x00
//   2.ShiftJIS 1~2byte一字，EOF为0x00
//   3.Unicode(UTF-16LE) 2byte一字，EOF为0x00 0x00

func DecodeString(bytes []byte, start, slen int, coding charset.Charset) (string, int) {
	end := start
	eofLen := 0 //
	charLen := 0

	switch coding {
	case charset.ShiftJIS:
		fallthrough
	case charset.UTF_8:
		eofLen = 1
		charLen = 1
	case charset.Unicode:
		fallthrough
	default:
		eofLen = 2
		charLen = 2
	}

	if slen == 0 {
		switch coding {
		case charset.ShiftJIS:
			fallthrough
		case charset.UTF_8:
			for end < len(bytes) && !(bytes[end] == 0) {
				end += charLen
			}
		case charset.Unicode:
			fallthrough
		default:
			for end+1 < len(bytes) && !(bytes[end] == 0 && bytes[end+1] == 0) {
				end += charLen
			}
		}
	} else {
		end = start + slen
	}

	str, _ := charset.ToUTF8(coding, bytes[start:end])
	return str, end + eofLen
}

func ToUint8(data byte) uint8 {
	return uint8(data)
}

func ToUint16(data []byte) uint16 {
	return binary.LittleEndian.Uint16(data)
}
func ToUint32(data []byte) uint32 {
	return binary.LittleEndian.Uint32(data)
}

func ToString(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

// AllToUint16 将数据转为uint16列表，若正好转化完则返回0，否则返回最后一位位置
func AllToUint16(data []byte) (list []uint16, end int) {
	dataLen := len(data)
	if dataLen%2 == 0 {
		end = -1
	} else {
		end = dataLen - 1
	}
	list = make([]uint16, 0, dataLen/2)
	for i := 0; i < (dataLen & ^1); i += 2 {
		list = append(list, binary.LittleEndian.Uint16(data[i:i+2]))
	}
	return list, end
}

// SetParam 参数转为字节
//
//	1.codeBytes 完整的参数字节数据
//	2.data[0] Paramter类型指针
//	3.data[1] start 可空，默认0。当前参数开始位置
//	4.data[2] size 可空，默认对于Paramter类型长度。当前参数字节长度
//	5.data[3] coding 可空，默认Unicode。LString类型编码
//	return start+size，即下个参数的start
func GetParam(codeBytes []byte, data ...interface{}) int {
	var start, size int
	var coding charset.Charset
	if len(data) >= 2 {
		start = data[1].(int)
	} else {
		start = 0
	}

	if len(data) >= 3 {
		size = data[2].(int)
	} else {
		size = 0
	}

	if len(data) >= 4 {
		coding = data[3].(charset.Charset)
	} else {
		coding = charset.Unicode
	}

	switch value := data[0].(type) {
	case *uint8:
		*value = codeBytes[start]
		return start + 1
	case *uint16:
		if size == 0 {
			size = 2
		}
		*value = ToUint16(codeBytes[start : start+size])
		return start + size
	case *uint32:
		if size == 0 {
			size = 4
		}
		*value = ToUint32(codeBytes[start : start+size])
		return start + size
	case *string:
		tmp, next := DecodeString(codeBytes, start, size, coding)
		*value = tmp
		return next
	case *lstring:
		l := int(ToUint16(codeBytes[start : start+2]))
		switch coding {
		case charset.Unicode, charset.ShiftJIS:
			size = l * 2
		case charset.UTF_8:
			size = 0x10000 - l
		}
		if l == 0 {
			size = 0
		}
		tmp, next := DecodeString(codeBytes, start+2, size, coding)
		*value = lstring(tmp)
		return next
	default:
		return start
	}
}
