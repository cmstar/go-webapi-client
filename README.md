# webapi client

一个简单的 GUI 客户端，用于调用 [go-webapi](https://github.com/cmstar/go-webapi) 库中相应协议开发的 Web API 。


## 安装

### 从源码

图形界面基于 [Fyne](https://fyne.io) 开发，编译源码需要本地有 `gcc` ，如果你本地没有，
可参考[这里](https://developer.fyne.io/started/#prerequisites)安装， Windows 建议安装 TDM-GCC 。

准备好环境后，执行：
```bash
go install github.com/cmstar/go-webapi-client/cmd/webapi-client@latest
```

在 Windows ，若直接双击运行 exe 会额外看到看到一个 cmd 窗口，如果不想看到这个窗口，可在安装时添加 `-H=windowsgui` 参数：
```bash
go install -ldflags -H=windowsgui github.com/cmstar/go-webapi-client/cmd/webapi-client@latest
```

编译较慢，需要耐心等待。


## 修改配置文件的存储路径

默认情况下，配置文件会被存储在用户的 home 目录的 .go-webapi-client 子目录：
- 在 Mac/Linux 是 `~/.go-webapi-client` 。
- 在 Windows 是 `%UserProfile%\.go-webapi-client` 。

可以通过启动程序时添加 `-c` 参数来指定你想要的位置：
```bash
webapi-client -c=/my/favor/path
```


## 功能扩展

主窗口 `client.MainWindow` 支持在多个 `client.Client` 间的切换，每个 `Client` 表示一个界面。

```go
package main

import (
	client "github.com/cmstar/go-webapi-client"
)

// 自己实现的 Client 接口。
var (
	myClient1 = NewMyClient1()
	myClient2 = NewMyClient2()
	myClient3 = NewMyClient3()
)

// 展示主窗体，可以通过菜单切换展示上面的三个 Client 。
func main() {
	client.RunClients([]client.Client{
		myClient1,
		myClient2,
		myClient3,
	})
}
```
