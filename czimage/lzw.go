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

func decompressLZW2(data []byte, size int) []byte {
	dictionary := make(map[int][]byte)
	for i := 0; i < 256; i++ {
		dictionary[i] = []byte{byte(i)}
	}
	dictionaryCount := len(dictionary)
	result := make([]byte, 0, size)

	dataSize := len(data)
	data = append(data, []byte{0, 0}...)
	bitIO := NewBitIO(data)
	w := dictionary[0]

	element := 0
	for {
		flag := bitIO.ReadBit(1)
		if flag == 0 {
			element = int(bitIO.ReadBit(15))
		} else {
			element = int(bitIO.ReadBit(18))
		}
		if bitIO.ByteOffset() > dataSize {
			break
		}
		var entry []byte
		if x, ok := dictionary[element]; ok {
			entry = make([]byte, len(x))
			copy(entry, x)
		} else if element == dictionaryCount {
			entry = append(w, w[0])
		} else {
			panic(fmt.Sprintf("Bad compressed element: %d", element))
		}
		result = append(result, entry...)
		w = append(w, entry[0])
		dictionary[dictionaryCount] = w
		dictionaryCount++
		w = entry
	}
	return result
}

func compressLZW2(data []byte, size int, last string) (count int, compressed []byte, lastElement string) {
	count = 0
	dictionary := make(map[string]uint64)
	for i := 0; i < 256; i++ {
		dictionary[string(byte(i))] = uint64(i)
	}
	dictionaryCount := uint64(len(dictionary) + 1)
	element := ""
	if len(last) != 0 {
		element = last
	}
	bitIO := NewBitIO(make([]byte, size+2))
	writeBit := func(code uint64) {
		if code > 0x7FFF {
			bitIO.WriteBit(1, 1)
			bitIO.WriteBit(code, 18)
		} else {
			bitIO.WriteBit(0, 1)
			bitIO.WriteBit(code, 15)
		}
	}

	for i, c := range data {
		_ = i
		entry := element + string(c)
		if _, isTrue := dictionary[entry]; isTrue {
			element = entry
		} else {
			writeBit(dictionary[element])
			dictionary[entry] = dictionaryCount
			element = string(c)
			dictionaryCount++
		}
		count++
		if size > 0 && bitIO.ByteSize() >= size {
			break
		}
	}
	lastElement = element
	if bitIO.ByteSize() == 0 {
		if len(lastElement) != 0 {
			// 数据在上一片压缩完毕，此次压缩无数据，剩余lastElement，拆分加入
			for _, c := range lastElement {
				writeBit(dictionary[string(c)])
			}
		}
		return count, bitIO.Bytes(), ""
	} else if bitIO.ByteSize() < size {
		// 数据压缩完毕，未达到size，则为最后一片，直接加入剩余数据
		if len(lastElement) != 0 {
			writeBit(dictionary[lastElement])
		}
		return count, bitIO.Bytes(), ""
	}
	// 数据压缩完毕，达到size，剩余数据交给下一片
	return count, bitIO.Bytes(), lastElement
}
