// Package client runs a window with the fyne Framework.
// It provides features for calling web APIs from github.com/cmstar/go-webapi.
package client

import (
	"os"
	"path"

	"fyne.io/fyne/v2"
)

const (
	DefaultWindowWidth  = 1440 // 默认的窗口宽度。
	DefaultWindowHeight = 800  // 默认的窗口高度。
)

// 描述一个 Web API 客户端的界面。
type Client interface {
	// 唯一的标识一个 [Client] ，配置的保存和读取需要用到此值。
	// 应符合 Go 变量命名规则。
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

// 返回默认的配置存储目录。默认存储在用户的 home 目录的 .go-webapi-client 子目录。
//   - 在 *nix 是 ~/.go-webapi-client
//   - 在 Windows 是 %UserProfile%\.go-webapi-client
func GetDefaultConfigDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	dir = path.Join(dir, ".go-webapi-client")
	return dir
}
