package platform

import (
	"fmt"
	"log"
	"testing"

	"git.vnnox.net/ncp/xframe/xca"
)

var (
	testObjName       = "testObjName"
	testComponentName = "authManager.uauth.novastar.tech"
	testPluginPath    = "./test.so"
)

func TestGetPlatformInfo(t *testing.T) {
	createObj()
	loadPlugin()
	pluginList, componentList := GetPlatformInfo()
	checkPluginPath := false
	checkPluginComponentName := false
	checkPluginObjName := false
	for _, pValue := range pluginList {
		//so在go version go1.15.3 linux/amd64下编译
		if pValue.Path == testPluginPath {
			checkPluginPath = true
		}
		for _, cValue := range pValue.Component {
			if cValue.Name == testComponentName {
				checkPluginComponentName = true
			}
			if len(cValue.ObjName) == 1 && cValue.ObjName[0] == testObjName {
				checkPluginObjName = true
			}
		}
	}
	if !checkPluginPath {
		t.Error("plugin import mode Path error")
	}
	if !checkPluginComponentName {
		t.Error("plugin import mode ComponentName error")
	}
	if !checkPluginObjName {
		t.Error("plugin import mode ObjName error")
	}

	checkComponentName := false
	checkComponentObjName := false
	for _, cValue := range componentList {
		if cValue.Name == ComponentName {
			checkComponentName = true
		}
		for _, oValue := range cValue.ObjName {
			if oValue == testObjName {
				checkComponentObjName = true
			}
		}
	}
	if !checkComponentName {
		t.Error("code import mode ComponentName error")
	}
	if !checkComponentObjName {
		t.Error("code import mode ObjName error")
	}
}

func loadPlugin() {
	err := xca.LoadModule(testPluginPath)
	if err != nil {
		log.Fatal(err)
	}
	xca.CreateNamedObject(testComponentName, testObjName)
}

func createObj() (PlatformInfoManager, error) {
	//  本地内置信息
	var localComponents = []map[string]interface{}{
		PlatformInfoManagerComponentInfo,
	}
	// 注册内置组件
	err := xca.RegisterComponents(localComponents)
	if err != nil {
		return nil, err
	}
	obj, ok := xca.CreateNamedObject(ComponentName, testObjName).(PlatformInfoManager)
	if !ok {
		return nil, fmt.Errorf("defined platformInfo error")
	}
	return obj, nil
}
