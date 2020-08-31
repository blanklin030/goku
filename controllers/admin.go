package controllers

import (
	"encoding/json"
	"goku/models/admin"
	"goku/types"
	"net/http"

	"github.com/google/uuid"

	"github.com/astaxie/beego"
)

// Operations about Users
type AdminController struct {
	beego.Controller
}

// @Title Login
// @Description Logs user into the system
// @Param	body		body 	types.AdminLogin	true		"body for user content"
// @Success 200 {object} types.APIResponseContext
// @Failure 403 user not exist
// @router /login [post]
func (u *AdminController) Login() {
	var params types.AdminLogin
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &params)
	userInfo, err := new(admin.Account).Login(params)
	if err != nil {
		u.RenderJson("", err.Error(), http.StatusForbidden)
		return
	}
	tokenInfo, err := new(admin.Session).Login(userInfo.Id)
	if err != nil {
		u.RenderJson("", err.Error(), http.StatusForbidden)
		return
	}
	// save cookie
	u.Ctx.SetCookie(types.COOKIE_TOKEN, tokenInfo.Token, tokenInfo.Expire, "/")
	u.Ctx.SetCookie(types.COOKIE_USER, params.UserName, tokenInfo.Expire, "/")
	u.RenderJson("login success", "", http.StatusOK)
}

// @Title logout
// @Description Logs out current logged in user session
// @Success 200 {object} types.APIResponseContext
// @router /logout [post]
func (u *AdminController) Logout() {
	// delete cookie
	u.Ctx.SetCookie(types.COOKIE_TOKEN, "", -1, "/")
	u.Ctx.SetCookie(types.COOKIE_USER, "", -1, "/")
	u.RenderJson("logout success", "", http.StatusOK)
}

// 输出status和data的json
func (b *AdminController) RenderJson(msg interface{}, error string, code int) {
	uid, _ := uuid.NewUUID()
	if code != http.StatusOK {
		resp := types.APIResponseContext{
			RequestID: uid.String(),
			Code:      code,
			Error:     error,
		}
		b.Ctx.Output.SetStatus(code)
		b.Ctx.Output.Body([]byte(resp.String()))
		b.Finish()
		b.StopRun()
	}
	out := types.APIResponseContext{
		RequestID: uid.String(),
		Code:      code,
		Data:      msg,
		Error:     error,
	}
	b.Data["json"] = out
	b.ServeJSON()
	b.StopRun()
}
