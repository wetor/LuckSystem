package czimage

import (
	"fmt"
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
	fmt.Println(1)
}
