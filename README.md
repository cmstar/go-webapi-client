# webapi client

一个简单的 GUI 客户端，用于调用 [go-webapi](https://github.com/cmstar/go-webapi) 库中相应协议开发的 Web API 。

完善中……

## 安装

### 从源码

GUI 框架基于 [Fyne](https://fyne.io) 开发，编译源码需要本地有 `gcc` ，如果你本地没有，
可参考[这里](https://developer.fyne.io/started/#prerequisites)安装， Windows 建议安装 TDM-GCC 。

准备好环境后，执行：
```bash
go install github.com/cmstar/go-webapi-client/cmd/webapi-client@latest
```

在 Windows ，若直接双击运行 exe 会额外看到看到一个 cmd 窗口，如果不想看到这个窗口，可在安装时添加 `-H=windowsgui` 参数：
```bash
go install -ldflags -H=windowsgui github.com/cmstar/go-webapi-client/cmd/webapi-client@latest
```
