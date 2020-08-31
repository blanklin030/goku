package user

import (
	"fmt"
	"goku/types"

	"github.com/astaxie/beego/orm"
)

type Profile struct {
	Id      int    `json:"id" orm:"pk"`
	UserId  int    `json:"user_id"`
	Gender  string `json:"gender"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Email   string `json:"email"`
}

func (p *Profile) TableName() string {
	return "tb_profile"
}

func (p *Profile) NewOrm() orm.Ormer {
	p.Clear()
	orm.RegisterModel(new(Profile))
	//eorm := orm.NewOrm()
	eorm := types.GetConn("goku_user")
	return eorm
}

func (p *Profile) Clear() {
	orm.ResetModelCache()
}

func (p *Profile) GetInfoByUserId(userId int) (*Profile, error) {
	err := p.NewOrm().QueryTable(p.TableName()).Filter("user_id", userId).One(p)
	if err != nil {
		return p, err
	}
	return p, err
}

func (p *Profile) Add(profile Profile) (int, error) {
	err := p.NewOrm().Read(&profile)
	if err == nil {
		return 0, fmt.Errorf("%s 已存在", profile.UserId)
	}
	id, err := p.NewOrm().Insert(&profile)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (p *Profile) Update(profile Profile) (int, error) {
	err := p.NewOrm().QueryTable(p.TableName()).Filter("user_id", profile.UserId).One(p)
	if err != nil {
		return 0, fmt.Errorf("%v 不存在", profile.UserId)
	}
	profile.Id = p.Id
	id, err := p.NewOrm().Update(&profile, "gender", "age", "address", "email")
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
