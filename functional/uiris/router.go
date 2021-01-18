// Copyright (C) 2020 The Nova Company Ltd.All rights reserved.

package uiris

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"git.vnnox.net/ncp/xframe/xca"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

// ComponentNameIrisRouterManager manager
const ComponentName = "routerManager.iris.novastar.tech"

var components = []map[string]interface{}{
	{
		"name":    ComponentName,
		"version": "1.0",
		"creator": func() interface{} { return newIrisRouter() },
	},
}

// IrisRouterManagerComponentInfo info
var IrisRouterManagerComponentInfo = map[string]interface{}{
	"name":       ComponentName,
	"version":    "1.0",
	"creator":    func() interface{} { return newIrisRouter() },
	"components": components,
}

type RouterConf struct {
	Method  string `json:"method"`
	UrlPath string `json:"urlPath"`
	// PluginPath string       `json:"pluginPath"`
	Component string       `json:"component"`
	ObjName   string       `json:"objName"`
	Handler   []string     `json:"handler"`
	Sub       []RouterConf `json:"sub"`
}

type IrisRouterManager interface {
	// 注册路由
	RegisterRouter([]RouterConf) (*iris.Application, error)
	//测试的handler
	TestLoginCheckHandle(iris.Context)
	//测试的handler
	TestHandle(iris.Context)
}

type irisRouter struct {
	IrisRouterManager
}

func GetComponents() map[string]interface{} {
	return IrisRouterManagerComponentInfo
}

func newIrisRouter() interface{} {
	c := new(irisRouter)
	return c
}

// 注册路由
func (i *irisRouter) RegisterRouter(conf []RouterConf) (app *iris.Application, err error) {
	app = iris.New()
	for _, parRouter := range conf {
		method := upperFirstChar(parRouter.Method)
		switch method {
		case "Get", "Post", "Put", "Delete", "Connect", "Head", "Options", "Patch", "Trace", "Any":
		default:
			panic(fmt.Errorf("did not support http method [%s]", method))
			// return nil, fmt.Errorf("Didn't support http method %s", method)
		}
		routerHandler, reflectHandler, err := i.getHandler(parRouter.Component, parRouter.ObjName, parRouter.Handler)
		if err != nil {
			return nil, err
		}
		if len(parRouter.Sub) > 0 {
			app.PartyFunc(parRouter.UrlPath, func(party iris.Party) {
				party.Use(routerHandler...)
				for _, subRouter := range parRouter.Sub {
					routerHandler, reflectHandler, err = i.getHandler(subRouter.Component, subRouter.ObjName, subRouter.Handler)
					if err != nil {
						return
					}
					reflectMethod := reflect.ValueOf(party).MethodByName(method)
					reflectMethod.Call(append([]reflect.Value{reflect.ValueOf(subRouter.UrlPath)}, reflectHandler...))
				}
			})
			continue
		}
		reflectMethod := reflect.ValueOf(app).MethodByName(method)
		reflectMethod.Call(append([]reflect.Value{reflect.ValueOf(parRouter.UrlPath)}, reflectHandler...))
	}
	return
}

func (i *irisRouter) getHandler(component, objectName string, handler []string) (contextHandler []context.Handler, reflectValue []reflect.Value, err error) {
	obj := xca.FindObject(objectName)
	if obj == nil {
		obj = xca.CreateNamedObject(component, objectName)
		if obj == nil {
			// panic(fmt.Errorf("Didn't component %s findObject %s", component, objectName))
			return nil, nil, fmt.Errorf("component [%s] didnot findObject [%s]", component, objectName)
		}
	}
	valueV := reflect.ValueOf(obj)
	valueT := reflect.TypeOf(obj)
	for _, pHandler := range handler {
		if _, ok := valueT.MethodByName(pHandler); !ok {
			return nil, nil, fmt.Errorf("component [%s]   objectName [%s]   methodName [%s] is not exist", component, objectName, pHandler)
		}
		var (
			handlerFunc func(ctx iris.Context)
			ok          bool
		)
		if handlerFunc, ok = valueV.MethodByName(pHandler).Interface().(func(ctx iris.Context)); !ok {
			return nil, nil, fmt.Errorf("cannot use component [%s]   objectName [%s]   methodName [%s] interface as func(ctx iris.Context)", component, objectName, pHandler)
		}
		contextHandler = append(contextHandler, handlerFunc)
		reflectValue = append(reflectValue, reflect.ValueOf(handlerFunc))
	}
	return contextHandler, reflectValue, nil
}

func (i *irisRouter) TestLoginCheckHandle(ctx iris.Context) {
	ctx.HTML("<h1>TestLoginCheckHandle</h1>")
	ctx.Next()
}

func (i *irisRouter) TestHandle(ctx iris.Context) {
	ctx.HTML(fmt.Sprintf("<h1> TestHandle %d</h1>", time.Now().UnixNano()))
}

func upperFirstChar(str string) string {
	if len(str) > 0 {
		str = strings.ToLower(str)
		runelist := []rune(str)
		if int(runelist[0]) >= 97 && int(runelist[0]) <= 122 {
			runelist[0] = rune(int(runelist[0]) - 32)
			str = string(runelist)
		}
	}
	return str
}
