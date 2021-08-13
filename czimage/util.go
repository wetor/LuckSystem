package czimage

import (
	"bytes"
	"encoding/binary"
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
		lzwBuf := make([]int, int(compressedSizeList[i]))
		for j := 0; j < int(compressedSizeList[i]); j++ {
			lzwBuf[j] = int(binary.LittleEndian.Uint16(data[offset : offset+2]))
			offset += 2
		}
		rawBuf := decompressLZW(lzwBuf, int(rawSizeList[i]))
		outputBuf.Write(rawBuf)
	}
	return outputBuf.Bytes()

}
