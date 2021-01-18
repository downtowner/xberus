# 简要说明

提供强大的log功能，暂时支持本地和控制台，后期考虑假如远程或者邮件等

# 详细说明

提供强大的log功能，暂时支持本地和控制台，后期考虑假如远程或者邮件等，下面对使用的方法做一个简单的说明

### 日志等级

支持八个日志等级，分别为：

```go
LevelEmergency = iota // 紧急级别
LevelAlert            // 报警级别
LevelCritical         // 严重错误级别
LevelError            // 错误级别
LevelWarning          // 警告级别
LevelNotice           // 注意级别
LevelInformational    // 报告级别
LevelDebug            // 除错级别
```

### 主要接口

```
SetLogger("file",`{"filename":"logs/error.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
```

设置logger的属性，意义分别如下：

```
filename 保存的文件名
maxlines 每个文件保存的最大行数，默认值 1000000
maxsize  每个文件保存的最大尺寸，默认值是 1 << 28, //256 MB
daily    是否按照每天 logrotate，默认是 true
maxdays  文件最多保存多少天，默认保存 7 天
rotate   是否开启 logrotate，默认是 true
level    日志保存的时候的级别，默认是 Trace 级别
perm     日志文件权限
```

不设置表示使用默认配置

```go
EnableFuncCallDepth(true) 
```

表示打印当前log所在文件和所在行

```
DelLogger("file")
```

删除输出类型，到控制台传入`console`，到日志传入`file`

```
SetLevel(0) 
```

设置日志等级，0表示只输出紧急，若配置为7，表示全部输出，若配置1只输出紧急和报警级别



使用实例参看eg.md文档



# 计划

后续计划可能扩展远程log或者邮件订阅log



欢迎大家交流使用,如有问题，请大家及时指正，联系邮箱 [xiaobing@novastar.tech](mailto:moubo@novastar.tech)