
## 建表示例

```sql
# URL权限表 1对1
CREATE TABLE `permissions` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `url` varchar(200) DEFAULT '' COMMENT '请求地址',
  `perms` varchar(100) DEFAULT NULL COMMENT '权限标识',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限表';

# 角色权限表 1角色对多权限标识
CREATE TABLE `role` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `role_name` varchar(30) NOT NULL COMMENT '角色名称',
  `perms` varchar(100) NOT NULL COMMENT '角色权限标识',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

```

## iris使用
```golang
package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"git.vnnox.net/ncp/xframe/functional/uauthority"
	"git.vnnox.net/ncp/xframe/xca"
	"github.com/kataras/iris/v12"
)

type TestAuthorityManager interface {
	// 注册url和权限的关系
	RegisterPermissions(permissionsMap map[string]string)
	// 修改url和权限的关系,这里是替换，不是追加
	UpdatePermissions(url string, permissions string)
	// 注册角色和权限的关系
	RegisterRole(roleMap map[string][]string)
	// 修改角色和权限标识的关系,这里是替换，不是追加
	UpdateRole(role string, permissions []string)
	// 权限认证
	PermissionsCheck(url string, userRole []string) bool
	// 获取角色列表 map[string][]string{"角色名字":[]string{"权限标识1","权限标识2"}}
	GetRoleMap() map[string][]string
	// 获取权限列表
	GetPermissionsMap() map[string]string
	// iris权限校验中间件 需要在前面的中间件中定义用户角色 ctx.Values().Set("userRole", []string{"用户角色"})
	IrisPermissionsCheck(ctx iris.Context)
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	//  本地内置信息
	var localComponents = []map[string]interface{}{
		uauthority.AuthorityManagerComponentInfo,
	}
	// 注册内置组件
	err := xca.RegisterComponents(localComponents)
	if err != nil {
		log.Println(err)
		return
	}
	authManager := xca.CreateNamedObject(uauthority.ComponentName, "auth")
	if _, ok := authManager.(TestAuthorityManager); !ok {
		log.Fatal("obj error")
		return
	}
	obj := authManager.(TestAuthorityManager)
	// 注册权限
	obj.RegisterPermissions(getPermissionsByDB())
	// 注册角色
	obj.RegisterRole(getRoleByDB())

	app := iris.New()
	loginCheck := func(ctx iris.Context) {
		// 这里一般都是根据用户token获取用户拥有哪些角色，然后将角色存入context中，由xca权限校验组件读取校验
		// 这里是模拟，所以使用了ctx.URLParam("role")去获取
		// userRole := []string{ctx.URLParam("role")}这里需要自己改成对应的获取方法
		ctx.Values().Set("userRole", userRole)
		ctx.Next()
	}
	// 使用xca权限校验组件
	system := app.Party("/system", loginCheck, obj.IrisPermissionsCheck)
	system.Get("/user", func(ctx iris.Context) {
		ctx.HTML(fmt.Sprintf("success %d", time.Now().UnixNano()))
	})
	system.Get("/config", func(ctx iris.Context) {
		ctx.HTML(fmt.Sprintf("success %d", time.Now().UnixNano()))
	})
	//模拟后台修改角色权限
	//curl http://localhost:8080/system/role/update?role=administrator -d "role=administrator&permissions=system:role:update"
	system.Post("/role/update", func(ctx iris.Context) {
		role := ctx.PostValue("role")
		permissions := ctx.PostValue("permissions")
		obj.UpdateRole(role, strings.Split(permissions, ","))
		ctx.HTML(fmt.Sprintf("success %d", time.Now().UnixNano()))
	})
	//模拟后台修改URL权限
	// curl http://localhost:8080/system/role/update?role=administrator -X POST -d "role=tester&permissions=system:config:view"
	system.Post("/permissions/update", func(ctx iris.Context) {
		url := ctx.PostValue("url")
		permissions := ctx.PostValue("permissions")
		obj.UpdatePermissions(url, permissions)
		ctx.HTML(fmt.Sprintf("success %d", time.Now().UnixNano()))
	})
	//查看角色
	system.Get("/role/list", func(ctx iris.Context) {
		ctx.JSON(obj.GetRoleMap())
	})
	//查看url和权限关系
	system.Get("/permissions/list", func(ctx iris.Context) {
		ctx.JSON(obj.GetPermissionsMap())
	})
	app.Listen(":8080")

}

//模拟数据库中查询出来的结果
func getRoleByDB() map[string][]string {
	// #表结构
	// CREATE TABLE `role` (
	//  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '角色ID',
	//  `role_name` varchar(30) NOT NULL COMMENT '角色名称',
	//  `perms` varchar(100) NOT NULL COMMENT '角色权限字符串',
	//  PRIMARY KEY (`id`)
	//   ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';
	// 1个角色可以对应多个权限
	// 这里的key是角色名  value是permissions表里的perms字段
	return map[string][]string{
		// 管理员角色
		"administrator": []string{"system:permissions:update", "system:role:update", "system:user:view", "system:config:view"},
		// 开发者角色
		"developer": []string{"system:config:view"},
		// 测试角色
		"tester": []string{},
	}
}

//模拟数据库中查询出来的结果
func getPermissionsByDB() map[string]string {
	// #表结构
	// CREATE TABLE `permissions` (
	//  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
	//  `url` varchar(200) DEFAULT '' COMMENT '请求地址',
	//  `perms` varchar(100) DEFAULT NULL COMMENT '权限标识',
	//  PRIMARY KEY (`id`)
	//   ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限表';
	// 1个url对应一个权限
	// 权限标识建议是这样的格式:将前面的/符号去掉，剩下的/变成:, 最后再加行为动作
	return map[string]string{
		"/system/user":               "system:user:view",
		"/system/config":             "system:config:view",
		"/system/permissions/update": "system:permissions:update",
		"/system/role/update":        "system:role:update",
	}
}
```

## 待完善


### 1.目前角色管理和url权限管理无web页面

### 2.目前接口权限校验没区分http Method

欢迎大家交流使用,如有问题，请大家及时指正，联系邮箱 [zousf@novastar.tech](mailto:zousf@novastar.tech)