package slimauth_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/cmstar/go-errx"
	"github.com/cmstar/go-httplib/headers"
	client "github.com/cmstar/go-webapi-client"
	"github.com/cmstar/go-webapi/slimauth"
)

const (
	_KEY    = "Key"
	_SECRET = "Secret"
	_URI    = "Uri"
	_PARAM  = "Param"
)

type SlimAuthClient struct {
	// 这几个值都是单向绑定，直接用 string 即可。
	key   string
	sec   string
	uri   string
	param string

	// 用于更新 res 。
	mu sync.Mutex

	// 直接用 string 控件不能实现双向绑定。
	res binding.String
}

var _ client.Client = (*SlimAuthClient)(nil)

// 创建一个 [*SlimAuthClient] 。
func NewClient() *SlimAuthClient {
	return &SlimAuthClient{
		res: binding.NewString(),
	}
}

func (x *SlimAuthClient) Name() string {
	return "SlimAuth"
}

func (x *SlimAuthClient) Title() string {
	return "SlimAuth"
}

func (x *SlimAuthClient) GetConfig() map[string]any {
	return map[string]any{
		_KEY:    x.key,
		_SECRET: x.sec,
		_URI:    x.uri,
		_PARAM:  x.param,
	}
}

func (x *SlimAuthClient) SetConfig(config map[string]any) {
	x.key = config[_KEY].(string)
	x.sec = config[_SECRET].(string)
	x.uri = config[_URI].(string)
	x.param = config[_PARAM].(string)
}

func (x *SlimAuthClient) Box() fyne.CanvasObject {
	paramInput := widget.NewMultiLineEntry()
	paramInput.Bind(binding.BindString(&x.param))

	requestForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Key", Widget: widget.NewEntryWithData(binding.BindString(&x.key))},
			{Text: "Secret", Widget: widget.NewEntryWithData(binding.BindString(&x.sec))},
			{Text: "URL", Widget: widget.NewEntryWithData(binding.BindString(&x.uri))},
			{Text: "Param", Widget: paramInput},
		},
		OnSubmit: x.onSubmit,
	}

	responseBox := widget.NewMultiLineEntry()
	responseBox.Bind(x.res)

	container := container.NewHSplit(
		container.NewVScroll(requestForm),
		container.NewVScroll(responseBox),
	)

	return container
}

func (x *SlimAuthClient) onSubmit() {
	// 采用异步请求。
	x.res.Set("requesting ...")

	go func() {
		// 发现更新 binding.String 速度太快会来不及反馈到界面上。等一下下。
		<-time.After(200 * time.Millisecond)

		responseText, err := x.performRequest()

		x.mu.Lock()
		if err != nil {
			x.res.Set(err.Error())
		} else {
			x.res.Set(responseText)
		}
		x.mu.Unlock()
	}()
}

func (x *SlimAuthClient) performRequest() (responseText string, err error) {
	defer func() {
		if err == nil {
			err = errx.PreserveRecover("", recover())
		}
	}()

	// Body must be a JSON.
	if !json.Valid([]byte(x.param)) {
		return "", fmt.Errorf("the request message is not a valid JSON")
	}

	requestBody := strings.NewReader(x.param)
	request, err := http.NewRequest(http.MethodPost, x.uri, requestBody)
	if err != nil {
		return "", err
	}

	request.Header.Set(headers.ContentType, "application/json")
	signResult := slimauth.AppendSign(request, x.key, x.sec, "", time.Now().Unix())
	if signResult.Type != slimauth.SignResultType_OK {
		return "", signResult.Cause
	}

	response, err := new(http.Client).Do(request)
	if err != nil {
		return "", err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	responseText, _ = x.identJson(responseBody)
	return
}

// 尝试格式化 JSON 。若给定过的不是合法的 JSON ，返回原值的字符串形式 + ok=false。
func (x *SlimAuthClient) identJson(v []byte) (res string, ok bool) {
	const ident = "    "
	buf := new(bytes.Buffer)
	err := json.Indent(buf, v, "", ident)
	if err != nil {
		return "", false
	}
	return buf.String(), true
}
