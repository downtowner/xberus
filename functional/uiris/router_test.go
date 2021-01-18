package uiris

import (
	"encoding/json"
	"log"
	"testing"

	"git.vnnox.net/ncp/xframe/xca"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

type TestIrisRouterManager interface {
	// 注册路由
	RegisterRouter([]RouterConf) (*iris.Application, error)
}

func TestRegisterRouter(t *testing.T) {
	//	本地内置信息
	var localComponents = []map[string]interface{}{
		IrisRouterManagerComponentInfo,
	}
	// 注册内置组件
	err := xca.RegisterComponents(localComponents)
	if err != nil {
		log.Println(err)
		return
	}
	routerManager := xca.CreateNamedObject("routerManager.iris.novastar.tech", "routerManager")
	//配置
	conf := `[
		{
			"method":"GET",
			"urlPath":"/user",
			"component":"routerManager.iris.novastar.tech",
			"objName":"routerManager",
			"handler":[
				"TestLoginCheckHandle"
			],
			"sub":[
				{
					"method":"GET",
					"urlPath":"/test",
					"component":"routerManager.iris.novastar.tech",
					"objName":"routerManager",
					"handler":[
						"TestHandle"
					]
				}
			]
		},
		{
			"method":"GET",
			"urlPath":"/test",
			"component":"routerManager.iris.novastar.tech",
			"objName":"routerManager",
			"handler":[
				"TestHandle"
			]
		}
	]`

	var routerConf []RouterConf
	if err := json.Unmarshal([]byte(conf), &routerConf); err != nil {
		t.Error(err)
	}
	app, err := routerManager.(TestIrisRouterManager).RegisterRouter(routerConf)
	if err != nil {
		t.Error(err)
	}
	// app.Listen(":8080")
	e := httptest.New(t, app)

	e.GET("/user/test").Expect().Status(httptest.StatusOK)
	e.GET("/test").Expect().Status(httptest.StatusOK)
	e.GET("/abc").Expect().Status(httptest.StatusNotFound)

}

func TestErrorRouter(t *testing.T) {
	//	本地内置信息
	var localComponents = []map[string]interface{}{
		IrisRouterManagerComponentInfo,
	}
	// 注册内置组件
	err := xca.RegisterComponents(localComponents)
	if err != nil {
		log.Println(err)
		return
	}
	routerManager := xca.CreateNamedObject("routerManager.iris.novastar.tech", "routerManager")
	//配置
	conf := `[
	{
		"method":"GET",
		"urlPath":"/test",
		"component":"routerManager.iris.novastar.tech",
		"objName":"routerManager",
		"handler":[
			"NotExistHandle"
		]
	}
]`

	var routerConf []RouterConf
	if err := json.Unmarshal([]byte(conf), &routerConf); err != nil {
		t.Error(err)
	}
	_, err = routerManager.(TestIrisRouterManager).RegisterRouter(routerConf)
	if err == nil {
		t.Error("NotExistHandle")
	}
}
