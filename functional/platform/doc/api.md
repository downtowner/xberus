# xca平台信息组件

## 应用场景

- 通过webUI查询加载了哪些插件，注册了哪些组件，创建了哪些命名对象
- 通过webUI查询通过插件方式注册了那些组件，每个组件下分别有哪些命名对象
- 通过webUI查询通过代码方式注册了那些组件，每个组件下分别有哪些命名对象

## 注意事项

- 必须调用xca.CreateNamedObject(cmpName, objName string)方法才能记录组件和命名对象的关系


## 接口

```golang
type PlatformInfoManager interface {
	// httpHandler
	HttpHandler(w http.ResponseWriter, r *http.Request)
	// iris路由包装
	IrisWrapRouter(app *iris.Application, path string)
}
```

### HttpHandler

`HttpHandler(w http.ResponseWriter, r *http.Request)`

作用:http handler,这里实现的就是标准库net/http的Handler,所有的web框架都可以基于当前handler做二次封装


### IrisWrapRouter

`IrisWrapRouter(app *iris.Application, path string)`

作用:iris包装HttpHandler

#### 传入参数:

app:iris v12 的对象，也就是iris.New()

path:绑定的路由

#### 返回参数


## 待完善

- 目前只做了iris路由包装

- 目前必须调用xca.CreateNamedObject(cmpName, objName string)方法才能记录组件和命名对象的关系

欢迎大家交流使用,如有问题，请大家及时指正，联系邮箱 [zousf@novastar.tech](mailto:zousf@novastar.tech)