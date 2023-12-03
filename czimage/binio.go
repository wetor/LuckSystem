package czimage

import (
	"encoding/binary"
)

type BitIO struct {
	data       []byte
	byteOffset int
	byteSize   int
	bitOffset  int
	bitSize    int
}

func NewBitIO(data []byte) *BitIO {
	return &BitIO{
		data: data,
	}
}

func (b *BitIO) ByteOffset() int {
	return b.byteOffset
}

func (b *BitIO) ByteSize() int {
	return b.byteSize
}

func (b *BitIO) Bytes() []byte {
	return b.data[:b.byteSize]
}

func (b *BitIO) ReadBit(bitLen int) uint64 {
	if bitLen > 8*8 {
		panic("不支持，最多64位")
	}
	if bitLen%8 == 0 && b.bitOffset == 0 {
		return b.Read(bitLen / 8)
	}
	var result uint64
	for i := 0; i < bitLen; i++ {
		// 从最低位开始读取
		bitValue := uint64((b.data[b.byteOffset] >> uint(b.bitOffset)) & 1)
		b.bitOffset++
		if b.bitOffset == 8 {
			b.byteOffset++
			b.bitOffset = 0
		}
		// 将读取的位放入结果中
		result |= bitValue << uint(i)
	}
	//fmt.Println(b.byteOffset, b.bitOffset)
	return result
}

func (b *BitIO) Read(byteLen int) uint64 {
	if byteLen > 8 {
		panic("不支持，最多64位")
	}
	paddedSlice := make([]byte, 8)
	copy(paddedSlice, b.data[b.byteOffset:b.byteOffset+byteLen])
	b.byteOffset += byteLen
	return binary.LittleEndian.Uint64(paddedSlice)
}

func (b *BitIO) WriteBit(data uint64, bitLen int) {
	if bitLen > 8*8 {
		panic("不支持，最多64位")
	}

	if bitLen%8 == 0 && b.bitOffset == 0 {
		b.Write(data, bitLen/8)
		return
	}

	for i := 0; i < bitLen; i++ {
		// 从 value 中获取要写入的位
		bitValue := (data >> uint(i)) & 1
		// 清除目标字节中的目标位
		b.data[b.byteOffset] &= ^(1 << uint(b.bitOffset))
		// 将 bitValue 写入目标位
		b.data[b.byteOffset] |= byte(bitValue << uint(b.bitOffset))

		b.bitOffset++
		if b.bitOffset == 8 {
			b.byteOffset++
			b.bitOffset = 0
		}
	}

	b.byteSize = b.byteOffset + (b.bitOffset+7)/8
}

func (b *BitIO) Write(data uint64, byteLen int) {
	if byteLen > 8 {
		panic("不支持，最多64位")
	}
	paddedSlice := make([]byte, 8)
	binary.LittleEndian.PutUint64(paddedSlice, data)
	copy(b.data[b.byteOffset:b.byteOffset+byteLen], paddedSlice[:byteLen])
	b.byteOffset += byteLen
	b.byteSize = b.byteOffset + (b.bitOffset+7)/8
}
