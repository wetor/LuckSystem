package czimage

import (
	"bytes"
	"encoding/binary"
	"os"
	"strconv"
)

func Decompress(data []byte) []byte {
	offset := 0
	fileCount := binary.LittleEndian.Uint32(data[offset : offset+4])
	offset += 4
	rawSizeList := make(map[int]uint32)
	compressedSizeList := make(map[int]uint32)
	outputInfo := CzOutputInfo{
		FileCount:           fileCount,
		TotalRawSize:        0,
		TotalCompressedSize: 0,
		BlockInfo:           make([]CzBlockInfo, fileCount),
	}

	for i := 0; i < int(fileCount); i++ {
		fileCompressedSize := binary.LittleEndian.Uint32(data[offset : offset+4])
		offset += 4
		fileRawSize := binary.LittleEndian.Uint32(data[offset : offset+4])
		offset += 4

		rawSizeList[i] = fileRawSize
		compressedSizeList[i] = fileCompressedSize

		outputInfo.TotalRawSize += fileRawSize
		outputInfo.TotalCompressedSize += fileCompressedSize
		outputInfo.BlockInfo[i] = CzBlockInfo{
			BlockIndex:     uint32(i),
			RawSize:        fileRawSize,
			CompressedSize: fileCompressedSize,
		}
	}

	// fmt.Println("uncompress info", outputInfo)
	outputBuf := &bytes.Buffer{}
	for i := 0; i < int(fileCount); i++ {
		lzwBuf := make([]uint16, int(compressedSizeList[i]))
		//offsetTemp := offset
		for j := 0; j < int(compressedSizeList[i]); j++ {
			lzwBuf[j] = binary.LittleEndian.Uint16(data[offset : offset+2])
			offset += 2
		}
		//os.WriteFile("../data/LB_EN/IMAGE/2.ori.lzw", data[offsetTemp:offset], 0666)
		rawBuf := decompressLZW(lzwBuf, int(rawSizeList[i]))
		//os.WriteFile("../data/LB_EN/IMAGE/2.ori", rawBuf, 0666)
		//panic("11")
		outputBuf.Write(rawBuf)
	}
	return outputBuf.Bytes()

}

func Decompress2(data []byte) []byte {
	offset := 0
	fileCount := binary.LittleEndian.Uint32(data[offset : offset+4])
	offset += 4
	rawSizeList := make(map[int]uint32)
	compressedSizeList := make(map[int]uint32)
	outputInfo := CzOutputInfo{
		FileCount:           fileCount,
		TotalRawSize:        0,
		TotalCompressedSize: 0,
		BlockInfo:           make([]CzBlockInfo, fileCount),
	}

	for i := 0; i < int(fileCount); i++ {
		fileCompressedSize := binary.LittleEndian.Uint32(data[offset : offset+4])
		offset += 4
		fileRawSize := binary.LittleEndian.Uint32(data[offset : offset+4])
		offset += 4

		rawSizeList[i] = fileRawSize
		compressedSizeList[i] = fileCompressedSize

		outputInfo.TotalRawSize += fileRawSize
		outputInfo.TotalCompressedSize += fileCompressedSize
		outputInfo.BlockInfo[i] = CzBlockInfo{
			BlockIndex:     uint32(i),
			RawSize:        fileRawSize,
			CompressedSize: fileCompressedSize,
		}
	}

	// fmt.Println("uncompress info", outputInfo)
	outputBuf := &bytes.Buffer{}
	for i := 0; i < int(fileCount); i++ {
		lzwBuf := make([]uint16, int(compressedSizeList[i]/2))
		//offsetTemp := offset
		for j := 0; j < int(compressedSizeList[i]/2); j++ {
			lzwBuf[j] = binary.LittleEndian.Uint16(data[offset : offset+2])
			offset += 2
		}
		//os.WriteFile("../data/LB_EN/IMAGE/2.ori.lzw", data[offsetTemp:offset], 0666)
		rawBuf := decompressLZW2(lzwBuf, int(rawSizeList[i]))
		os.WriteFile("../data/Other/CZ2/ゴシック14.raw.1", rawBuf, 0666)
		//panic("11")
		outputBuf.Write(rawBuf)
		break
	}

	//// fmt.Println("uncompress info", outputInfo)
	//compressedSize := 0
	//rawSize := 0
	//outputBuf := &bytes.Buffer{}
	//for i := 0; i < int(fileCount); i++ {
	//	compressedSize += int(compressedSizeList[i]) / 2
	//	rawSize += int(rawSizeList[i])
	//}
	//lzwBuf := make([]uint16, compressedSize)
	////offsetTemp := offset
	//for j := 0; j < compressedSize; j++ {
	//	lzwBuf[j] = binary.LittleEndian.Uint16(data[offset : offset+2])
	//	offset += 2
	//}
	////os.WriteFile("../data/LB_EN/IMAGE/2.ori.lzw", data[offsetTemp:offset], 0666)
	//rawBuf := decompressLZW2(lzwBuf, rawSize)
	////os.WriteFile("../data/Other/CZ2/ゴシック14.raw.1", rawBuf, 0666)
	////panic("11")
	//outputBuf.Write(rawBuf)
	return outputBuf.Bytes()

}
func Compress(data []byte, maxCount int) []byte {

	for {
		code := 256
		count := 0
		dictionary := make(map[string]int)
		for i := 0; i < 256; i++ {
			dictionary[strconv.Itoa(i)] = i
		}

		currChar := ""
		result := make([]int, 0)
		for _, c := range data {
			phrase := currChar + string(c)
			if _, isTrue := dictionary[phrase]; isTrue {
				currChar = phrase
			} else {
				result = append(result, dictionary[currChar])
				dictionary[phrase] = code
				code++
				currChar = string(c)
			}
			count++
			if len(result) == maxCount {
				break
			}

		}
	}
	return nil
}
