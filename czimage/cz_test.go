package czimage

import (
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
		cz.Save("../data/LB_EN/IMAGE/" + strconv.Itoa(i) + ".png")
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
		cz.Save("../data/LB_EN/IMAGE/" + name + ".png")
		fmt.Println()

	}

}
func TestLineDiff(t *testing.T) {
	restruct.EnableExprBeta()

	data, _ := os.ReadFile("../data/LB_EN/IMAGE/2.cz3")
	cz, err := LoadCzImage(data)
	if err != nil {
		panic(err)
	}
	cz3 := cz.(*Cz3Image)
	//os.WriteFile("../data/LB_EN/IMAGE/2.ld", []byte(cz3.Image.(*image.RGBA).Pix), 0666)
	f, _ := os.Open("../data/LB_EN/IMAGE/2.png")
	defer f.Close()
	img, _ := png.Decode(f)
	data1 := DiffLine(&cz3.CzHeader, img)
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
