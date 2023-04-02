package client

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

/*
主窗口，形态为：

----------------------------------------------
|<WindowTitle>                               |
|--------------------------------------------| <-|
|Client Config Help                          |   <Menu>
|--------------------------------------------| <-|
|ClientTitle |                               |   |
|ConfigList  |                               |   |
| |- config1 |                               |   |
| |- config2 |                               |   |
| |- config3 |                               |   |
| |- config4 |         client.Box()          |   <MainContent>
|            |                               |   |
|------------|                               |   |
|Config      |                               |   |
|Operation   |                               |   |
---------------------------------------------- <-|
^            ^                               ^
|<ConfigArea>|          <ClientBox>          |

结构层级如下：
	<Window>						主窗体。
		<WindowTitle>				窗体标题，可跟着选中的 Client 变化。
		<Menu>						菜单。
			<Client>				可以在这个菜单里选择要展示哪个 Client ，每个 Client 一个菜单项。
		<MainContent>				窗体的主容器，当前展示的 Client 的配置和主界面。
			<ConfigArea>			当前 Client 的配置。
				<ClientTitle>		展示当前的 Client.Title() 。
				<ConfigList>		当前 Client 的配置列表，每个 Client 可以有一组配置，基于 Client.Name() 从配置文件里获取。
				<ConfigOperation>	对于当前配置的操作：保存、删除、移动。
			<ClientBox>				展示当前的 Client.Box() 。
*/

// 程序的主界面。
type MainWindow struct {
	configManager *ConfigManager
	app           fyne.App
	win           fyne.Window
	width         float32
	height        float32
	clients       []Client

	// <ConfigArea> 的数据，每次切换 Client 时初始化。
	configAreaData struct {
		title       binding.String             // 绑定当前 Client 的 Title() 。
		configKeys  []string                   // 当前 Client 的所有配置的 key 。
		keys        binding.ExternalStringList // 绑定 configKeys 。
		selectedKey binding.String             // configKeys 中当前被选中的 key 。
	}

	// <ClientBox> 的数据，每次切换 Client 时初始化。
	clientBoxData struct {
		client    Client          // 当前的 Client 。
		container *fyne.Container // 装 client.Box() 的容器。
	}
}

type MainWindowOption struct {
	ConfigPath string   // 指定存储配置的目录。若为空，则使用 [GetDefaultConfigDir] 。
	Width      float32  // 窗口的宽度。若为0，套用 [DefaultWindowWidth] 。
	Height     float32  // 窗口的高度。若为0，套用 [DefaultWindowHeight] 。
	Clients    []Client // 给定可在窗口内切换的 [Client] 。
}

// 创建主界面。
func NewMainWindow(op *MainWindowOption) *MainWindow {
	if len(op.Clients) == 0 {
		panic("there must be at least one Client")
	}

	configPath := op.ConfigPath
	if configPath == "" {
		configPath = GetDefaultConfigDir()
	}

	a := app.New()
	w := a.NewWindow("Test client for go-webapi")
	m := &MainWindow{
		configManager: NewConfigManager(configPath),
		app:           a,
		win:           w,
		width:         op.Width,
		height:        op.Height,
		clients:       op.Clients,
	}
	m.configAreaData.title = binding.NewString()
	m.configAreaData.selectedKey = binding.NewString()
	m.configAreaData.keys = binding.BindStringList(&[]string{})
	m.clientBoxData.container = container.NewMax()

	if m.width <= 0 {
		m.width = DefaultWindowWidth
	}

	if m.height <= 0 {
		m.width = DefaultWindowHeight
	}

	w.SetMainMenu(m.makeMenu())
	w.SetContent(m.makeMainContent())
	return m
}

// 运行并显示窗口。
func (x *MainWindow) ShowAndRun() {
	x.showClient(x.clients[0])
	x.win.Resize(fyne.NewSize(x.width, x.height))
	x.win.ShowAndRun()
}

func (x *MainWindow) makeMenu() *fyne.MainMenu {
	clientItems := make([]*fyne.MenuItem, 0, len(x.clients))
	for _, c := range x.clients {
		menu := fyne.NewMenuItem(c.Name(), func() {
			x.showClient(c)
		})
		clientItems = append(clientItems, menu)
	}

	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("Client", clientItems...),
	)
	return mainMenu
}

func (x *MainWindow) makeMainContent() fyne.CanvasObject {
	content := container.NewHSplit(
		x.makeConfigArea(),
		x.clientBoxData.container,
	)

	// <ConfigArea> 比较小，主体空间留给 <ClientBox> 。
	content.Offset = 0.2
	return content
}

func (x *MainWindow) makeConfigArea() fyne.CanvasObject {
	createListItem := func() fyne.CanvasObject {
		// 每项样式为： [ICON] LABEL
		return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel(""))
	}
	updateListItem := func(i binding.DataItem, o fyne.CanvasObject) {
		// o 对应上面 createListItem 返回的容器，索引0为 ICON ，索引1为 LABEL 部分。
		label := o.(*fyne.Container).Objects[1].(*widget.Label)
		label.Bind(i.(binding.String))
	}
	list := widget.NewListWithData(x.configAreaData.keys, createListItem, updateListItem)

	// 当选中一个配置，将配置应用到 <ClientBox> 上。
	list.OnSelected = func(id widget.ListItemID) {
		item, _ := x.configAreaData.keys.GetItem(id)
		key, _ := item.(binding.String).Get()

		conf := x.configManager.Load(x.clientBoxData.client.Name(), key)
		if conf == nil {
			panic("cannot load the config, the file may have been removed")
		}

		x.configAreaData.selectedKey.Set(key)
		x.clientBoxData.client.SetConfig(conf)
	}

	return container.NewBorder(
		/* top		*/ widget.NewLabelWithData(x.configAreaData.title),
		/* bottom	*/ x.makeConfigOperation(),
		/* left		*/ nil,
		/* right	*/ nil,
		/* center	*/ list,
	)
}

func (x *MainWindow) makeConfigOperation() fyne.CanvasObject {
	btnSave := widget.NewButton("SAVE", func() {
		key, _ := x.configAreaData.selectedKey.Get()
		if key == "" {
			return
		}

		c := x.clientBoxData.client
		x.configManager.Save(c.Name(), key, c.GetConfig())
		x.reloadConfig(c.Name())
	})

	btnDelete := widget.NewButton("DELETE", func() {
		key, _ := x.configAreaData.selectedKey.Get()
		if key == "" {
			return
		}

		c := x.clientBoxData.client
		x.configManager.Remove(c.Name(), key)

		x.reloadConfig(c.Name())
	})

	return container.NewVBox(
		widget.NewSeparator(),
		widget.NewEntryWithData(x.configAreaData.selectedKey),
		container.NewGridWithColumns(2, btnSave, btnDelete),
	)
}

func (x *MainWindow) showClient(client Client) {
	x.reloadConfig(client.Name())

	x.configAreaData.title.Set(client.Title())
	x.configAreaData.selectedKey.Set("")

	// 如果 container.Objects 没有发生变化， fyne 不会刷新界面。
	x.clientBoxData.client = client
	x.clientBoxData.container.Objects = []fyne.CanvasObject{client.Box()}
}

func (x *MainWindow) reloadConfig(clientName string) {
	keys := x.configManager.ListKeys(clientName)
	x.configAreaData.configKeys = keys
	x.configAreaData.keys.Set(keys)
}
