// Package client runs a window with the fyne Framework.
// It provides features for calling web APIs from github.com/cmstar/go-webapi.
package client

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

// 描述一个 Web API 客户端的界面。
type Client interface {
	// 唯一的标识一个 [Client] ，配置的保存和读取需要用到此值。
	Name() string

	// 界面展示的标题。
	Title() string

	// 读取当前界面的元素。
	Box() fyne.CanvasObject

	// 读取当前界面的配置。用于配置的保存和读取功能。
	GetConfig() map[string]any

	// 设置当前界面的配置。用于配置的保存和读取功能。
	SetConfig(config map[string]any)
}

type RunOption struct {
	Width   float32  // 窗口的宽度。若为0，套用默认值。
	Height  float32  // 窗口的高度。若为0，套用默认值。
	Clients []Client // 要展示的 [Client] 。
}

// 运行程序。
func Run(op *RunOption) {
	if len(op.Clients) == 0 {
		panic("there must be at least one Client")
	}

	a := app.New()
	w := a.NewWindow("Test client for go-webapi")

	w.Resize(fyne.NewSize(
		determineSize(op.Width, 880),
		determineSize(op.Height, 550),
	))

	client := op.Clients[0]
	w.SetContent(client.Box())
	w.ShowAndRun()
}

func determineSize(v, dft float32) float32 {
	if v > 0 {
		return v
	}
	return dft
}
