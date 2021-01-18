package platform

import (
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func TestIrisWrapRouter(t *testing.T) {
	obj, err := createObj()
	if err != nil {
		t.Fatal(err)
	}
	app := iris.New()
	url := "/test"
	obj.IrisWrapRouter(app, url)

	e := httptest.New(t, app)
	e.GET(url).Expect().Status(httptest.StatusOK)
}
