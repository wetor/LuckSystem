package czimage

import (
	"fmt"
)

// compressLZW lzw压缩
//
//	Description lzw压缩一块
//	Param data []byte 未压缩数据
//	Param size int 压缩后数据的大小限制
//	Param last string 上一个lzw压缩剩余的element
//	Return count 使用数据量
//	Return compressed 压缩后的数据
//	Return lastElement lzw压缩剩余的element
func compressLZW(data []byte, size int, last string) (count int, compressed []uint16, lastElement string) {
	count = 0
	dictionary := make(map[string]uint16)
	for i := 0; i < 256; i++ {
		dictionary[string(byte(i))] = uint16(i)
	}
	dictionaryCount := uint16(len(dictionary) + 1)
	element := ""
	if len(last) != 0 {
		element = last
	}
	compressed = make([]uint16, 0, size)
	for _, c := range data {
		entry := element + string(c)
		if _, isTrue := dictionary[entry]; isTrue {
			element = entry
		} else {
			compressed = append(compressed, dictionary[element])
			dictionary[entry] = dictionaryCount
			element = string(c)
			dictionaryCount++
		}
		count++
		if size > 0 && len(compressed) == size {
			break
		}
	}
	lastElement = element
	if len(compressed) == 0 {
		if len(lastElement) != 0 {
			// 数据在上一片压缩完毕，此次压缩无数据，剩余lastElement，拆分加入
			for _, c := range lastElement {
				compressed = append(compressed, dictionary[string(c)])
			}
		}
		return count, compressed, ""
	} else if len(compressed) < size {
		// 数据压缩完毕，未达到size，则为最后一片，直接加入剩余数据
		if len(lastElement) != 0 {
			compressed = append(compressed, dictionary[lastElement])
		}
		return count, compressed, ""
	}
	// 数据压缩完毕，达到size，剩余数据交给下一片
	return count, compressed, lastElement
}

// decompressLZW lzw解压
//
//	Description lzw解压一块
//	Param compressed []uint16 压缩的数据
//	Param size int 未压缩数据大小，可超过
//	Return []byte 解压后的数据
func decompressLZW(compressed []uint16, size int) []byte {

	dictionary := make(map[uint16][]byte)
	for i := 0; i < 256; i++ {
		dictionary[uint16(i)] = []byte{byte(i)}
	}
	dictionaryCount := uint16(len(dictionary))
	w := dictionary[compressed[0]]
	decompressed := make([]byte, 0, size)
	for _, element := range compressed {
		var entry []byte
		if x, ok := dictionary[element]; ok {
			entry = make([]byte, len(x))
			copy(entry, x)
		} else if element == dictionaryCount {
			entry = make([]byte, len(w), len(w)+1)
			copy(entry, w)
			entry = append(entry, w[0])
		} else {
			panic(fmt.Sprintf("Bad compressed element: %d", element))
		}
		decompressed = append(decompressed, entry...)
		w = append(w, entry[0])
		dictionary[dictionaryCount] = w
		dictionaryCount++

		w = entry
	}
	return decompressed
}

// DecompressLZWByAsm
//
//	CZ2的LZW解压算法，暂时汇编实现
func DecompressLZWByAsm(data []byte, size int) []byte {
	var lastBit, resultIndex, dataIndex, resultSize int
	var dataSize, dictIndex int
	resultSize = size // 解压长度
	resultIndex = 0   // 解压指针
	result := make([]byte, resultSize)

	dataSize = len(data) // 待解压长度
	dataIndex = 0        // 待解压指针
	data = append(data, []byte{0, 0}...)

	posDict := map[int]int{}

	for {
		posDict[dictIndex] = resultIndex
		dictIndex++

		flag := int(data[dataIndex]) & (1 << lastBit)

		lastBit++
		if lastBit >= 8 {
			lastBit = 0
			dataIndex++
		}
		code := int(data[dataIndex]) >> lastBit
		dataIndex++

		if flag == 0 {
			code |= int(data[dataIndex]) << (8 - lastBit) & 0x7FFF // 低15位
			if lastBit > 1 {
				dataIndex++
				code |= int(data[dataIndex]) << (16 - lastBit) & 0x7FFF // 低15位
			} else if lastBit == 1 {
				dataIndex++
			}
			lastBit += 15
		} else {
			code |= int(data[dataIndex]) << (8 - lastBit)
			dataIndex++
			code |= int(data[dataIndex]) << (16 - lastBit) & 0x3FFFF // 低18位
			if lastBit > 6 {
				dataIndex++
				code |= int(data[dataIndex]) << (24 - lastBit) & 0x3FFFF // 低18位
			} else if lastBit == 6 {
				dataIndex++
			}
			lastBit += 18
		}
		lastBit &= 7

		if dataIndex > dataSize {
			break
		}
		if code < 0x100 { // cmp r8d,0x100
			result[resultIndex] = byte(code)
			resultIndex++
		} else {
			dictValueEnd := posDict[code-256]   // mov r8,[r13+rax*8-00000800]
			dictValueIndex := posDict[code-257] // mov rdx,[r13+rax*8-00000808]
			dictValueSize := dictValueEnd - dictValueIndex + 1

			resultIndex = writeResult(result, resultIndex, resultSize,
				dictValueIndex, dictValueSize, dictValueEnd)
		}

	}
	return result
}

func decompressLZW2(data []byte, size int) []byte {
	dictionary := make(map[int][]byte)
	for i := 0; i < 256; i++ {
		dictionary[i] = []byte{byte(i)}
	}
	dictionaryCount := len(dictionary)
	decompressed := make([]byte, 0, size)

	dataSize := len(data)
	data = append(data, []byte{0, 0}...)
	bitIO := NewBitIO(data)
	w := dictionary[0]

	element := 0
	for {
		var entry []byte
		flag := bitIO.ReadBit(1)
		if flag == 0 {
			element = int(bitIO.ReadBit(15))
		} else {
			element = int(bitIO.ReadBit(18))
		}
		if bitIO.ByteOffset() > dataSize {
			break
		}
		if x, ok := dictionary[element]; ok {
			entry = make([]byte, len(x))
			copy(entry, x)
		} else if element == dictionaryCount {
			entry = make([]byte, len(w), len(w)+1)
			copy(entry, w)
			entry = append(entry, w[0])
		} else {
			panic(fmt.Sprintf("Bad compressed element: %d", element))
		}
		decompressed = append(decompressed, entry...)
		w = append(w, entry[0])
		dictionary[dictionaryCount] = w
		dictionaryCount++
		w = entry
	}
	return decompressed
}

func DecompressLZWByAsm2(data []byte, size int) []byte {
	var resultIndex, resultSize int
	var dictIndex, code int
	resultSize = size // 解压长度
	resultIndex = 0   // 解压指针
	result := make([]byte, resultSize)
	dataSize := len(data)
	data = append(data, []byte{0, 0}...)
	posDict := map[int]int{}
	bitIO := NewBitIO(data)
	for {
		posDict[dictIndex] = resultIndex
		dictIndex++

		flag := bitIO.ReadBit(1)
		if flag == 0 {
			code = int(bitIO.ReadBit(15))
		} else {
			code = int(bitIO.ReadBit(18))
		}

		if bitIO.ByteOffset() > dataSize {
			break
		}
		if code < 0x100 { // cmp r8d,0x100
			result[resultIndex] = byte(code)
			resultIndex++
		} else {
			dictValueEnd := posDict[code-256]
			dictValueIndex := posDict[code-257]
			dictValueSize := dictValueEnd - dictValueIndex + 1

			resultIndex = writeResult(result, resultIndex, resultSize,
				dictValueIndex, dictValueSize, dictValueEnd)
		}

	}
	return result
}

func writeResult(result []byte, index, resultSize int,
	dictValueIndex, dictValueSize, dictValueEnd int) (resultIndex int) {
	resultIndex = index
	if dictValueSize+resultIndex >= resultSize {
		dictValueSize = resultSize - resultIndex
		dictValueEnd = dictValueSize + dictValueIndex
	}
	if dictValueSize&1 != 0 {
		result[resultIndex] = result[dictValueIndex]
		dictValueIndex++
		resultIndex++
	}
	if dictValueSize&2 != 0 {
		for i := 0; i < 2; i++ {
			result[resultIndex+i] = result[dictValueIndex+i]
		}
		resultIndex += 2
		dictValueIndex += 2
	}
	if dictValueSize&4 != 0 {
		for i := 0; i < 4; i++ {
			result[resultIndex+i] = result[dictValueIndex+i]
		}
		resultIndex += 4
		dictValueIndex += 4
	}
	if dictValueSize&8 != 0 {
		for i := 0; i < 8; i++ {
			result[resultIndex+i] = result[dictValueIndex+i]
		}
		resultIndex += 8
		dictValueIndex += 8
	}

	if dictValueIndex < dictValueEnd {
		copySize := int(int64(dictValueSize) &^ int64(0xF)) // and r14,-10
		if resultIndex != copySize+dictValueIndex {
			dictValueSize = resultIndex
			for i := 0; i < copySize; i++ {
				result[resultIndex+i] = result[dictValueIndex+i]
			}
			resultIndex += copySize
		} else {
			for dictValueIndex < dictValueEnd {
				for i := 0; i < 0x10; i++ {
					result[resultIndex+i] = result[dictValueIndex+i]
				}
				resultIndex += 0x10
				dictValueIndex += 0x10
			}
		}
	}
	return resultIndex
}
