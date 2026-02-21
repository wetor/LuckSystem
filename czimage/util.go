package czimage

import (
	"bytes"
	"encoding/binary"
	"github.com/golang/glog"
	"image"
	"image/draw"
	"io"

	"github.com/go-restruct/restruct"
)

func FillImage(src image.Image, width, height int) (dst *image.NRGBA) {
	dst = image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.Draw(dst, dst.Bounds().Add(image.Pt(0, 0)), src, src.Bounds().Min, draw.Src)
	return dst
}

// GetOutputInfo 读取分块信息
//
//	Description 读取分块信息
//	Param data []byte
//	Return outputInfo
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

// WriteStruct 写入结构体
//
//	Description 写入结构体
//	Param writer io.Writer
//	Param list ...interface{}
//	Return error
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

// Decompress 解压数据
//
//	Description
//	Param data []byte 压缩的数据
//	Param outputInfo *CzOutputInfo 分块信息
//	Return []byte
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
		//panic("11")
		outputBuf.Write(rawBuf)
	}
	//os.WriteFile("../data/LB_EN/IMAGE/32.ori", outputBuf.Bytes(), 0666)
	return outputBuf.Bytes()

}

// Decompress2 解压数据 CZ2专用
//
//	Description
//	Param data []byte 压缩的数据
//	Param outputInfo *CzOutputInfo 分块信息
//	Return []byte
func Decompress2(data []byte, outputInfo *CzOutputInfo) []byte {
	offset := 0

	outputBuf := &bytes.Buffer{}
	for _, block := range outputInfo.BlockInfo {
		offsetTemp := offset
		offset += int(block.CompressedSize)
		//_ = os.WriteFile(fmt.Sprintf("C:\\Users\\wetor\\Desktop\\Prototype\\CZ2\\32\\%d_asm.src.lzw", i),
		//	data[offsetTemp:offset], 0666)
		rawBuf := decompressLZW2(data[offsetTemp:offset], int(block.RawSize))
		//_ = os.WriteFile(fmt.Sprintf("C:\\Users\\wetor\\Desktop\\Prototype\\CZ2\\32\\%d_asm.src.out", i),
		//	rawBuf, 0666)
		outputBuf.Write(rawBuf)
	}
	return outputBuf.Bytes()

}

// Compress 压缩数据
//
//	Description 压缩数据
//	Param data []byte 未压缩数据
//	Param size int 分块大小
//	Return compressed
//	Return outputInfo
func Compress(data []byte, size int) (compressed []byte, outputInfo *CzOutputInfo) {

	if size == 0 {
		size = 0xFEFD
	}
	var partData []uint16
	offset := 0
	count := 0
	last := ""
	prevCarry := 0
	tmp := make([]byte, 2)
	outputBuf := &bytes.Buffer{}
	outputInfo = &CzOutputInfo{
		TotalRawSize: len(data),
		BlockInfo:    make([]CzBlockInfo, 0),
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

		// FIX: Correct RawSize accounting for LZW carry-over between blocks.
		//
		// compressLZW's carry-over (lastElement) always represents exactly 1 data
		// byte when non-empty, because break only triggers after the else-branch
		// which resets element to string(c) (a single input byte).
		//
		// CRITICAL: We must NOT use len(last) to count data bytes, because Go's
		// string(byte(c)) produces UTF-8: for c > 127, len(string(c)) == 2,
		// but it still represents only 1 original data byte. Using len(last)
		// would cause ±1 RawSize errors on blocks where the carry byte crosses
		// the 128 boundary.
		//
		// RawSize = data bytes carried in from previous block
		//         + data bytes consumed from input (count)
		//         - data bytes carried out to next block
		carry := 0
		if len(last) > 0 {
			carry = 1
		}
		rawSize := prevCarry + count - carry

		outputInfo.BlockInfo = append(outputInfo.BlockInfo, CzBlockInfo{
			CompressedSize: uint32(len(partData)),
			RawSize:        uint32(rawSize),
		})
		outputInfo.FileCount++
		prevCarry = carry
	}
	outputInfo.TotalCompressedSize = outputBuf.Len() / 2

	return outputBuf.Bytes(), outputInfo
}

// Compress2 压缩数据 CZ2专用
//
//	Description 压缩数据
//	Param data []byte 未压缩数据

// CompressWithRawSizes PATCH YOREMI: Compresse les données en préservant les RawSize originaux
// Cette fonction garantit que la répartition des blocs LZW est identique à celle du CZ3 original
// Cela résout le problème d'artefacts dans le jeu qui attend une structure de blocs exacte
func CompressWithRawSizes(data []byte, blockSize int, targetRawSizes []int) (compressed []byte, outputInfo *CzOutputInfo) {

	if blockSize == 0 {
		blockSize = 0xFEFD
	}
	var partData []uint16
	offset := 0
	count := 0
	last := ""
	tmp := make([]byte, 2)
	outputBuf := &bytes.Buffer{}
	outputInfo = &CzOutputInfo{
		TotalRawSize: len(data),
		BlockInfo:    make([]CzBlockInfo, 0),
	}
	
	// Compresser bloc par bloc en respectant les RawSize originaux
	for blockIdx, targetSize := range targetRawSizes {
		if offset >= len(data) {
			break
		}
		
		// Calculer la fin de ce chunk (limité aux données restantes)
		chunkEnd := offset + targetSize
		if chunkEnd > len(data) {
			chunkEnd = len(data)
		}
		
		// Compresser EXACTEMENT targetSize bytes (size=0 pour illimité)
		count, partData, last = compressLZW(data[offset:chunkEnd], 0, last)
		
		// Vérification: count devrait être égal à targetSize
		expectedCount := chunkEnd - offset
		if count != expectedCount {
			glog.V(2).Infof("Block %d: compressed %d bytes, expected %d\n", blockIdx, count, expectedCount)
		}
		
		offset = chunkEnd
		
		// Écrire les codes LZW compressés
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

//	Param size int 分块大小
//	Return compressed
//	Return outputInfo
func Compress2(data []byte, size int) (compressed []byte, outputInfo *CzOutputInfo) {

	if size == 0 {
		size = 0x87BDF
	}
	var partData []byte
	offset := 0
	count := 0
	last := ""
	prevCarry := 0
	outputBuf := &bytes.Buffer{}
	outputInfo = &CzOutputInfo{
		TotalRawSize: len(data),
		BlockInfo:    make([]CzBlockInfo, 0),
	}
	for {
		count, partData, last = compressLZW2(data[offset:], size, last)
		if count == 0 {
			break
		}
		offset += count
		outputBuf.Write(partData)

		// FIX: same carry-over correction as Compress()
		// See Compress() comments for full explanation.
		carry := 0
		if len(last) > 0 {
			carry = 1
		}
		rawSize := prevCarry + count - carry

		outputInfo.BlockInfo = append(outputInfo.BlockInfo, CzBlockInfo{
			CompressedSize: uint32(len(partData)),
			RawSize:        uint32(rawSize),
		})
		outputInfo.FileCount++
		prevCarry = carry
	}
	outputInfo.TotalCompressedSize = outputBuf.Len()
	return outputBuf.Bytes(), outputInfo
}

// ImageToNRGBA convert image.Image to image.NRGBA
func ImageToNRGBA(im image.Image) *image.NRGBA {
	dst := image.NewNRGBA(im.Bounds())
	draw.Draw(dst, im.Bounds(), im, im.Bounds().Min, draw.Src)
	return dst
}
