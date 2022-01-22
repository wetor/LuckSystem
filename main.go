package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/rivo/tview"
)

// glog level:
//   1
//   2 提示信息
//   3 vm运行信息
//   4 vm调试信息
//   5 vm错误信息，不影响运行
//   6 调试信息，一些运行时输出
//   7 错误信息，不影响运行
//   8 错误信息（不panic，完全可忽略）

func main() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "log")
	flag.Set("v", "10")
	flag.Parse()
	glog.V(8).Infoln("helloworld")
	app := tview.NewApplication()
	list := tview.NewList().
		AddItem("List item 1", "Some explanatory text", 'a', nil).
		AddItem("List item 2", "Some explanatory text", 'b', nil).
		AddItem("List item 3", "Some explanatory text", 'c', nil).
		AddItem("List item 4", "Some explanatory text", 'd', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})
	if err := app.SetRoot(list, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
