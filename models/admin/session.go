package admin

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"goku/types"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
)

type Session struct {
	Id        int    `json:"id" orm:"pk"`
	AccountId int    `json:"account_id"`
	Expire    int    `json:"expire"`
	Token     string `json:"token"`
}

func (t *Session) TableName() string {
	return "tb_session"
}

func (t *Session) NewOrm() orm.Ormer {
	t.Clear()
	orm.RegisterModel(new(Session))

	//eorm := orm.NewOrm()
	eorm := types.GetConn("default")
	return eorm
}

func (t *Session) Clear() {
	orm.ResetModelCache()
}

func (t *Session) Login(accountId int) (*Session, error) {
	conn := t.NewOrm()
	err := conn.QueryTable(t.TableName()).Filter("account_id", accountId).One(t)
	duration, _ := time.ParseDuration("168h")
	expire := time.Now().Add(duration).Unix()
	md5Helper := md5.New()
	md5Helper.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))
	token := hex.EncodeToString(md5Helper.Sum(nil))
	params := Session{
		AccountId: accountId,
		Expire:    int(expire),
		Token:     token,
	}
	// add Session
	if err != nil {
		id, err := conn.Insert(&params)
		params.Id = int(id)
		return &params, err
	}
	// update Session expire
	params.Id = t.Id
	_, err = conn.Update(&params, "expire")
	return t, nil
}

func (t *Session) Authority(accountId int, session string) (*Session, error) {
	err := t.NewOrm().QueryTable(t.TableName()).Filter("account_id", accountId).One(t)
	if err != nil {
		return t, err
	}
	if t.Token != session {
		return t, fmt.Errorf("session is error")
	}
	// check expire Session
	now := time.Now().Unix()
	if t.Expire < int(now) {
		return t, fmt.Errorf("Session is expire,please login again")
	}
	return t, nil
}
