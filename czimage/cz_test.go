package czimage

import (
	"crypto/md5"
	"encoding/hex"
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

func TestLZWdict2(t *testing.T) {

	lzwData, err := os.ReadFile("testdata/明朝32_0.lzw")
	if err != nil {
		panic(err)
	}
	data, err := os.ReadFile("testdata/明朝32_0.bin")
	if err != nil {
		panic(err)
	}
	h := md5.New()
	h.Write(data)
	dstMD5 := hex.EncodeToString(h.Sum(nil))

	result := decompressLZW2(lzwData, len(data))
	h.Reset()
	h.Write(result)
	resMD5 := hex.EncodeToString(h.Sum(nil))

	if resMD5 != dstMD5 {
		panic("不匹配")
	} else {
		fmt.Println(resMD5, dstMD5)
	}
}

func TestNewBinIO(t *testing.T) {
	bytes := []byte{0xFF, 0xF0, 0xFF, 0xFF, 0xFF}
	b := NewBitIO(bytes)
	for i := 0; i < len(bytes)*8/12; i++ {
		fmt.Println(b.ReadBit(12))
	}

}
