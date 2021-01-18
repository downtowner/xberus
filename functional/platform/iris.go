package platform

import (
	"net/http"
	"strings"

	"github.com/kataras/iris/v12"
)

// iris路由包装
// 其实也是基于HttpHandler的包装，其他web框架可以模仿然后自己包装一下
func (p *platformInfo) IrisWrapRouter(app *iris.Application, path string) {
	app.WrapRouter(func(w http.ResponseWriter, r *http.Request, router http.HandlerFunc) {
		if strings.Contains(r.URL.Path, path) {
			p.HttpHandler(w, r)
			return
		}
		router(w, r)
	})
}
