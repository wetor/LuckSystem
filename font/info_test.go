package font

import (
	"flag"
	"fmt"
	"os"
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
func TestInfo(t *testing.T) {
	restruct.EnableExprBeta()
	list := []string{"info32"}
	for _, name := range list {
		file := "../data/LB_EN/IMAGE/" + name
		data, _ := os.ReadFile(file)
		info := LoadFontInfo(data)
		info.Export(file + ".txt")
		fmt.Println(info.CharNum, " ", len(info.IndexUnicode))
		info.Write(file + ".out")
	}
}
func TestInfo2(t *testing.T) {
	restruct.EnableExprBeta()
	file := "../data/Other/Font/info32e.info"
	data, _ := os.ReadFile(file)
	info := LoadFontInfo(data)
	info.Export(file + ".txt")
	fmt.Println(info.CharNum, " ", len(info.IndexUnicode))
	info.Write(file + ".out")

}
func TestStr(t *testing.T) {
	aaa := []int{1, 2, 3}
	index := 3
	bbb := make([]int, 10)
	copy(bbb, aaa[:index])

	fmt.Println(bbb)
}
