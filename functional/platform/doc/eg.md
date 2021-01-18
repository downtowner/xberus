# xca平台信息

## 注意事项

- 必须调用xca.CreateNamedObject(cmpName, objName string)方法才能记录组件和命名对象的关系

## 使用示例

自定义注册路由，然后通过a标签或者iframe的方式展示页面即可

```golang
package main

import (
	"fmt"
	"log"
	"net/http"

	"git.vnnox.net/ncp/xframe/functional/platform"
	"git.vnnox.net/ncp/xframe/xca"
	"github.com/kataras/iris/v12"
)

type PlatformInfoManager interface {
	// httpHandler
	HttpHandler(w http.ResponseWriter, r *http.Request)
	// iris路由包装
	IrisWrapRouter(app *iris.Application, path string)
}

func main() {
	log.SetFlags(log.Llongfile)
	obj, err := createObj()
	if err != nil {
		log.Fatal(err)
	}
	TestLoadPlugin()
	app := iris.New()
	//注册路由
	url := "/xca/info"
	obj.IrisWrapRouter(app, url)
	app.Run(iris.Addr(":8080"))
}

func createObj() (PlatformInfoManager, error) {
	//  本地内置信息
	var localComponents = []map[string]interface{}{
		platform.PlatformInfoManagerComponentInfo,
	}
	// 注册内置组件
	err := xca.RegisterComponents(localComponents)
	if err != nil {
		return nil, err
	}
	obj, ok := xca.CreateNamedObject(platform.ComponentName, "objName").(PlatformInfoManager)
	if !ok {
		return nil, fmt.Errorf("defined platformInfo error")
	}
	return obj, nil
}

func TestLoadPlugin() {
	err := xca.LoadModule("./uauthority.so")
	if err != nil {
		log.Println(err)
	}
	xca.CreateNamedObject("authManager.uauth.novastar.tech", "obj.uauth")
}


```

curl请求接口

```bash
curl http://localhost:8080/xca/info?api=1
#返回
{
    "code":0,
    "msg":"操作成功",
    "data":{
        "plugin":[//加载了哪些插件
            {
                "name":"./uauthority.so",//插件路径
                "component":[//插件注册了哪些组件
                    {
                        "name":"authManager.uauth.novastar.tech",//组件名
                        "objName":[//对象名称列表
                            "obj.uauth"
                        ]
                    }
                ]
            }
        ],
        "component":[//通过代码的方式注册了哪些组件
            {
                "name":"objectManager.xca.novastar.tech",//组件名
                "objName":[

                ]
            },
            {
                "name":"platformInfo.xca.novastar.tech",//组件名
                "objName":[//对象名称列表
                    "objName"
                ]
            }
        ]
    }
}
```

## 待完善


- 目前只做了iris路由包装

- 目前必须调用xca.CreateNamedObject(cmpName, objName string)方法才能记录组件和命名对象的关系

欢迎大家交流使用,如有问题，请大家及时指正，联系邮箱 [zousf@novastar.tech](mailto:zousf@novastar.tech)