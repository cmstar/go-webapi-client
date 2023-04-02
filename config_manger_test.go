package client

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const _CONFIG_PATH = "./.config_test"

func clearAllConfig() {
	fileInfo, err := os.Stat(_CONFIG_PATH)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}

		panic(err)
	}

	if !fileInfo.IsDir() {
		panic(fmt.Sprintf("not a directory %q, there may be something wrong", _CONFIG_PATH))
	}

	err = os.RemoveAll(_CONFIG_PATH)
	if err != nil {
		panic(err)
	}
}

func TestConfigManager(t *testing.T) {
	clearAllConfig()
	m := NewConfigManager(_CONFIG_PATH)
	r := require.New(t)

	r.Empty(m.ListKeys("x"))
	r.Empty(m.ListKeys("y"))

	// Save & ListKeys
	m.Save("x", "c", map[string]any{"xc": "3"})
	m.Save("x", "nil", nil)
	m.Save("x", "empty", map[string]any{})
	m.Save("x", "a", map[string]any{"xa": "1"})
	r.Equal([]string{"a", "c", "empty", "nil"}, m.ListKeys("x"))

	m.Save("z", "a", map[string]any{"za": "2"})
	m.Save("z", "b", map[string]any{"zb": "3"})
	r.Equal([]string{"a", "b"}, m.ListKeys("z"))

	// Load
	r.Equal(map[string]any{"xa": "1"}, m.Load("x", "a"))
	r.Nil(m.Load("x", "nil"))
	r.Equal(map[string]any{}, m.Load("x", "empty"))
	r.Equal(map[string]any{"xc": "3"}, m.Load("x", "c"))
	r.Equal(map[string]any{"za": "2"}, m.Load("z", "a"))
	r.Equal(map[string]any{"zb": "3"}, m.Load("z", "b"))
	r.Nil(m.Load("not-exist", "a"))

	// Remove
	m.Remove("not-exist", "a") // Do nothing.

	m.Remove("x", "a")
	r.Equal([]string{"c", "empty", "nil"}, m.ListKeys("x"))
	r.Nil(m.Load("x", "a"))

	m.Remove("z", "b")
	r.Equal([]string{"a"}, m.ListKeys("z"))
	r.Nil(m.Load("z", "b"))
}
