package handle

import (
	"strings"

	"github.com/gohutool/boot4go-docker-ui/db"
	. "github.com/gohutool/boot4go-docker-ui/log"
	. "github.com/gohutool/boot4go-docker-ui/model"
	. "github.com/gohutool/boot4go-util"
	httputil "github.com/gohutool/boot4go-util/http"
	. "github.com/gohutool/boot4go-util/jwt"
	routing "github.com/qiangxue/fasthttp-routing"
)

/*
*
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : user.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/13 15:15
* 修改历史 : 1. [2022/5/13 15:15] 创建文件 by LongYong
*/
type userHandler struct {
}

var UserHandler = &userHandler{}

func (u *userHandler) InitRouter(router *routing.Router, routerGroup *routing.RouteGroup) {
	router.Post("/login", u.Login)
	routerGroup.Post("/user/pwd", u.Pwd)
	routerGroup.Get("/user/list", u.ListUsers)
	routerGroup.Post("/user/create", u.CreateUser)
	routerGroup.Post("/user/delete", u.DeleteUser)
	routerGroup.Post("/user/resetpwd", u.ResetPwd)
	router.Get("/logout", u.Logout)
}

func (u *userHandler) Login(context *routing.Context) error {
	username := httputil.GetParams(context.RequestCtx, "username", "")
	password := httputil.GetParams(context.RequestCtx, "password", "")

	if username == "" || password == "" {
		httputil.Result.Fail("请填写登录用户名和用户密码").Response(context.RequestCtx)
		return nil
	}

	username = MD5(username)

	user := db.GetUser(username)

	if user == nil {
		httputil.Result.Fail("登录用户名和用户密码不正确").Response(context.RequestCtx)
		return nil
	}

	if user.UserPassword != SaltMd5(password, user.UserID) {
		httputil.Result.Fail("登录用户名和用户密码不正确").Response(context.RequestCtx)
		return nil
	}

	token := GenToken(user.UserID, Issuer, Issuer, TokenExpire)

	rtn := make(map[string]string)
	rtn["token"] = token
	rtn["userid"] = user.UserID

	httputil.Result.Success(rtn, "OK").Response(context.RequestCtx)

	Logger.Debug("%v %v %v", user.UserID, username, password)

	return nil

}

func (u *userHandler) Logout(context *routing.Context) error {
	Logger.Debug("%v", "Logout")
	httputil.Result.Success("", "OK").Response(context.RequestCtx)
	return nil
}

func (u *userHandler) ListUsers(context *routing.Context) error {
	users, err := db.ListUsers()
	if err != nil {
		httputil.Result.Fail("获取用户列表失败:" + err.Error()).Response(context.RequestCtx)
		return nil
	}

	// Return rows as-is; front-end grid will bind by field name
	httputil.Result.Success(users, "OK").Response(context.RequestCtx)
	return nil
}

func (u *userHandler) CreateUser(context *routing.Context) error {
	username := httputil.GetParams(context.RequestCtx, "username", "")
	password := httputil.GetParams(context.RequestCtx, "password", "")
	password2 := httputil.GetParams(context.RequestCtx, "password2", "")

	if username == "" || password == "" || password2 == "" {
		httputil.Result.Fail("请填写用户名和密码").Response(context.RequestCtx)
		return nil
	}

	if password != password2 {
		httputil.Result.Fail("两次输入的密码不一致").Response(context.RequestCtx)
		return nil
	}

	if len(password) < 6 || len(password) > 12 {
		httputil.Result.Fail("密码长度需为6-12位").Response(context.RequestCtx)
		return nil
	}

	userID := MD5(username)
	if db.GetUser(userID) != nil {
		httputil.Result.Fail("用户已存在").Response(context.RequestCtx)
		return nil
	}

	if err := db.CreateUser(username, password); err != nil {
		httputil.Result.Fail("创建用户失败:" + err.Error()).Response(context.RequestCtx)
		return nil
	}

	httputil.Result.Success("", "OK").Response(context.RequestCtx)
	return nil
}

func (u *userHandler) DeleteUser(context *routing.Context) error {
	userid := httputil.GetParams(context.RequestCtx, "userid", "")
	if userid == "" {
		httputil.Result.Fail("请选择要删除的用户").Response(context.RequestCtx)
		return nil
	}

	opUserID := ""
	opUserName := ""
	authHeader := string(context.Request.Header.Peek("Authorization"))
	if len(authHeader) > 0 {
		tokenStr := authHeader
		if strings.HasPrefix(strings.ToLower(tokenStr), "bearer ") {
			tokenStr = strings.TrimSpace(tokenStr[len("Bearer "):])
		}
		if len(tokenStr) > 0 {
			if sub, err := CheckToken(Issuer, tokenStr, func(subject string) (any, error) { return subject, nil }); err == nil {
				if v, ok := sub.(string); ok {
					opUserID = v
					if opUser := db.GetUser(opUserID); opUser != nil {
						opUserName = opUser.UserName
					}
				}
			}
		}
	}

	username := ""
	if user := db.GetUser(userid); user != nil {
		username = user.UserName
	}

	if err := db.DeleteUser(userid); err != nil {
		Logger.Critical("DeleteUser failed op_userid=%v op_username=%v userid=%v username=%v ip=%v err=%v", opUserID, opUserName, userid, username, context.RequestCtx.RemoteIP(), err)
		httputil.Result.Fail("删除用户失败:" + err.Error()).Response(context.RequestCtx)
		return nil
	}

	Logger.Info("DeleteUser ok op_userid=%v op_username=%v userid=%v username=%v ip=%v", opUserID, opUserName, userid, username, context.RequestCtx.RemoteIP())

	httputil.Result.Success("", "OK").Response(context.RequestCtx)
	return nil
}

func (u *userHandler) ResetPwd(context *routing.Context) error {
	userid := httputil.GetParams(context.RequestCtx, "userid", "")
	password1 := httputil.GetParams(context.RequestCtx, "password1", "")
	password2 := httputil.GetParams(context.RequestCtx, "password2", "")

	if userid == "" {
		httputil.Result.Fail("请选择要修改密码的用户").Response(context.RequestCtx)
		return nil
	}
	if password1 == "" || password2 == "" {
		httputil.Result.Fail("请填写新密码并确认").Response(context.RequestCtx)
		return nil
	}
	if password1 != password2 {
		httputil.Result.Fail("两次输入的新密码不一致").Response(context.RequestCtx)
		return nil
	}
	if len(password1) < 6 || len(password1) > 12 {
		httputil.Result.Fail("新密码长度需为6-12位").Response(context.RequestCtx)
		return nil
	}

	if err := db.UpdatePwd(userid, password1); err != nil {
		httputil.Result.Fail("修改密码失败:" + err.Error()).Response(context.RequestCtx)
		return nil
	}

	httputil.Result.Success("", "OK").Response(context.RequestCtx)
	return nil
}

func (u *userHandler) Pwd(context *routing.Context) error {

	id := httputil.GetParams(context.RequestCtx, "id", "")
	user := db.GetUser(id)
	password := httputil.GetParams(context.RequestCtx, "password", "")
	password1 := httputil.GetParams(context.RequestCtx, "password1", "")
	password2 := httputil.GetParams(context.RequestCtx, "password2", "")

	if password1 != password2 {
		httputil.Result.Fail("登录用户密码不一致，请确认密码一致").Response(context.RequestCtx)
		return nil
	}

	if user.UserPassword != SaltMd5(password, user.UserID) {
		httputil.Result.Fail("登录用户名和用户密码不正确").Response(context.RequestCtx)
		return nil
	}

	if error := db.UpdatePwd(id, password1); error != nil {
		httputil.Result.Fail("登录用户密码修改本版本暂时不支持:" + error.Error()).Response(context.RequestCtx)
		Logger.Debug("%v", error)
	} else {
		httputil.Result.Success("", "OK").Response(context.RequestCtx)
	}

	return nil
}
