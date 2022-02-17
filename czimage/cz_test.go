package czimage

import (
	"flag"
	"fmt"
	"image/png"
	"os"
	"strconv"
	"testing"

	"github.com/go-restruct/restruct"
)

func TestCZ3(t *testing.T) {
	restruct.EnableExprBeta()
	for i := 0; i < 10; i++ {
		data, _ := os.ReadFile("../data/LB_EN/IMAGE/" + strconv.Itoa(i) + ".cz3")
		cz, err := LoadCzImage(data)
		if err != nil {
			panic(err)
		}
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
		cz, err := LoadCzImage(data)
		if err != nil {
			panic(err)
		}
		w, _ := os.Create("../data/LB_EN/IMAGE/" + name + ".png")
		cz.Export(w)
		w.Close()
		fmt.Println()

	}

}
func TestCZ1_2(t *testing.T) {
	restruct.EnableExprBeta()

	data, _ := os.ReadFile("../data/LB_EN/FONT/明朝32_onlyRedraw")
	cz, err := LoadCzImage(data)
	if err != nil {
		panic(err)
	}
	w, _ := os.Create("../data/LB_EN/FONT/明朝32_onlyRedraw.png")
	cz.Export(w)
	w.Close()
	fmt.Println()

}
func TestLineDiff(t *testing.T) {
	restruct.EnableExprBeta()

	data, _ := os.ReadFile("../data/LB_EN/IMAGE/2.cz3")
	cz, err := LoadCzImage(data)
	//cz.GetImage()
	if err != nil {
		panic(err)
	}
	cz3 := cz.(*Cz3Image)
	//os.WriteFile("../data/LB_EN/IMAGE/2.ld", []byte(cz3.Image.(*image.RGBA).Pix), 0666)
	f, _ := os.Open("../data/LB_EN/IMAGE/2.png")
	defer f.Close()
	img, _ := png.Decode(f)
	data1 := DiffLine(cz3.CzHeader, img)
	fmt.Println(len(data1))
	os.WriteFile("../data/LB_EN/IMAGE/2.dl", data1, 0666)

}

//func TestCZ2(t *testing.T) {
//	restruct.EnableExprBeta()
//	data, _ := os.ReadFile("../data/Other/CZ2/ゴシック14.cz2")
//	cz, err := LoadCzImage(data)
//	if err != nil {
//		panic(err)
//	}
//	cz.Save("../data/Other/CZ2/ゴシック14.png")
//	fmt.Println()
//
//}

func TestCz3Image_Import(t *testing.T) {
	restruct.EnableExprBeta()
	data, _ := os.ReadFile("../data/LB_EN/IMAGE/4.cz3")
	cz, err := LoadCzImage(data)
	if err != nil {
		panic(err)
	}
	w, _ := os.Create("../data/LB_EN/IMAGE/4.png")
	cz.Export(w)
	w.Close()

	r, _ := os.Open("../data/LB_EN/IMAGE/4.png")
	defer r.Close()
	w, _ = os.Create("../data/LB_EN/IMAGE/4.png.cz3")
	defer w.Close()
	cz.Import(r)
	cz.Write(w)
	fmt.Println()

}

func TestCz1Image_Import(t *testing.T) {
	restruct.EnableExprBeta()
	data, _ := os.ReadFile("../data/LB_EN/IMAGE/明朝20.cz1")
	cz, err := LoadCzImage(data)
	if err != nil {
		panic(err)
	}
	w, _ := os.Create("../data/LB_EN/IMAGE/明朝20.png")
	cz.Export(w)
	w.Close()

	r, _ := os.Open("../data/LB_EN/IMAGE/明朝20.png")
	defer r.Close()
	w, _ = os.Create("../data/LB_EN/IMAGE/明朝20.png.cz1")
	defer w.Close()
	cz.Import(r)
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
	data, _ := os.ReadFile("../data/LB_EN/IMAGE/10.cz0")
	cz, err := LoadCzImage(data)
	if err != nil {
		panic(err)
	}

	w, _ := os.Create("../data/LB_EN/IMAGE/10.cz0.png")
	cz.Export(w)
	w.Close()
}
func TestMain(m *testing.M) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "log")
	flag.Set("v", "10")
	flag.Parse()

	ret := m.Run()
	os.Exit(ret)
}
func TestCz0Image_Import(t *testing.T) {
	restruct.EnableExprBeta()
	data, _ := os.ReadFile("../data/LB_EN/IMAGE/10.cz0")
	cz, err := LoadCzImage(data)
	if err != nil {
		panic(err)
	}

	w, _ := os.Create("../data/LB_EN/IMAGE/10.cz0.png")
	cz.Export(w)
	w.Close()

	r, _ := os.Open("../data/LB_EN/IMAGE/10.cz0.png")
	defer r.Close()
	w, _ = os.Create("../data/LB_EN/IMAGE/10.cz0.png.cz0")
	defer w.Close()
	cz.Import(r, w)
	fmt.Println()
}
