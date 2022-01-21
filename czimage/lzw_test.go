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
	com := compressLZW(data)
	fmt.Println(len(com))
	buf := bytes.NewBuffer(nil)
	tmp := make([]byte, 2)
	for _, d := range com {
		binary.LittleEndian.PutUint16(tmp, d)
		buf.Write(tmp)
	}
	os.WriteFile("../data/LB_EN/IMAGE/9.new.lzw", buf.Bytes(), 0666)
	//fmt.Println(buf)

	// 解压
	//r := lzw.NewReader(buf, lzw.LSB, 8)
	//defer r.Close()
	//io.Copy(os.Stdout, r)
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
