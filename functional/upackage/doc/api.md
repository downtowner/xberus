# 简要说明

该包基于tcp传输层的协议对socket接收的数据进行解码，数据内容协议支持字节流、PB、json、xml等协议。

# 详细说明

该包是对tcp协议直接进行解码,解码格式如下:

|          | 模式段       | 标识长度 | 内容长度 | 内容             |
| -------- | ------------ | -------- | -------- | ---------------- |
| ID模式   | 1byte        | 2bytes   | 4bytes   | 理论内容最长32GB |
|          | Header长度7  |          |          |                  |
| 命令模式 | 1byte        | 8bytes   | 4bytes   | 理论内容最长32GB |
|          | header长度13 |          |          |                  |

该包支持两种协议模式，ID数字模式和CMD字符串模式,可以根据不同的场景需求进行配置

### ID模式

在ID模式下，协议的头部的总长度为7个字节，第一个字节表示该头部应该用那种模式进行解析，1表示ID模式，2表示命令模式，接下的二个字节表示消息标识，最后的四个字节表示内容长度，最大支持32GB。

### CMD模式

在cmd模式下，协议的头部总长度为13个字节，第一个字节表示该头部应该用那种模式进行解析，2表示cmd模式，接下的8个字节表示消息的标识，最后四个字节表示内容的长度，最大支付32GB



# 用法

### 编码过程

```go
//创建对象
p :=NewBPackage{}
//添加ID
p.AddPackageID(12)
//添加数据
p.AddUint8(3)
...
//添加完成标志
p.Done()
//获取数据
p.GetData()

//
```

### 解码过程

```go
//创建对象
p :=NewBPackage{}
//把数据导入该包中
p.AddBytes(...)
//读取标识
p.ReadPackageID()/p.ReadPackageCmd()
//读取数据长度
p.ReadDataLength()
//读取包体内容
p.Read...
```

注意： 解码的顺序必须和编码顺序一致。



欢迎大家交流使用,如有问题请大家及时指正，联系邮箱 [xiaobing@novastar.tech](mailto:moubo@novastar.tech)




