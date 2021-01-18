# uiris
## 需求场景

- 想一个路由对应一个组件的handler
- 想一个路由对应一个组件的一组handler
- 路由下的handler经常需要按需扩展（例如加权限校验或者登录授权之类），代码实现层不方便抽象分解
- 随着项目增加，代码层写死的路径可能会重复不方便排查
- router的代码随意遍布项目各个地方，不便查执行顺序也不便修改
- 想抽离之前写过的hander在另一个项目中使用，又得重复对接，而且没有统一标准，每次都要重新写

如果上面其中一点触碰到你的痛点，那建议你使用当前动态路由组件。

## 配置详情

配置可以是json可以是toml,只要最终解析出是对应的结构体即可
```golang
type RouterConf struct {
	Method    string       `json:"method"`
	UrlPath   string       `json:"urlPath"`
	Component string       `json:"component"`
	ObjName   string       `json:"objName"`
	Handler   []string     `json:"handler"`
	Sub       []RouterConf `json:"sub"`
}
```
示例的json配置
```json
[
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
]
```
### 参数说明

| 字段名        | 数据类型   |  字段说明  |
| --------   | :-----:  | :----:  |
| method      | string   |http method,可选项:Get,Post,Put,Delete,Connect,Head,Options,Patch,Trace,Any|
| urlPath        |   string   |   路由   |
| component        |    string    |  插件组件名  |
| objName        | string    |  实例化的对象名，名字一样，对象只会实例化一次  |
| handler        |    array    |  路由处理的handler,必须是组件里的方法,可以是多个，相当于app.Get(路由, handler)这里的handler的意思  |
| sub        |    array    |  如果数量大于0，则使用路由分组,相当于app.Party  |


## 使用示例

```golang
package main

import (
	"encoding/json"
	"log"
	"git.vnnox.net/ncp/xframe/functional/uiris"
	"git.vnnox.net/ncp/xframe/xca"

	"github.com/kataras/iris/v12"
)

type IrisRouterManager interface {
	// 注册路由
	RegisterRouter([]uiris.RouterConf) (*iris.Application, error)
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	//	本地内置信息
	var localComponents = []map[string]interface{}{
		uiris.IrisRouterManagerComponentInfo,
	}
	// 注册内置组件
	err := xca.RegisterComponents(localComponents)
	if err != nil {
		log.Println(err)
		return
	}
	routerManager := xca.CreateNamedObject(uiris.ComponentName, "routerManager")

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

	var routerConf []uiris.RouterConf
	if err := json.Unmarshal([]byte(conf), &routerConf); err != nil {
		log.Fatalf("router config decode error %s", err.Error())
		return
	}
	app, err := routerManager.(IrisRouterManager).RegisterRouter(routerConf)
	if err != nil {
		panic(err)
	}
	app.Listen(":8080")

}

```

## 待完善


欢迎大家交流使用,如有问题，请大家及时指正，联系邮箱 [zousf@novastar.tech](mailto:zousf@novastar.tech)