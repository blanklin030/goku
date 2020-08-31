package controllers

import (
	"encoding/json"
	"goku/models/user"
	"goku/types"
	"net/http"
)

// Operations about Users
type UserController struct {
	BaseController
}

// @Title CreateUser
// @Description create users
// @Param	GOKU_SSO_TOKEN		header 	string	true	"set header of token"
// @Param	GOKU_SSO_USER		header 	string	true	"set header of user"
// @Param	GOKU_REQUEST_ID		header 	string	true	"set header of request id"
// @Param	body		body 	types.AddUser	true		"body for user content"
// @Success 200 {object} types.APIResponseContext
// @Failure 500 {object} types.APIResponseContext
// @router / [post]
func (u *UserController) Post() {
	var params types.AddUser
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &params)
	if err != nil {
		u.RenderJson(0, err.Error(), http.StatusInternalServerError)
	}
	model := new(user.User)
	uid, err := model.Add(params)
	if err != nil {
		u.RenderJson("", err.Error(), http.StatusInternalServerError)
	} else {
		u.RenderJson(uid, "", http.StatusOK)
	}
}

// @Title GetAll
// @Description get all Users
// @Param	GOKU_SSO_TOKEN		header 	string	true	"set header of token"
// @Param	GOKU_SSO_USER		header 	string	true	"set header of user"
// @Param	GOKU_REQUEST_ID		header 	string	true	"set header of request id"
// @Success 200 {object} types.APIResponseContext
// @router / [get]
func (u *UserController) GetAll() {
	model := new(user.User)
	users, err := model.GetList()
	if err != nil {
		u.RenderJson("", err.Error(), http.StatusInternalServerError)
	} else {
		u.RenderJson(users, "", http.StatusOK)
	}
}

// @Title Get
// @Description get user by uid
// @Param	GOKU_SSO_TOKEN		header 	string	true	"set header of token"
// @Param	GOKU_SSO_USER		header 	string	true	"set header of user"
// @Param	GOKU_REQUEST_ID		header 	string	true	"set header of request id"
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 500 {object} models.User
// @router /:uid [get]
func (u *UserController) Get() {
	model := new(user.User)
	uid, err := u.GetInt(":uid")
	if err != nil {
		u.RenderJson("", err.Error(), http.StatusInternalServerError)
	}
	userInfo, err := model.GetInfo(uid)
	if err != nil {
		u.RenderJson("", err.Error(), http.StatusInternalServerError)
	} else {
		u.RenderJson(userInfo, "", http.StatusOK)
	}
}

// @Title Update
// @Description update the user
// @Param	GOKU_SSO_TOKEN		header 	string	true	"set header of token"
// @Param	GOKU_SSO_USER		header 	string	true	"set header of user"
// @Param	GOKU_REQUEST_ID		header 	string	true	"set header of request id"
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	types.AddUser	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /:uid [put]
func (u *UserController) Put() {
	model := new(user.User)
	uid, err := u.GetInt(":uid")
	if err != nil {
		u.RenderJson("", err.Error(), http.StatusInternalServerError)
	}
	var params types.AddUser
	err = json.Unmarshal(u.Ctx.Input.RequestBody, &params)
	if err != nil {
		u.RenderJson(0, err.Error(), http.StatusInternalServerError)
	}
	userInfo, err := model.Update(uid, params)
	if err != nil {
		u.RenderJson("", err.Error(), http.StatusInternalServerError)
	} else {
		u.RenderJson(userInfo, "", http.StatusOK)
	}
}

// @Title Delete
// @Description delete the user
// @Param	GOKU_SSO_TOKEN		header 	string	true	"set header of token"
// @Param	GOKU_SSO_USER		header 	string	true	"set header of user"
// @Param	GOKU_REQUEST_ID		header 	string	true	"set header of request id"
// @Param	uid		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /:uid [delete]
func (u *UserController) Delete() {
	uid, err := u.GetInt(":uid")
	if err != nil {
		u.RenderJson("", err.Error(), http.StatusInternalServerError)
	}
	model := new(user.User)
	err = model.Delete(uid)
	if err != nil {
		u.RenderJson("", err.Error(), http.StatusInternalServerError)
	} else {
		u.RenderJson("delete success!", "", http.StatusOK)
	}
}
