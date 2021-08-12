package font

import (
	"fmt"
	"image/png"
	"lucascript/charset"
	"lucascript/pak"
	"os"
	"testing"

	"github.com/go-restruct/restruct"
)

func TestFont(t *testing.T) {

	restruct.EnableExprBeta()
	pak := pak.NewPak(&pak.PakFileOptions{
		FileName: "../data/LB_EN/FONT.PAK",
		Coding:   charset.UTF_8,
	})
	err := pak.Open()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("pak header", pak.PakHeader)

	font := LoadLucaFont(pak, "モダン", 36)
	img := font.GetStringImage("理樹@「…は？　こんな夜に？　どこで？」")
	f, _ := os.Create("../data/LB_EN/IMAGE/str/str.png")
	png.Encode(f, img)
	f.Close()
	// imgs := font.GetStringImageList("真人@「…戦いさ」")

	// for i, img := range imgs {
	// 	f, _ := os.Create("../data/LB_EN/IMAGE/str/" + strconv.Itoa(i) + ".png")
	// 	png.Encode(f, img)
	// 	f.Close()
	// }

}
