package czimage

import (
	"encoding/binary"
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

func decompressLZW22(compressed []uint16, size int) []byte {
	dictionary := make(map[uint16][]byte)
	for i := 0; i < 256; i++ {
		dictionary[uint16(i)] = []byte{byte(i)}
	}
	dictionaryCount := uint16(len(dictionary))
	w := dictionary[compressed[0]>>1]

	decompressed := make([]byte, 0, size)
	for i, element := range compressed {
		_ = i
		element = element >> 1
		var entry []byte
		if x, ok := dictionary[element]; ok {
			entry = make([]byte, len(x))
			copy(entry, x)
			// entry = dictionary[element]
		} else if element == dictionaryCount {
			entry = make([]byte, len(w), len(w)+1)
			copy(entry, w)
			entry = append(entry, w[0])
			// entry = dictionary[element-1] + dictionary[element-1][0]
		} else {
			panic(fmt.Sprintf("Bad compressed element: %d", element))
		}
		if len(decompressed) >= size {
			return decompressed
		}
		decompressed = append(decompressed, entry...)
		w = append(w, entry[0])
		dictionary[dictionaryCount] = w
		dictionaryCount++
		// dictionary[element] = dictionary[element-1] + dictionary[element-1][0]

		fmt.Printf("i: %d, e: %d, dict_count: %d, w: %v, entry: %v\n", i, element, dictionaryCount, w, entry)
		if i > 500 {
			panic("1")
		}
		w = entry
	}
	return decompressed
}

func decompressLZW222(compressed []uint16, size int) []byte {
	dictionary := make(map[uint16][]byte)
	for i := 0; i < 256; i++ {
		dictionary[uint16(i)] = []byte{byte(i)}
	}
	dictionaryCount := uint16(len(dictionary))
	w := dictionary[compressed[0]>>1]
	decompressed := make([]byte, 0, size)
	for i, element := range compressed {
		_ = i
		dictionaryCount &= 32767
		//dictionaryCount |= uint16(i)
		element = element >> 1
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

		//if len(decompressed) >= size {
		//	return decompressed
		//}
		decompressed = append(decompressed, entry...)
		w = append(w, entry[0])
		if _, ok := dictionary[dictionaryCount]; !ok {
			dictionary[dictionaryCount] = w
			dictionaryCount++
		}

		w = entry
	}
	return decompressed
}

func DecompressLZWByAsm(ptr []byte, size int) []byte {
	dict := map[int][]byte{}
	for i := 0; i < 256; i++ {
		dict[i] = []byte{byte(i)}
	}
	dictCount := 256
	var r8d_data, w []byte

	var rax, rbx, rcx, rdx, rsi, rdi, rbp int
	var r8d, r9d, r10d, r11d, r12d, r14d int
	rbp = size // 解压长度
	result := make([]byte, rbp)
	r12d = len(ptr)
	rsi = 0
	ptr = append(ptr, []byte{0, 0, 0, 0}...)

F0150:
	rcx = rdi

	//r8_map[r15d] = rbx - 1
	//rdx_map[r15d] = rbx - r14d
	//r15d++
	rdx = int(ptr[rsi])
	rax = 1
	rdi++
	rax <<= rcx & 0xFF // shl eax,cl
	rax &= rdx

	if rdi < 8 {
		goto F0175
	}
	rdx = int(ptr[rsi+1]) // movzx edx,byte ptr [rsi+01]
	rdi = 0               // xor edi,edi
	rsi++                 // inc rsi
F0175:
	r10d = 8
	r11d = rdx & 0xFF // movzx r11d,dl
	rcx = rdi & 0xFF  // movzx ecx,dil
	r10d -= rdi
	r11d >>= rcx & 0xFF

	rsi++ // inc rsi
	r9d = rdi
	if rax != 0 {
		goto F01E2
	}
	r8d = int(ptr[rsi])
	rcx = r10d
	r8d <<= rcx & 0xFF // shl r8d,cl
	rdi += 0x0F
	r8d &= 0x7FFF
	r8d |= r11d
	if rdi <= 0x10 {
		goto F01CE
	}

	rdx = int(ptr[rsi+1]) // movzx edx,byte ptr [rsi+01]
	rsi++                 // inc rsi
	rcx = 0x10            // mov ecx,00000010
	rcx -= r9d            // sub ecx,r9d
	rdx <<= rcx & 0xFF    // shl edx,cl
	rdx &= 0x7FFF         // and edx,00007FFF
	r8d |= rdx            // or r8d,edx
	rdi &= 7              // and edi,07
	goto F0243            // jmp LOOPERS.exe+F0243
F01CE:
	if rdi != 0x10 { // jle F01CE
		goto F01D8
	}
	rsi++
	rdi &= 7
	goto F0243

F01D8:
	// cmp edi,08
	if rdi < 0x8 {
		goto F0243
	}
	rdi &= 7
	goto F0243
F01E2:
	rdx = int(ptr[rsi])   // movzx edx,byte ptr [rsi]
	rax = rdi + 0x12      // lea eax,[rdi+12]
	r8d = int(ptr[rsi+1]) // movzx r8d,byte ptr [rsi+01]
	rsi++                 // inc rsi
	rcx = 0x10            // mov ecx,00000010
	rdi = rax             // mov edi,eax
	rcx -= r9d            // sub ecx,r9d
	r8d <<= rcx & 0xFF    // shl r8d,cl
	rcx = r10d            // mov ecx,r10d
	r8d &= 0x3FFFF        // and r8d,0003FFFF
	rdx <<= rcx & 0xFF    // shl edx,cl
	r8d |= rdx            // or r8d,edx
	r8d |= r11d           // or r8d,r11d
	if rax <= 0x18 {      // cmp eax,18
		goto F0230 // jle LOOPERS.exe+F0230
	}

	rdx = int(ptr[rsi+1]) // movzx edx,byte ptr [rsi+01]
	rsi++                 // inc rsi
	rcx = 0x18            // mov ecx,00000018
	rcx -= r9d            // sub ecx,r9d
	rdx <<= rcx & 0xFF    // shl edx,cl
	rdx &= 0x3FFFF        // and edx,0003FFFF
	r8d |= rdx            // or r8d,edx
	goto F023C            // jmp LOOPERS.exe+F023C
F0230:
	if rax != 0x18 {
		goto F0237 // jne LOOPERS.exe+F0237
	}
	rsi++      // inc rsi
	goto F023C // jmp LOOPERS.exe+F023C
F0237:
	if rax < 0x8 { // cmp eax,08
		goto F0243 // jl LOOPERS.exe+F0243
	}
F023C:
	rax = rdi                         // mov eax,edi
	rax = int(int32(rax) &^ int32(7)) // and eax,-08
	//rax &= 248
	rdi -= rax // sub edi,eax
F0243:
	if rsi > r12d {
		goto END
	}
	if r8d >= 0x100 { // cmp r8d,0x100
		goto F0260
	}

	result[rbx] = byte(r8d) // mov [rbx],r8l
	rbx++                   // inc rbx
	r8d_data = []byte{byte(r8d)}

	w = append(w, byte(r8d))
	dict[dictCount] = w
	dictCount++
	w = []byte{byte(r8d)}
	//fmt.Printf("i: %v, e: %d, dict_count: %d, w: %v, entry: %v\n", ok, rax, dictCount, dict[dictCount-1], r8d_data)

	goto F0150
F0260:
	rax = r8d // movsxd rax,r8d

	if x, ok := dict[rax]; ok {
		r8d_data = make([]byte, len(x))
		copy(r8d_data, x)
		// entry = dictionary[element]
	} else if rax == dictCount {
		r8d_data = make([]byte, len(w)+1)
		copy(r8d_data, append(w, w[0]))
		// entry = dictionary[element-1] + dictionary[element-1][0]
	} else {
		panic(fmt.Sprintf("Bad compressed element: %d", rax))
	}
	//fmt.Printf("i: %v, e: %d, dict_count: %d, w: %v, entry: %v\n", ok, rax, dictCount, dict[dictCount-1], r8d_data)
	w = append(w, r8d_data[0])
	dict[dictCount] = w
	dictCount++
	w = make([]byte, len(r8d_data))
	copy(w, r8d_data)

	//if rsi >= 543220 {
	//	fmt.Println(1)
	//}

	r8d = len(r8d_data) - 1
	rdx = 0

	//r8d = r8_map[rax]  // mov r8,[r13+rax*8-00000800]
	//rdx = rdx_map[rax] // mov rdx,[r13+rax*8-00000808]
	rcx = r8d
	rcx -= rdx
	rcx++
	rax = rcx
	rax += rbx // add rax,rbx
	if rax < rbp {
		goto F028F
	}
	rcx = rbp  // mov ecx,ebp
	rcx -= rbx // sub ecx,ebx
	r8d = rcx  // movsxd  r8,ecx
	r8d += rdx // add r8,rdx
F028F:
	if rcx&1 == 0 {
		goto F029F
	}
	result[rbx] = r8d_data[rdx]
	rdx++
	rbx++
F029F:
	if rcx&2 == 0 {
		goto F02B8
	}
	for i := 0; i < 2; i++ {
		result[rbx+i] = r8d_data[rdx+i]
	}
	rbx += 2
	rdx += 2
F02B8:
	if rcx&4 == 0 {
		goto F02DF
	}
	for i := 0; i < 4; i++ {
		result[rbx+i] = r8d_data[rdx+i]
	}
	rbx += 4
	rdx += 4
F02DF:
	if rcx&8 == 0 {
		goto F0322
	}
	for i := 0; i < 8; i++ {
		result[rbx+i] = r8d_data[rdx+i]
	}
	rbx += 8
	rdx += 8
F0322:
	if rdx >= r8d {
		goto F0150
	}
	r14d = rcx
	r14d = int(int64(r14d) &^ int64(0xF)) // and r14,-10
	//if r14d > 255 {
	//	r14d = 255
	//}
	//r14d &= 240
	rax = r14d + rdx
	if rbx > rax {
		goto F03D0
	}
	rax = r14d + rbx
	if rax < rdx {
		goto F03D0
	}
	//nop dword ptr [rax+00]
F0350:
	for i := 0; i < 0x10; i++ {
		result[rbx+i] = r8d_data[rdx+i]
	}
	rbx += 0x10
	rdx += 0x10
	if rdx < r8d { //cmp rdx,r8
		goto F0350 //jb LOOPERS.exe+F0350
	}
	goto F0150 //jmp LOOPERS.exe+F0150
F03D0:
	r8d = r14d
	rcx = rbx
	// TODO: call 3A2700
	for i := rdx; i < r14d; i++ {
		result[rbx+i] = r8d_data[i]
	}
	// ret
	rbx += r14d
	goto F0150

END:
	return result
}

func uint16SliceToByteSlice(input []uint16) []byte {
	var output []byte
	for _, u := range input {
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, u)
		output = append(output, b...)
	}
	return output
}
