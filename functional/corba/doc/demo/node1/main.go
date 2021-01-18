package main

import (
	"fmt"
	"time"

	"git.vnnox.net/ncp/xframe/functional/corba"
	"git.vnnox.net/ncp/xframe/xca"
)

// PluginA 组件A
type PluginA struct {
	Name  string
	Age   int
	Score int
}

// GetInfo 获得组件信息
func (c *PluginA) GetInfo() string {
	return fmt.Sprint("PluginA:{", "Name: ", c.Name, " Age: ", c.Age, " Score: ", c.Score, "}")
}

// HandleMessage 消息handle
func (c *PluginA) HandleMessage(msg xca.Message) (interface{}, error) {
	fmt.Println("A say ! ")
	return "success", nil
}

var localComponents = []map[string]interface{}{
	{"name": "PluginA.driver.mocar.autopard.com", "catalog": "service", "version": "1.0", "creator": func() interface{} { return new(PluginA) }},
}

func createObject() {
	xca.CreateNamedObject("PluginA.driver.mocar.autopard.com", "plugina")
}
func main() {

	// XCA
	xca.RegisterComponents(localComponents)
	createObject()

	// 运行消息的子节点
	corba.RunSubNodeDefault("127.0.0.1:8080")

	// 注册对象
	corba.RegisterObject("plugina")

	// 等候子节点监听运行成功，您也可以通过corba.GetState()主动获得运行状态
	time.Sleep(3 * time.Second)

	// 发消息
	msg := corba.Message{}
	msg.Code = corba.SystemSendMessage
	msg.ObjectName = "pluginb"

	simpleMsg := xca.SimpleMessage{}
	simpleMsg.Data = "sssss"
	msg.Data = simpleMsg

	_, err := corba.SendMessage(&msg)
	if err != nil {
		fmt.Println(err)
	}

	// 等候 下一轮的消息
	select {}

}
