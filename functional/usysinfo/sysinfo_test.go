package usysinfo

import (
	"git.vnnox.net/ncp/xframe/xca"
)

func createObj() SysinfoManager {
	//	本地内置信息
	var localComponents = []map[string]interface{}{
		SysinfoManagerComponentInfo,
	}
	// 注册内置组件
	err := xca.RegisterComponents(localComponents)
	if err != nil {
		return nil
	}
	obj, ok := xca.CreateNamedObject(ComponentName, ComponentName).(SysinfoManager)
	if !ok {
		return nil
	}
	return obj
}
