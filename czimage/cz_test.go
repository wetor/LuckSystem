package czimage

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
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
