package usysinfo

import (
	"net/http"
	"strings"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func TestIrisHandler(t *testing.T) {
	obj := createObj()
	obj.NewStorage()
	// obj.CrontabWrite(obj.NewStorage(), 7)

	app := iris.New()
	url := "/sysinfo"
	//iris路由包装
	app.WrapRouter(func(w http.ResponseWriter, r *http.Request, router http.HandlerFunc) {
		path := r.URL.Path
		if strings.Contains(path, url) {
			obj.HttpHandler(w, r)
			return
		}
		router(w, r)
	})
	e := httptest.New(t, app)
	e.GET(url).Expect().Status(httptest.StatusOK)
}
