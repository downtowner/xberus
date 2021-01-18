// Copyright (C) 2020 The Nova Company Ltd.All rights reserved.

package platform

import (
	"encoding/json"
	"net/http"

	"git.vnnox.net/ncp/xframe/xca"
	"github.com/kataras/iris/v12"
)

// ComponentNameManager manager
const ComponentName = "platformInfo.xca.novastar.tech"

var components = []map[string]interface{}{
	{
		"name":    ComponentName,
		"version": "1.0",
		"creator": func() interface{} { return newPlatformInfo() },
	},
}

// AuthorityManagerComponentInfo info
var PlatformInfoManagerComponentInfo = map[string]interface{}{
	"name":       ComponentName,
	"version":    "1.0",
	"creator":    func() interface{} { return newPlatformInfo() },
	"components": components,
}

type PlatformInfoManager interface {
	// httpHandler
	HttpHandler(w http.ResponseWriter, r *http.Request)
	// iris路由包装
	IrisWrapRouter(app *iris.Application, path string)
}

type platformInfo struct {
	PlatformInfoManager
}

type Component struct {
	Name    string   `json:"name"`
	ObjName []string `json:"objName"`
}
type Plugin struct {
	Path      string      `json:"name"`
	Component []Component `json:"component"`
}

func GetComponents() map[string]interface{} {
	return PlatformInfoManagerComponentInfo
}

func newPlatformInfo() interface{} {
	c := &platformInfo{}
	return c
}

// httpHandler 目前iris的handler是基于当前方法封装
// 如果使用其他web框架也可以自己封装
func (p *platformInfo) HttpHandler(w http.ResponseWriter, r *http.Request) {
	api, ok := r.URL.Query()["api"]
	if !ok || len(api) < 1 {
		p.htmlHandle(w, r)
		return
	}
	p.apiHandle(w, r)
	return
}

func (p *platformInfo) htmlHandle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>test</h1>"))
}

func (p *platformInfo) apiHandle(w http.ResponseWriter, r *http.Request) {
	type ResponseResult struct {
		Plugin    []Plugin    `json:"plugin"`
		Component []Component `json:"component"`
	}
	pluginList, componentList := GetPlatformInfo()
	responseResult := ResponseResult{
		Plugin:    pluginList,
		Component: componentList,
	}
	// {"code":0,"msg":"操作成功","data":{"plugin":[{"name":"./uauthority.so","component":[{"name":"authManager.uauth.novastar.tech","objName":["obj.uauth"]}]}],"component":[{"name":"objectManager.xca.novastar.tech","objName":[]},{"name":"platformInfo.xca.novastar.tech","objName":["obj123"]}]}}
	WSuccess(w, responseResult)
}

func GetPlatformInfo() ([]Plugin, []Component) {
	pluginNames := xca.GetPluginNames()
	//用于过滤插件注册的组件
	pluginComponent := map[string]struct{}{}

	pluginList := make([]Plugin, 0)
	for _, path := range pluginNames {
		componentNames := xca.FindComponentNameByPluginName(path)
		componentList := make([]Component, 0)
		for _, cName := range componentNames {
			pluginComponent[cName] = struct{}{}
			componentInfo := Component{}
			componentInfo.Name = cName
			componentInfo.ObjName = xca.FindObjNameByComponentName(cName)
			componentList = append(componentList, componentInfo)
		}
		pluginInfo := Plugin{
			Path:      path,
			Component: componentList,
		}
		pluginList = append(pluginList, pluginInfo)
	}
	componentNames := xca.GetComponentNames()
	componentList := make([]Component, 0)

	for _, cName := range componentNames {
		if _, ok := pluginComponent[cName]; ok {
			continue
		}
		componentInfo := Component{}
		componentInfo.Name = cName
		componentInfo.ObjName = xca.FindObjNameByComponentName(cName)
		componentList = append(componentList, componentInfo)
	}
	return pluginList, componentList
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func wResponse(w http.ResponseWriter, resp Response) {
	jsonResp, err := json.Marshal(resp)
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
