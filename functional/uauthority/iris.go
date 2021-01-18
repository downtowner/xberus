package uauthority

import (
	"log"
	"net/http"

	"github.com/kataras/iris/v12"
)

// iris中间件
// 需要在前面中间件 定义userRole
// ctx.Values().Set("userRole", []string{"用户角色"})
func (a *authority) IrisPermissionsCheck(ctx iris.Context) {
	path := ctx.Path()
	userRole := ctx.Values().Get("userRole")
	if _, ok := userRole.([]string); !ok {
		errorText := `xca tips: userRole Type must []string. You should use like this: ctx.Values().Set("userRole", []string{"user role"})`
		log.Println(errorText)
		ctx.StatusCode(http.StatusForbidden)
		ctx.Write([]byte(errorText))
		return
	}
	if !a.PermissionsCheck(path, userRole.([]string)) {
		ctx.StatusCode(http.StatusForbidden)
		return
	}
	ctx.Next()
	return
}
