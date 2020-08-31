package user

import (
	"fmt"
	"goku/types"

	"github.com/astaxie/beego/orm"
)

type User struct {
	Id       int    `json:"id" orm:"pk"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) TableName() string {
	return "tb_user"
}

func (u *User) NewOrm() orm.Ormer {
	u.Clear()
	orm.RegisterModel(new(User))
	//eorm := orm.NewOrm()
	eorm := types.GetConn("goku_user")
	return eorm
}

func (u *User) Clear() {
	orm.ResetModelCache()
}

func (u *User) Add(params types.AddUser) (int, error) {
	conn := u.NewOrm()
	var user = User{
		Username: params.UserName,
		Password: params.Password,
	}
	var profile = Profile{
		Gender:  params.Gender,
		Age:     params.Age,
		Address: params.Address,
		Email:   params.Email,
	}
	err := conn.Read(&user)
	// exist
	if err == nil {
		return 0, fmt.Errorf("%s 已存在", user.Username)
	}
	id, err := conn.Insert(&user)
	if err != nil {
		return 0, err
	}
	profile.UserId = int(id)
	_, err = new(Profile).Add(profile)
	if err != nil {
		return 0, err
	}
	return int(id), err
}

func (u *User) Update(userId int, params types.AddUser) (int, error) {
	conn := u.NewOrm()
	var user = User{
		Id:       userId,
		Username: params.UserName,
		Password: params.Password,
	}
	var profile = Profile{
		Gender:  params.Gender,
		Age:     params.Age,
		Address: params.Address,
		Email:   params.Email,
		UserId:  userId,
	}
	err := conn.Read(&user)
	// not exist
	if err != nil {
		return 0, fmt.Errorf("%s 不存在", user.Username)
	}
	_, err = conn.Update(&user, "username", "password")
	if err != nil {
		return 0, err
	}
	_, err = new(Profile).Update(profile)
	if err != nil {
		return 0, err
	}
	return userId, err
}

func (u *User) Delete(userId int) error {
	conn := u.NewOrm()
	var params = User{
		Id: userId,
	}
	if _, err := conn.Delete(&params); err != nil {
		return err
	}
	return nil
}

func (u *User) GetInfo(userId int) (*types.UserInfo, error) {
	conn := u.NewOrm()
	err := conn.QueryTable(u.TableName()).Filter("id", userId).One(u)
	if err != nil {
		return &types.UserInfo{}, err
	}
	profile, err := new(Profile).GetInfoByUserId(userId)
	if err != nil {
		return &types.UserInfo{}, err
	}
	var info = &types.UserInfo{
		UserID:   u.Id,
		Username: u.Username,
		Gender:   profile.Gender,
		Age:      profile.Age,
		Email:    profile.Email,
		Address:  profile.Address,
	}
	return info, err
}

func (u *User) GetList() ([]*types.UserInfo, error) {
	conn := u.NewOrm()
	var list = make([]*User, 0)
	var res = make([]*types.UserInfo, 0)
	_, err := conn.QueryTable(u.TableName()).All(&list)
	if err != nil {
		return res, err
	}
	for _, value := range list {
		profile, err := new(Profile).GetInfoByUserId(value.Id)
		if err != nil {
			continue
		}
		var item = &types.UserInfo{
			UserID:   value.Id,
			Username: value.Username,
			Gender:   profile.Gender,
			Age:      profile.Age,
			Email:    profile.Email,
			Address:  profile.Address,
		}
		res = append(res, item)
	}
	return res, err
}
