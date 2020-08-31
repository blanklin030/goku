package admin

import (
	"goku/types"

	"github.com/astaxie/beego/orm"
)

type Account struct {
	Id       int    `json:"id" orm:"pk"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *Account) TableName() string {
	return "tb_account"
}

func (u *Account) NewOrm() orm.Ormer {
	u.Clear()
	orm.RegisterModel(new(Account))
	//eorm := orm.NewOrm()
	eorm := types.GetConn("default")
	return eorm
}

func (u *Account) Clear() {
	orm.ResetModelCache()
}

func (u *Account) Login(params types.AdminLogin) (*Account, error) {
	err := u.NewOrm().QueryTable(u.TableName()).Filter("username", params.UserName).Filter("password", params.Password).One(u)
	if err != nil {
		return u, err
	}
	return u, err
}

func (u *Account) Authority(params types.APIRequestContext) (*Account, error) {
	err := u.NewOrm().QueryTable(u.TableName()).Filter("username", params.UserName).One(u)
	if err != nil {
		return u, err
	}
	return u, nil
}
