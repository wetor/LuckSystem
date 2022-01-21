package czimage

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"testing"
)

func TestLzwCompress(t *testing.T) {
	data, _ := os.ReadFile("../data/LB_EN/IMAGE/9.ori")
	n, com, _ := compressLZW(data, 0, "")
	fmt.Println(len(com), n)
	buf := bytes.NewBuffer(nil)
	tmp := make([]byte, 2)
	for _, d := range com {
		binary.LittleEndian.PutUint16(tmp, d)
		buf.Write(tmp)
	}
	os.WriteFile("../data/LB_EN/IMAGE/9.new.lzw", buf.Bytes(), 0666)
	//65276 386904
}
func TestLzwDeCompress(t *testing.T) {

	data, _ := os.ReadFile("../data/LB_EN/IMAGE/9.new.lzw")

	com := make([]uint16, len(data)/2)
	for j := 0; j < len(data); j += 2 {
		com[j/2] = binary.LittleEndian.Uint16(data[j : j+2])
	}
	decom := decompressLZW(com, len(data))
	fmt.Println(len(decom))

	os.WriteFile("../data/LB_EN/IMAGE/9.new", decom, 0666)

}
func TestLzwCompressPart(t *testing.T) {
	//65000 385376
	//657 1528
	data, _ := os.ReadFile("../data/LB_EN/IMAGE/9.ori")
	n, com, l := compressLZW(data, 65000, "")
	fmt.Println(len(com), n)
	buf := bytes.NewBuffer(nil)
	tmp := make([]byte, 2)
	for _, d := range com {
		binary.LittleEndian.PutUint16(tmp, d)
		buf.Write(tmp)
	}
	os.WriteFile("../data/LB_EN/IMAGE/9.new.lzw", buf.Bytes(), 0666)

	n, com, l = compressLZW(data[n:], 65000, l)
	fmt.Println(len(com), n)
	buf = bytes.NewBuffer(nil)
	for _, d := range com {
		binary.LittleEndian.PutUint16(tmp, d)
		buf.Write(tmp)
	}
	os.WriteFile("../data/LB_EN/IMAGE/9.new.2.lzw", buf.Bytes(), 0666)

}
func TestLzwDeCompressPart(t *testing.T) {

	data, _ := os.ReadFile("../data/LB_EN/IMAGE/9.new.lzw")

	com := make([]uint16, len(data)/2)
	for j := 0; j < len(data); j += 2 {
		com[j/2] = binary.LittleEndian.Uint16(data[j : j+2])
	}
	decom := bytes.NewBuffer(nil)
	decom.Write(decompressLZW(com, len(data)))
	fmt.Println(decom.Len())

	data, _ = os.ReadFile("../data/LB_EN/IMAGE/9.new.2.lzw")

	com = make([]uint16, len(data)/2)
	for j := 0; j < len(data); j += 2 {
		com[j/2] = binary.LittleEndian.Uint16(data[j : j+2])
	}
	decom.Write(decompressLZW(com, len(data)))
	fmt.Println(decom.Len())
	os.WriteFile("../data/LB_EN/IMAGE/9.new", decom.Bytes(), 0666)

}

func TestCompress(t *testing.T) {
	data, _ := os.ReadFile("../data/LB_EN/IMAGE/2.dl")
	f, info := Compress(data, 0)
	fmt.Println(info)
	os.WriteFile("../data/LB_EN/IMAGE/2.dl.lzw", f, 0666)
}
