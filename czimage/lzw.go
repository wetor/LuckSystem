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
	var rax, resultIndex, rcx, rdx, dataIndex, rdi, resultSize int
	var r8d, r9d, r10d, r11d, dataSize, r14d, dictIndex int
	resultSize = size // 解压长度
	resultIndex = 0   // 解压指针
	result := make([]byte, resultSize)

	dataSize = len(data) // 待解压长度
	dataIndex = 0        // 待解压指针
	data = append(data, []byte{0, 0}...)

	posDict := map[int]int{}

	for {
		rcx = rdi
		posDict[dictIndex] = resultIndex
		dictIndex++
		rdx = int(data[dataIndex])
		rax = 1
		rdi++
		rax <<= rcx & 0xFF // shl eax,cl
		rax &= rdx

		if rdi >= 8 {
			rdx = int(data[dataIndex+1]) // movzx edx,byte ptr [rsi+01]
			rdi = 0                      // xor edi,edi
			dataIndex++                  // inc rsi
		}
		r10d = 8
		r11d = rdx & 0xFF // movzx r11d,dl
		rcx = rdi & 0xFF  // movzx ecx,dil
		r10d -= rdi
		r11d >>= rcx & 0xFF

		dataIndex++ // inc rsi
		r9d = rdi
		if rax == 0 {
			r8d = int(data[dataIndex])
			rcx = r10d
			r8d <<= rcx & 0xFF // shl r8d,cl
			rdi += 0x0F
			r8d &= 0x7FFF
			r8d |= r11d
			if rdi > 0x10 {
				rdx = int(data[dataIndex+1]) // movzx edx,byte ptr [rsi+01]
				dataIndex++                  // inc rsi
				rcx = 0x10                   // mov ecx,00000010
				rcx -= r9d                   // sub ecx,r9d
				rdx <<= rcx & 0xFF           // shl edx,cl
				rdx &= 0x7FFF                // and edx,00007FFF
				r8d |= rdx                   // or r8d,edx
				rdi &= 7                     // and edi,07
			} else if rdi == 0x10 {
				dataIndex++
				rdi &= 7
			} else if rdi >= 0x8 { // cmp edi,08
				rdi &= 7
			}
		} else {
			rdx = int(data[dataIndex])   // movzx edx,byte ptr [rsi]
			rax = rdi + 0x12             // lea eax,[rdi+12]
			r8d = int(data[dataIndex+1]) // movzx r8d,byte ptr [rsi+01]
			dataIndex++                  // inc rsi
			rcx = 0x10                   // mov ecx,00000010
			rdi = rax                    // mov edi,eax
			rcx -= r9d                   // sub ecx,r9d
			r8d <<= rcx & 0xFF           // shl r8d,cl
			rcx = r10d                   // mov ecx,r10d
			r8d &= 0x3FFFF               // and r8d,0003FFFF
			rdx <<= rcx & 0xFF           // shl edx,cl
			r8d |= rdx                   // or r8d,edx
			r8d |= r11d                  // or r8d,r11d

			flag := false
			if rax > 0x18 { // cmp eax,18
				rdx = int(data[dataIndex+1]) // movzx edx,byte ptr [rsi+01]
				dataIndex++                  // inc rsi
				rcx = 0x18                   // mov ecx,00000018
				rcx -= r9d                   // sub ecx,r9d`
				rdx <<= rcx & 0xFF           // shl edx,cl
				rdx &= 0x3FFFF               // and edx,0003FFFF
				r8d |= rdx                   // or r8d,edx
			} else if rax == 0x18 {
				dataIndex++ // inc rsi
			} else if rax < 0x8 { // cmp eax,08
				flag = true
			}
			if !flag {
				rax = rdi                         // mov eax,edi
				rax = int(int32(rax) &^ int32(7)) // and eax,-08
				rdi -= rax                        // sub edi,eax
			}
		}
		if dataIndex > dataSize {
			break
		}
		if r8d < 0x100 { // cmp r8d,0x100
			result[resultIndex] = byte(r8d) // mov [rbx],r8l
			resultIndex++                   // inc rbx
			continue
		}
		rax = r8d // movsxd rax,r8d

		r8d = posDict[rax-256] // mov r8,[r13+rax*8-00000800]
		rdx = posDict[rax-257] // mov rdx,[r13+rax*8-00000808]
		rcx = r8d
		rcx -= rdx
		rcx++

		rax = rcx
		rax += resultIndex // add rax,rbx
		if rax >= resultSize {
			rcx = resultSize   // mov ecx,ebp
			rcx -= resultIndex // sub ecx,ebx
			r8d = rcx          // movsxd  r8,ecx
			r8d += rdx         // add r8,rdx
		}
		if rcx&1 != 0 {
			result[resultIndex] = result[rdx]
			rdx++
			resultIndex++
		}
		if rcx&2 != 0 {
			for i := 0; i < 2; i++ {
				result[resultIndex+i] = result[rdx+i]
			}
			resultIndex += 2
			rdx += 2
		}
		if rcx&4 != 0 {
			for i := 0; i < 4; i++ {
				result[resultIndex+i] = result[rdx+i]
			}
			resultIndex += 4
			rdx += 4
		}
		if rcx&8 != 0 {
			for i := 0; i < 8; i++ {
				result[resultIndex+i] = result[rdx+i]
			}
			resultIndex += 8
			rdx += 8
		}

		if rdx >= r8d {
			continue
		}
		r14d = rcx
		r14d = int(int64(r14d) &^ int64(0xF)) // and r14,-10
		rax = r14d + rdx
		if resultIndex != rax {
			r8d = r14d
			rcx = resultIndex
			for i := 0; i < r14d; i++ {
				result[resultIndex+i] = result[rdx+i]
			}
			// ret
			resultIndex += r14d
			continue
		}
		//nop dword ptr [rax+00]
		for rdx < r8d {
			for i := 0; i < 0x10; i++ {
				result[resultIndex+i] = result[rdx+i]
			}
			resultIndex += 0x10
			rdx += 0x10
			//cmp rdx,r8
			//jb LOOPERS.exe+F0350
		}
		continue //jmp LOOPERS.exe+F0150
	}
	return result
}
