package usysinfo

import (
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func TestIrisWrapRouter(t *testing.T) {
	obj := createObj()
	obj.NewStorage()
	// obj.CrontabWrite(obj.NewStorage(), 7)

	app := iris.New()
	url := "/sysinfo"
	obj.IrisWrapRouter(app, url)

	e := httptest.New(t, app)
	e.GET(url).Expect().Status(httptest.StatusOK)
}
