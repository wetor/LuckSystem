package czimage

import (
	"bytes"
	"encoding/binary"
	"github.com/go-restruct/restruct"
	"io"
	"os"
)

//func GetOutputInfo(data []byte) (outputInfo *CzOutputInfo) {
//	fileCount := int(binary.LittleEndian.Uint32(data[0:4]))
//	outputInfo = &CzOutputInfo{
//		Offset:              4,
//		FileCount:           fileCount,
//		TotalRawSize:        0,
//		TotalCompressedSize: 0,
//		BlockInfo:           make([]CzBlockInfo, fileCount),
//	}
//
//	for i := 0; i < int(fileCount); i++ {
//		fileCompressedSize := int(binary.LittleEndian.Uint32(data[outputInfo.Offset : outputInfo.Offset+4]))
//		outputInfo.Offset += 4
//		fileRawSize := int(binary.LittleEndian.Uint32(data[outputInfo.Offset : outputInfo.Offset+4]))
//		outputInfo.Offset += 4
//
//		outputInfo.TotalRawSize += fileRawSize
//		outputInfo.TotalCompressedSize += fileCompressedSize
//		outputInfo.BlockInfo[i] = CzBlockInfo{
//			RawSize:        fileRawSize,
//			CompressedSize: fileCompressedSize,
//		}
//	}
//	return outputInfo
//}

func GetOutputInfo(data []byte) (outputInfo *CzOutputInfo) {
	outputInfo = &CzOutputInfo{}
	err := restruct.Unpack(data, binary.LittleEndian, outputInfo)
	if err != nil {
		panic(err)
	}
	for _, block := range outputInfo.BlockInfo {
		outputInfo.TotalRawSize += int(block.RawSize)
		outputInfo.TotalCompressedSize += int(block.CompressedSize)
	}
	outputInfo.Offset = 4 + int(outputInfo.FileCount)*8
	return outputInfo
}

func WriteStruct(writer io.Writer, list ...interface{}) error {
	for _, v := range list {
		temp, err := restruct.Pack(binary.LittleEndian, v)
		if err != nil {
			return err
		}
		writer.Write(temp)
	}
	return nil
}

func Decompress(data []byte, outputInfo *CzOutputInfo) []byte {
	offset := 0

	// fmt.Println("uncompress info", outputInfo)
	outputBuf := &bytes.Buffer{}
	for _, block := range outputInfo.BlockInfo {
		lzwBuf := make([]uint16, int(block.CompressedSize))
		//offsetTemp := offset
		for j := 0; j < int(block.CompressedSize); j++ {
			lzwBuf[j] = binary.LittleEndian.Uint16(data[offset : offset+2])
			offset += 2
		}
		//os.WriteFile("../data/LB_EN/IMAGE/2.ori.lzw", data[offsetTemp:offset], 0666)
		rawBuf := decompressLZW(lzwBuf, int(block.RawSize))
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
	rawSizeList := make(map[int]int)
	compressedSizeList := make(map[int]int)
	outputInfo := CzOutputInfo{
		FileCount:           fileCount,
		TotalRawSize:        0,
		TotalCompressedSize: 0,
		BlockInfo:           make([]CzBlockInfo, fileCount),
	}

	for i := 0; i < int(fileCount); i++ {
		fileCompressedSize := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
		offset += 4
		fileRawSize := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
		offset += 4

		rawSizeList[i] = fileRawSize
		compressedSizeList[i] = fileCompressedSize

		outputInfo.TotalRawSize += fileRawSize
		outputInfo.TotalCompressedSize += fileCompressedSize
		outputInfo.BlockInfo[i] = CzBlockInfo{
			RawSize:        uint32(fileRawSize),
			CompressedSize: uint32(fileCompressedSize),
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
func Compress(data []byte, size int) (compressed []byte, outputInfo *CzOutputInfo) {

	if size == 0 {
		size = 0xFEFD
	}
	//if len(outputInfo.BlockInfo) != 0 {
	//	blockSize = int(outputInfo.BlockInfo[0].CompressedSize)
	//}
	var partData []uint16
	offset := 0
	count := 0
	last := ""
	tmp := make([]byte, 2)
	outputBuf := &bytes.Buffer{}
	outputInfo = &CzOutputInfo{
		Offset:              0,
		FileCount:           0,
		TotalRawSize:        len(data),
		TotalCompressedSize: 0,
		BlockInfo:           make([]CzBlockInfo, 0),
	}
	for {
		count, partData, last = compressLZW(data[offset:], size, last)
		if count == 0 {
			break
		}
		offset += count
		for _, d := range partData {
			binary.LittleEndian.PutUint16(tmp, d)
			outputBuf.Write(tmp)
		}

		outputInfo.BlockInfo = append(outputInfo.BlockInfo, CzBlockInfo{
			CompressedSize: uint32(len(partData)),
			RawSize:        uint32(count),
		})
		outputInfo.FileCount++
	}
	outputInfo.TotalCompressedSize = outputBuf.Len() / 2

	return outputBuf.Bytes(), outputInfo
}
