package controllers

import (
	"goku/models/admin"
	"goku/types"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
	"github.com/google/uuid"
)

// BaseController 基础控制器
type BaseController struct {
	beego.Controller
	RequestID string
	UserName  string
	UserId    int
}

func (b *BaseController) Error(code int, text string) {
	resp := types.APIResponseContext{
		RequestID: b.RequestID,
		Code:      code,
		Error:     text,
	}
	b.Ctx.Output.SetStatus(code)
	b.Ctx.Output.Body([]byte(resp.String()))
	b.Finish()
}

func (b *BaseController) Prepare() {
	uid, _ := uuid.NewUUID()
	b.RequestID = uid.String()

	// require auth
	if strings.HasPrefix(b.Ctx.Request.URL.Path, "/v1") {
		username := b.Ctx.Input.Header(types.COOKIE_USER)
		token := b.Ctx.Input.Header(types.COOKIE_TOKEN)
		requestId := b.Ctx.Input.Header(types.COOKIE_REQUEST)
		if len(username) == 0 {
			b.Error(http.StatusUnauthorized, "miss username of header")
			return
		}
		if len(token) == 0 {
			b.Error(http.StatusUnauthorized, "miss token of header")
			return
		}
		if len(requestId) == 0 {
			b.Error(http.StatusUnauthorized, "miss requestId of header")
			return
		}
		b.Authority(username, token, requestId)
	}
}

func (b *BaseController) Authority(username string, token string, requestId string) {
	if b.Ctx.GetCookie(types.COOKIE_TOKEN) != token {
		b.Error(http.StatusUnauthorized, "cookie:token error")
		return
	}
	if b.Ctx.GetCookie(types.COOKIE_USER) != username {
		b.Error(http.StatusUnauthorized, "cookie:username error")
		return
	}
	service := new(admin.Account)
	params := types.APIRequestContext{
		UserName:  username,
		Token:     token,
		RequestID: requestId,
	}
	userInfo, err := service.Authority(params)
	if err != nil {
		b.Error(http.StatusUnauthorized, err.Error())
		return
	}
	sessionService := new(admin.Session)
	_, err = sessionService.Authority(userInfo.Id, params.Token)
	if err != nil {
		b.Error(http.StatusUnauthorized, err.Error())
		return
	}
	b.UserName = username
	b.UserId = userInfo.Id
}

// 输出status和data的json
func (b *BaseController) RenderJson(msg interface{}, error string, code int) {
	if code != http.StatusOK {
		b.Error(code, error)
		b.StopRun()
	}
	out := types.APIResponseContext{
		RequestID: b.RequestID,
		Code:      code,
		Data:      msg,
		Error:     error,
	}
	b.Data["json"] = out
	b.ServeJSON()
	b.StopRun()
}
