package client

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
)

// 用于读写配置。
//
// 配置存放在 rootPath 所指向的目录下，每个 ClientName 一个子目录，
// 每个 key 一个 .json 文件。
//
//	root
//	|- ClientName1
//	|  |- key1.json
//	|  |- key2.json
//	|- ClientName2
//	|  |- key1.json
//	|  |- key2.json
//
// 注意：
//   - key 可以在不同的 ClientName 下重复。
//   - Windows 平台的文件名是大小写不敏感的；*nix 则是敏感的。
type ConfigManager struct {
	rootPath string
}

// 创建一个 [ConfigManager] ，给定存放配置文件的根目录的路径。
func NewConfigManager(rootPath string) *ConfigManager {
	return &ConfigManager{rootPath}
}

// 返回一个 clientName 下的所有配置项的 key ，按字典顺序排列。
//
// clientName 子目录不存在时，返回 nil ；无效的配置会被忽略。
func (x *ConfigManager) ListKeys(clientName string) []string {
	p := x.getClientDirPath(clientName)
	files, err := os.ReadDir(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		panic(err)
	}

	trimExt := func(name, ext string) string {
		if len(name) < len(ext) {
			return ""
		}

		tail := name[len(name)-len(ext):]
		if !strings.EqualFold(tail, ext) {
			return ""
		}

		return name[:len(name)-len(ext)]
	}

	keys := make([]string, 0, len(files))
	for _, v := range files {
		key := trimExt(v.Name(), ".json")
		if key == "" {
			continue
		}
		keys = append(keys, key)
	}

	// os.ReadDir 读取出来本来应该是排序好的，但 API 并没有这个保证，这里再排序一下。
	sort.Strings(keys)
	return keys
}

// 读取一个 clientName 下指定 key 的配置的值。
// 不会返回 nil ，若对应配置不存， 返回 nil 。
// 应先通过 ListKeys 获取相关的数据。
func (x *ConfigManager) Load(clientName, key string) map[string]any {
	p := x.getKeyFilePath(clientName, key)
	content, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		panic(err)
	}

	var res map[string]any
	err = json.Unmarshal(content, &res)
	if err != nil {
		panic(err)
	}

	return res
}

// 保存一个配置。若 key 在 clientName 下的配置中已存在，则覆盖原配置。
// 若配置目录或文件不存在，会被创建出来。
func (x *ConfigManager) Save(clientName, key string, conf map[string]any) {
	// 确保目录存在。
	dir := x.getClientDirPath(clientName)
	err := os.MkdirAll(dir, 0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	content, err := json.Marshal(conf)
	if err != nil {
		panic(err)
	}

	p := x.getKeyFilePath(clientName, key)
	err = os.WriteFile(p, content, 0644)
	if err != nil {
		panic(err)
	}
}

// 移除 clientName 下指定 key 的配置。若配置不存在，操作被忽略。
func (x *ConfigManager) Remove(clientName, key string) {
	p := x.getKeyFilePath(clientName, key)
	err := os.Remove(p)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
}

func (x *ConfigManager) getClientDirPath(clientName string) string {
	x.mustBeValidName(clientName)

	res := path.Join(x.rootPath, clientName)
	return res
}

func (x *ConfigManager) getKeyFilePath(clientName, key string) string {
	x.mustBeValidName(clientName)
	x.mustBeValidName(key)

	res := path.Join(x.rootPath, clientName, key+".json")
	return res
}

// 校验文件名的有效性，以 Windows 为准，它的限制比较多， Linux 只要求不要包含斜杠。
func (x *ConfigManager) mustBeValidName(v string) {
	if len(v) == 0 {
		panic("name cannot be empty")
	}

	onlyDot := true
	for _, c := range v {
		if c == '\\' ||
			c == '/' ||
			c == ':' ||
			c == '*' ||
			c == '?' ||
			c == '"' ||
			c == '<' ||
			c == '>' ||
			c == '|' {
			panic(fmt.Sprintf(`the given name %q is not a valid file name`, v))
		}

		if c != '.' {
			onlyDot = false
		}
	}

	if onlyDot {
		panic(fmt.Sprintf(`the given name %s is a relative name`, v))
	}
}
