package corba

var ctrNode ICenterNode
var subNode ISubNode
var messageManager IMessageManager

func init() {
	ctrNode = NewCenterNode()
	subNode = NewSubNode()
	messageManager = &MessageManager{}
}

// RunCenterNode 运行中心节点
// 参数:
//     address 中心节点要监听的地址，如":8080"
func RunCenterNode(address string) {
	go func() {
		ctrNode.Run(address)
	}()
}

// // SendMessageForCenter 发送消息
// func SendMessageForCenter(nodeTag string, msg IMessage) error {
// 	return ctrNode.SendMessage(nodeTag, msg)
// }

// SetSubNodeTag 设置子节点的唯一标识
func SetSubNodeTag(tag string) error {
	return subNode.SetNodeTag(tag)
}

// RunSubNode 运行子节点
// 参数:
//     address 中心节点的通信地址，如"127.0.0.1:8080"
func RunSubNode(address string) {
	go func() {
		subNode.Run(address)
	}()

}

// RunSubNodeDefault 运行子节点，按默认节点名运行
// 参数:
//     address 中心节点的通信地址，如"127.0.0.1:8080"
func RunSubNodeDefault(address string) {

	// 分布式消息组件加载
	tag := CreateSubNodeTag()
	err := SetSubNodeTag(tag)
	if err != nil {
		panic(err)
	}

	go func() {
		subNode.Run(address)
	}()

}

// CreateSubNodeTag 生成一个子节点tag
func CreateSubNodeTag() string {
	return subNode.CreateTag()
}

// GetSubState 子节点的运行状态
func GetSubState() bool {
	return subNode.GetState()
}

// RegisterObject 注册对象到中心节点
func RegisterObject(obj string) error {
	return subNode.RegisterObject(obj)
}

// SendMessage 发送消息
func SendMessage(msg IMessage) (interface{}, error) {
	return subNode.SendMessage(msg)
}
