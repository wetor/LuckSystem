package czimage

import (
	"fmt"
)

func compressLZW(data []byte) []uint16 {
	dictionary := make(map[string]uint16)
	for i := 0; i < 256; i++ {
		dictionary[string(byte(i))] = uint16(i)
	}
	code := uint16(len(dictionary) + 1)
	currChar := ""
	result := make([]uint16, 0)
	for _, c := range data {
		phrase := currChar + string(c)
		if _, isTrue := dictionary[phrase]; isTrue {
			currChar = phrase
		} else {
			result = append(result, dictionary[currChar])
			dictionary[phrase] = code
			currChar = string(c)
			code++
		}
	}
	if currChar != "" {
		result = append(result, dictionary[currChar])
	}
	return result
}

// decompressLZW cz1 cz3
func decompressLZW(compressed []uint16, size int) []byte {

	dictionary := make(map[uint16][]byte)
	for i := 0; i < 256; i++ {
		dictionary[uint16(i)] = []byte{byte(i)}
	}
	code := uint16(len(dictionary))
	w := dictionary[compressed[0]]
	decompressed := make([]byte, 0, size)
	for _, element := range compressed {
		var entry []byte
		if x, ok := dictionary[element]; ok {
			entry = make([]byte, len(x))
			copy(entry, x)
		} else if element == code {
			entry = make([]byte, len(w), len(w)+1)
			copy(entry, w)
			entry = append(entry, w[0])
		} else {
			panic(fmt.Sprintf("Bad compressed element: %d", element))
		}
		decompressed = append(decompressed, entry...)
		w = append(w, entry[0])
		dictionary[code] = w
		code++

		w = entry
	}
	return decompressed
}

// decompressLZW2 cz2
func decompressLZW2(compressed []uint16, size int) []byte {

	dictionary := make(map[uint16][]byte)
	for i := 0; i < 256; i++ {
		dictionary[uint16(i)] = []byte{byte(i)}
	}
	code := uint16(len(dictionary))
	w := dictionary[compressed[0]]
	decompressed := make([]byte, 0, size)
	for _, element := range compressed {
		element /= 2
		var entry []byte
		if x, ok := dictionary[element]; ok {
			entry = make([]byte, len(x))
			copy(entry, x)
		} else if element == code {
			entry = make([]byte, len(w), len(w)+1)
			copy(entry, w)
			entry = append(entry, w[0])
		} else {
			panic(fmt.Sprintf("Bad compressed element: %d", element))
		}
		decompressed = append(decompressed, entry...)
		w = append(w, entry[0])
		dictionary[code] = w
		code++

		w = entry
	}
	return decompressed
}

// decompressLZW_2 fast 有问题
func decompressLZW_2(compressed []uint16, size int) []byte {

	dictionary := make(map[uint16]string)
	for i := 0; i < 256; i++ {
		dictionary[uint16(i)] = string(byte(i))
	}
	code := uint16(len(dictionary))
	w := dictionary[compressed[0]]
	decompressed := make([]byte, 0, size)
	for _, element := range compressed {
		var entry string
		if x, ok := dictionary[element]; ok {
			entry = x
		} else if element == code {
			entry = w + string(w[0])
		} else {
			panic(fmt.Sprintf("Bad compressed element: %d", element))
		}
		decompressed = append(decompressed, entry...)
		dictionary[code] = w + string(entry[0])
		code++

		w = entry
	}
	return decompressed
}
