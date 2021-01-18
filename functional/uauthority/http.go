package uauthority

import (
	"bytes"
	"encoding/json"
	"net/http"
	"runtime"
	"strings"

	"github.com/rakyll/statik/fs"
)

const (
	_ = iota
	jsonDecodeError
)

var ErrorHttpCode = map[int]string{
	jsonDecodeError: "json decode 错误",
}

// 静态资源处理
func (a *authority) assets(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.SplitAfterN(r.URL.Path, a.urlPath, 2)
	if len(urlPath) != 2 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("urlPath != 2"))
		return
	}
	staticsData, err := a.assetsReadFile("/" + urlPath[1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if len(r.URL.Path) >= 4 {
		suffix := strings.SplitAfterN(r.URL.Path[len(r.URL.Path)-4:], ".", 2)
		if len(suffix) == 2 {
			switch suffix[1] {
			case "css":
				w.Header().Add("content-type", "text/css")
			case "gif":
				w.Header().Add("content-type", "image/gif")
			case "ico":
				w.Header().Add("content-type", "image/x-icon")
			case "js":
				w.Header().Add("content-type", "application/javascript")
			case "json":
				w.Header().Add("content-type", "application/javascript")
			}
		}
	}
	w.Write(staticsData)
}

func (a *authority) assetsReadFile(name string) ([]byte, error) {
	if runtime.GOOS == "windows" {
		return fs.ReadFile(a.fileSystem, strings.ReplaceAll(name, "/", "\\"))
	}
	return fs.ReadFile(a.fileSystem, name)
}

func (a *authority) htmlHandle(w http.ResponseWriter, r *http.Request) {
	originContext, err := a.assetsReadFile("/index.html")
	if err != nil {
		originContext = []byte("fileSystem is error")
	}
	htmlContext := bytes.ReplaceAll(originContext, []byte("./"), []byte(r.URL.Path+"/"))
	w.Write(htmlContext)
}

func (a *authority) apiHandle(w http.ResponseWriter, r *http.Request) {
	WSuccess(w, "apiHandle")
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func wResponse(w http.ResponseWriter, resp Response) {
	jsonResp, err := json.Marshal(resp)
	w.Header().Add("content-type", "application/json")
	if err != nil {
		w.Write([]byte("{\"code\":0, \"msg\":\"json转码错误\"}"))
		return
	}
	w.Write(jsonResp)
}

func WSuccess(w http.ResponseWriter, data interface{}) {
	response := Response{
		Msg:  "操作成功",
		Data: data,
	}
	wResponse(w, response)
}

func WFailure(w http.ResponseWriter, code int) {
	response := Response{
		Code: -1,
		Msg:  "未定义的错误",
		Data: map[string]interface{}{},
	}
	if _, ok := ErrorHttpCode[code]; ok {
		response.Code = code
		response.Msg = ErrorHttpCode[code]
	}
	wResponse(w, response)
}
