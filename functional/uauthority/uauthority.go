// Copyright (C) 2020 The Nova Company Ltd.All rights reserved.

package uauthority

import (
	"net/http"
	"regexp"
	"sync"

	"github.com/kataras/iris/v12"
)

// ComponentName manager
const ComponentName = "authManager.uauth.novastar.tech"

var components = []map[string]interface{}{
	{
		"name":    ComponentName,
		"version": "1.0",
		"creator": func() interface{} { return newAuth() },
	},
}

// AuthorityManagerComponentInfo info
var AuthorityManagerComponentInfo = map[string]interface{}{
	"name":       ComponentName,
	"version":    "1.0",
	"creator":    func() interface{} { return newAuth() },
	"components": components,
}

type AuthorityManager interface {
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

type authority struct {
	AuthorityManager
	permissionsMap permissionsMap
	roleMap        roleMap
	// 后台静态文件
	fileSystem http.FileSystem
	// 后台url路由绑定路径
	urlPath string
}

type permissionsMap struct {
	Locker      sync.RWMutex
	Permissions map[string]string
}

type roleMap struct {
	Locker sync.RWMutex
	Role   map[string][]string
}

func GetComponents() map[string]interface{} {
	return AuthorityManagerComponentInfo
}

func newAuth() interface{} {
	c := &authority{
		permissionsMap: permissionsMap{},
		roleMap:        roleMap{},
	}
	return c
}

// 注册url和权限的关系
// permissionsMap map[string]string{"url":"权限标识"}
func (a *authority) RegisterPermissions(permissionsMap map[string]string) {
	a.permissionsMap.Locker.Lock()
	defer a.permissionsMap.Locker.Unlock()
	a.permissionsMap.Permissions = permissionsMap
}

// 更新url和权限的关系 这里是替换，不是追加
func (a *authority) UpdatePermissions(url string, permissions string) {
	a.permissionsMap.Locker.Lock()
	defer a.permissionsMap.Locker.Unlock()
	a.permissionsMap.Permissions[url] = permissions
}

// 注册角色和权限的关系
// roleMap map[string][]string{"admin":[]string{"权限标识1","权限标识2"}}
func (a *authority) RegisterRole(roleMap map[string][]string) {
	a.roleMap.Locker.Lock()
	defer a.roleMap.Locker.Unlock()
	a.roleMap.Role = roleMap
}

// 更新角色和权限的关系
// 这里是替换，不是追加
func (a *authority) UpdateRole(role string, permissions []string) {
	a.roleMap.Locker.Lock()
	defer a.roleMap.Locker.Unlock()
	a.roleMap.Role[role] = permissions
}

// 权限校验
func (a *authority) PermissionsCheck(url string, userRole []string) bool {
	var permissions string
	var ok bool
	a.permissionsMap.Locker.RLock()
	defer a.permissionsMap.Locker.RUnlock()
	if permissions, ok = a.permissionsMap.Permissions[url]; !ok {
		return true
	}
	rolePermissions := map[string]struct{}{}
	a.roleMap.Locker.RLock()
	defer a.roleMap.Locker.RUnlock()
	for _, role := range userRole {
		for _, tempPermissions := range a.roleMap.Role[role] {
			rolePermissions[tempPermissions] = struct{}{}
		}
	}
	if _, ok = rolePermissions[permissions]; !ok {
		return false
	}
	return true
}

// 获取角色列表
func (a *authority) GetRoleMap() map[string][]string {
	a.roleMap.Locker.RLock()
	defer a.roleMap.Locker.RUnlock()
	return a.roleMap.Role
}

// 获取权限列表
func (a *authority) GetPermissionsMap() map[string]string {
	a.permissionsMap.Locker.RLock()
	defer a.permissionsMap.Locker.RUnlock()
	return a.permissionsMap.Permissions
}

// 管理后台 现在先不暴露 不知道是否需要前端需求
func (a *authority) managementHandler(w http.ResponseWriter, r *http.Request) {
	// 静态资源处理
	if regexp.MustCompile(`\.(js|css|json|ico|gif)+$`).MatchString(r.URL.Path) {
		a.assets(w, r)
		return
	}

	api, ok := r.URL.Query()["api"]
	if !ok || len(api) < 1 {
		// 页面
		if len(a.urlPath) == 0 {
			a.urlPath = r.URL.Path
		}
		a.htmlHandle(w, r)
		return
	}
	// api
	a.apiHandle(w, r)
	return
}
