package czimage

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/go-restruct/restruct"
)

func TestMain(m *testing.M) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "log")
	flag.Set("v", "10")
	flag.Parse()

	ret := m.Run()
	os.Exit(ret)
}
func TestCZ3(t *testing.T) {
	restruct.EnableExprBeta()
	for i := 0; i < 10; i++ {
		data, _ := os.ReadFile("../data/LB_EN/IMAGE/" + strconv.Itoa(i) + ".cz3")
		cz := LoadCzImage(data)

		w, _ := os.Create("../data/LB_EN/IMAGE/" + strconv.Itoa(i) + ".png")
		cz.Export(w)
		w.Close()
		fmt.Println()
	}

}
func TestCZ1(t *testing.T) {
	restruct.EnableExprBeta()
	list := []string{"明朝24", "明朝20"}
	for _, name := range list {

		data, _ := os.ReadFile("../data/LB_EN/IMAGE/" + name + ".cz1")
		cz := LoadCzImage(data)

		w, _ := os.Create("../data/LB_EN/IMAGE/" + name + ".png")
		cz.Export(w)
		w.Close()
		fmt.Println()

	}
}

func TestCZ1_2(t *testing.T) {
	restruct.EnableExprBeta()

	data, _ := os.ReadFile("../data/LB_EN/FONT/明朝32_onlyRedraw")
	cz := LoadCzImage(data)

	w, _ := os.Create("../data/LB_EN/FONT/明朝32_onlyRedraw.png")
	cz.Export(w)
	w.Close()
	fmt.Println()

}
func TestLineDiff(t *testing.T) {
	restruct.EnableExprBeta()

	data, _ := os.ReadFile("../data/LB_EN/IMAGE/2.cz3")
	cz := LoadCzImage(data)
	//cz.GetImage()

	cz3 := cz.(*Cz3Image)
	//os.WriteFile("../data/LB_EN/IMAGE/2.ld", []byte(cz3.Image.(*image.NRGBA).Pix), 0666)
	f, _ := os.Open("../data/LB_EN/IMAGE/2.png")
	defer f.Close()
	img, _ := png.Decode(f)
	pic, ok := img.(*image.NRGBA)
	if !ok {
		pic = ImageToNRGBA(img)
	}
	data1 := DiffLine(cz3.CzHeader, pic)
	fmt.Println(len(data1))
	os.WriteFile("../data/LB_EN/IMAGE/2.dl", data1, 0666)

}

//func TestCZ2(t *testing.T) {
//	restruct.EnableExprBeta()
//	data, _ := os.ReadFile("../data/Other/CZ2/ゴシック14.cz2")
//	cz := LoadCzImage(data)
//	if err != nil {
//		panic(err)
//	}
//	cz.Save("../data/Other/CZ2/ゴシック14.png")
//	fmt.Println()
//
//}

func TestCz3Image_Import(t *testing.T) {
	restruct.EnableExprBeta()
	filename := "../data/Other/CZ3/TITLE03"
	data, _ := os.ReadFile(filename)
	cz := LoadCzImage(data)

	w, _ := os.Create(filename + ".png")
	cz.GetImage()
	cz.Export(w)
	w.Close()

	r, _ := os.Open(filename + ".png")
	defer r.Close()
	w, _ = os.Create(filename + ".cz3")
	defer w.Close()
	cz.Import(r, false)
	cz.Write(w)
	fmt.Println()

}

func TestCz1Image_Import(t *testing.T) {
	restruct.EnableExprBeta()
	data, _ := os.ReadFile("../data/LB_EN/IMAGE/明朝20.cz1")
	cz := LoadCzImage(data)

	w, _ := os.Create("../data/LB_EN/IMAGE/明朝20.png")
	cz.GetImage()
	cz.Export(w)
	w.Close()

	r, _ := os.Open("../data/LB_EN/IMAGE/明朝20.png")
	defer r.Close()
	w, _ = os.Create("../data/LB_EN/IMAGE/明朝20.png.cz1")
	defer w.Close()
	cz.Import(r, false)
	cz.Write(w)
	fmt.Println()
	//data, _ = os.ReadFile("../data/LB_EN/IMAGE/明朝20.png.cz1")
	//cz, err = LoadCzImage(data)
	//if err != nil {
	//	panic(err)
	//}
	//cz.Export("../data/LB_EN/IMAGE/明朝20.cz1.png")

}
func TestCz0Image_Export(t *testing.T) {
	restruct.EnableExprBeta()
	data, _ := os.ReadFile("../data/LB_EN/FONT/明朝32")
	cz := LoadCzImage(data)

	w, _ := os.Create("../data/LB_EN/IMAGE/明朝32.png")
	cz.GetImage()
	cz.Export(w)
	w.Close()
}

func TestCz0Image_Import(t *testing.T) {
	restruct.EnableExprBeta()
	data, _ := os.ReadFile("../data/LB_EN/IMAGE/10.cz0")
	cz := LoadCzImage(data)

	w, _ := os.Create("../data/LB_EN/IMAGE/10.cz0.png")
	cz.GetImage()
	cz.Export(w)
	w.Close()

	r, _ := os.Open("../data/LB_EN/IMAGE/10.cz0.png")
	defer r.Close()
	w, _ = os.Create("../data/LB_EN/IMAGE/10.cz0.png.cz0")
	defer w.Close()
	cz.Import(r, false)
	cz.Write(w)
	fmt.Println()
}

func TestCZ2(t *testing.T) {
	restruct.EnableExprBeta()
	dir := "C:/Users/wetor/Desktop/Prototype/CZ2/32"
	list := []string{"明朝32"}
	for _, name := range list {

		data, _ := os.ReadFile(path.Join(dir, name))
		cz := LoadCzImage(data)
		w, _ := os.Create(path.Join(dir, name+".png"))
		err := cz.Export(w)
		if err != nil {
			panic(err)
		}
		w.Close()
		fmt.Println()

	}
}

func TestCZ1112(t *testing.T) {
	rdi := 3
	fmt.Println(int64(rdi)&^10, int64(rdi)&240)
}

func TestCZ22(t *testing.T) {
	restruct.EnableExprBeta()
	dir := "C:/Users/wetor/Desktop/Prototype/"
	list := []string{"cz2_12"}
	for _, name := range list {

		data, _ := os.ReadFile(path.Join(dir, name))
		cz := LoadCzImage(data)
		w, _ := os.Create(path.Join(dir, name+".png"))
		err := cz.Export(w)
		if err != nil {
			panic(err)
		}
		w.Close()
		fmt.Println()

	}
}

func TestLZWo(t *testing.T) {
	compressedData := []byte{0x00, 0x00, 0x02, 0x02, 0x04, 0x02, 0x06, 0x02, 0x08, 0x02, 0x0A, 0x02, 0x0C, 0x02, 0x0E, 0x02, 0x10, 0x02, 0x12, 0x02}
	uint16s := readUint16File(compressedData)
	// 创建一个Reader来读取压缩数据
	decompressedData := decompressLZW22(uint16s, 2000)
	fmt.Println(decompressedData)

	compressedData = []byte{0x1E, 0x01, 0xFE, 0x01, 0xFE, 0x01, 0xF0, 0x01, 0xE2, 0x01, 0xD2, 0x01,
		0x0C, 0x02, 0x0E, 0x02, 0x0C, 0x02, 0xEA, 0x00, 0x12, 0x00, 0x04, 0x02, 0x12, 0x00, 0x80,
		0x00, 0x8E, 0x01, 0x9C, 0x00, 0x00, 0x00, 0x22, 0x02, 0x24, 0x02, 0x26, 0x02, 0x14, 0x00,
		0x14, 0x00, 0x4C, 0x00, 0xAC, 0x01}

	uint16s = readUint16File(compressedData)
	// 创建一个Reader来读取压缩数据
	decompressedData = decompressLZW22(uint16s, 10000)
	fmt.Println(decompressedData)
}

// readUint16File 从文件中读取[]uint16数据
func readUint16File(data []byte) []uint16 {

	var result []uint16
	for i := 0; i < len(data); i += 2 {
		val := binary.LittleEndian.Uint16(data[i : i+2])
		result = append(result, val)
	}

	return result
}

func TestLZWCE(t *testing.T) {
	ptr := []byte{0xF4, 0x24, 0xCA, 0x6E, 0x9F, 0x00, 0x61, 0x2C,
		0x42, 0x76, 0x93, 0xCE}

	dst := []byte{0xFF, 0xFF, 0xEF,
		0x1E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x13, 0xD8, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xEF, 0x1E,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x98}

	result := DecompressLZWByAsm(ptr, len(dst))
	if bytes.Compare(result, dst) != 0 {
		fmt.Printf("result:\t%v,\ndst:\t%v\n", result, dst)
	}
}

func TestLZWdict(t *testing.T) {
	ptr := []byte{0x00, 0x00, 0x02, 0x02, 0x04, 0x02, 0x06, 0x02,
		0x08, 0x02, 0x0A, 0x02, 0x0C, 0x02, 0x0E, 0x02, 0x10, 0x02,
		0x12, 0x02, 0x14, 0x02, 0x16, 0x02, 0x18, 0x02, 0x1A, 0x02,
		0x1C, 0x02, 0x1E, 0x02, 0x20, 0x02, 0x22, 0x02, 0x24, 0x02,
		0x26, 0x02, 0x28, 0x02, 0x2A, 0x02, 0x2C, 0x02, 0x2E, 0x02,
		0x30, 0x02, 0x32, 0x02, 0x34, 0x02, 0x36, 0x02, 0x38, 0x02}

	dst := []byte{0xFF, 0xFF, 0xEF,
		0x1E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x13, 0xD8, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xEF, 0x1E,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x98}

	result := DecompressLZWByAsm(ptr, len(dst))
	if bytes.Compare(result, dst) != 0 {
		fmt.Printf("result:\t%v,\ndst:\t%v\n", result, dst)
	}
}

func TestLZWdict2(t *testing.T) {

	data, _ := os.ReadFile("C:\\Users\\wetor\\Desktop\\Prototype\\CZ2\\0.src.lzw")

	result := DecompressLZWByAsm(data, 1892588)
	os.WriteFile("C:\\Users\\wetor\\Desktop\\Prototype\\CZ2\\0.lzw.out", result, 0666)
}
