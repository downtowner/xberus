package uauthority

import (
	"log"
	"math/rand"
	"testing"
	"time"

	"git.vnnox.net/ncp/xframe/xca"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type TestAuthorityManager interface {
	// 注册url和权限的关系
	RegisterPermissions(permissionsMap map[string]string)
	// 修改url和权限的关系
	UpdatePermissions(url string, permissions string)
	// 注册角色和权限的关系
	RegisterRole(roleMap map[string][]string)
	// 修改角色和权限标识的关系
	UpdateRole(role string, permissions []string)
	// 权限认证
	PermissionsCheck(url string, userRole []string) bool
	// 获取角色列表
	GetRoleMap() map[string][]string
	// 获取权限列表
	GetPermissionsMap() map[string]string
}

func createObj() TestAuthorityManager {
	//	本地内置信息
	var localComponents = []map[string]interface{}{
		AuthorityManagerComponentInfo,
	}
	// 注册内置组件
	err := xca.RegisterComponents(localComponents)
	if err != nil {
		return nil
	}
	obj, ok := xca.CreateNamedObject(ComponentName, ComponentName).(TestAuthorityManager)
	if !ok {
		return nil
	}
	return obj
}

func TestPermissionsCheck(t *testing.T) {
	obj := createObj()
	if obj == nil {
		t.Fatal("create obj error")
	}
	obj.RegisterPermissions(getPermissionsByDB())
	obj.RegisterRole(getRoleByDB())
	// administrator角色的用户是拥有所有权限的
	testURL := "/system/user"
	userRole := []string{"administrator"}
	if !obj.PermissionsCheck(testURL, userRole) {
		t.Error("权限校验失败")
	}
	// developer角色的用户是没有/system/user的权限的
	testURL = "/system/user"
	userRole = []string{"developer"}
	if obj.PermissionsCheck(testURL, userRole) {
		t.Error("权限校验失败")
	}

	testUpdateURL := "/test/user"
	testRole := "test:user:view"
	role := "administrator"
	userRole = []string{role}
	emptyRole := []string{}

	obj.UpdatePermissions(testUpdateURL, testRole)
	obj.UpdateRole(role, []string{testRole})
	if !obj.PermissionsCheck(testUpdateURL, userRole) {
		t.Error("权限校验失败")
	}
	if obj.PermissionsCheck(testUpdateURL, emptyRole) {
		t.Error("权限校验失败")
	}

}

//模拟数据库中查询出来的结果
func getRoleByDB() map[string][]string {
	// #表结构
	// CREATE TABLE `role` (
	// 	`id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '角色ID',
	// 	`role_name` varchar(30) NOT NULL COMMENT '角色名称',
	// 	`perms` varchar(100) NOT NULL COMMENT '角色权限字符串',
	// 	PRIMARY KEY (`id`)
	//   ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';
	// 1个角色可以对应多个权限
	// 这里的key是角色名  value是permissions表里的perms字段
	return map[string][]string{
		// 管理员角色
		"administrator": []string{"system:user:view", "system:config:view"},
		// 开发者角色
		"developer": []string{"system:config:view"},
	}
}

//模拟数据库中查询出来的结果
func getPermissionsByDB() map[string]string {
	// #表结构
	// CREATE TABLE `permissions` (
	// 	`id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
	// 	`url` varchar(200) DEFAULT '' COMMENT '请求地址',
	// 	`perms` varchar(100) DEFAULT NULL COMMENT '权限标识',
	// 	PRIMARY KEY (`id`)
	//   ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限表';
	// 1个url对应一个权限
	// 权限标识建议是这样的格式:将前面的/符号去掉，剩下的/变成:, 最后再加行为动作
	return map[string]string{
		"/system/user":   "system:user:view",
		"/system/config": "system:config:view",
	}
}

//测试单角色有1000个权限
func BenchmarkPermissionsCheck1000Role(b *testing.B) {
	obj := createObj()
	if obj == nil {
		log.Fatal("create obj error")
	}
	roleNum := 1000
	testURL := "/test"
	userRole := []string{"tester"}
	testRole := make([]string, roleNum)
	for i := 0; i < roleNum; i++ {
		testRole[i] = RandStringRunes(10)
	}
	role := map[string][]string{
		"tester": testRole,
	}
	obj.RegisterPermissions(map[string]string{testURL: "tester"})
	obj.RegisterRole(role)
	for n := 0; n < b.N; n++ {
		obj.PermissionsCheck(testURL, userRole)
	}
}

//测试用户有1000个角色 1个角色1000个权限  时间复杂度O(1000*1000)
func BenchmarkPermissionsCheckUserHas1000Role(b *testing.B) {
	obj := createObj()
	if obj == nil {
		log.Fatal("create obj error")
	}
	roleNum := 1000
	testURL := "/test"
	userRole := make([]string, roleNum)
	testRole := make([]string, roleNum)
	for i := 0; i < roleNum; i++ {
		testRole[i] = RandStringRunes(10)
		userRole[i] = RandStringRunes(10)
	}
	userRole[0] = "tester"
	role := map[string][]string{
		"tester": testRole,
	}
	obj.RegisterRole(role)
	obj.RegisterPermissions(map[string]string{testURL: "tester"})
	for n := 0; n < b.N; n++ {
		obj.PermissionsCheck(testURL, userRole)
	}
}

func RandStringRunes(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
