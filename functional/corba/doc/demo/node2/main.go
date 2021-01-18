package main

import (
	"fmt"

	"git.vnnox.net/ncp/xframe/functional/corba"
	"git.vnnox.net/ncp/xframe/xca"
)

// PluginB 组件B
type PluginB struct {
	Name  string
	Age   int
	Score int
}

// GetInfo 组件的信息
func (c *PluginB) GetInfo() string {
	return fmt.Sprint("PluginB:{", "Name: ", c.Name, " Age: ", c.Age, " Score: ", c.Score, "}")
}

// HandleMessage 消息handle
func (c *PluginB) HandleMessage(msg xca.Message) (interface{}, error) {
	fmt.Println("B say ! ")
	fmt.Println(msg)
	return "success", nil
}

var localComponents = []map[string]interface{}{
	{"name": "PluginB.driver.mocar.autopard.com", "catalog": "service", "version": "1.0", "creator": func() interface{} { return new(PluginB) }},
}

func createObject() {
	xca.CreateNamedObject("PluginB.driver.mocar.autopard.com", "pluginb")
}
func main() {

	// XCA
	xca.RegisterComponents(localComponents)
	createObject()

	// 运行消息的子节点
	corba.RunSubNodeDefault("127.0.0.1:8080")

	// 注册对象
	corba.RegisterObject("pluginb")

	// 等候 下一轮的消息
	select {}

}
