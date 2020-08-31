package types

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var db = make(map[string]*sql.DB)

// Config mysql数据库配置
type MysqlConfig struct {
	Host         string
	Port         string
	Database     string
	Username     string
	Password     string
	TimeZone     string
	MaxIdleConns int
	MaxOpenConns int
}

type Manager struct {
	Env      string
	Config   map[string]MysqlConfig
	TypeName string
	Default  string
	Database map[string]*sql.DB
}

func (m *Manager) GetEnv() string {
	return beego.AppConfig.String("runmode")
}

func (m *Manager) GetFileConfig() (config.Configer, error) {
	filePath := "conf/" + m.TypeName + "/" + m.Env + ".conf"
	logs.Info("read config path: " + filePath)
	cfg, err := config.NewConfig("ini", filePath)
	return cfg, err
}

func (m *Manager) Init(flag string) {
	m.TypeName = flag
	m.Env = m.GetEnv()
	m.InitConfig()
	m.InitDatabase()
}

func (m *Manager) InitConfig() {
	cfg, err := m.GetFileConfig()
	if err != nil {
		logs.Error(err)
		panic(err)
	}
	m.Default = cfg.String(m.TypeName + "::default")
	dbList := strings.Split(cfg.String(m.TypeName+"::database"), ",")
	if len(dbList) == 0 {
		panic("配置有错误")
	}
	m.Config = make(map[string]MysqlConfig)
	for _, value := range dbList {
		m.Config[value] = MysqlConfig{
			Host:         cfg.String(value + "::host"),
			Port:         cfg.String(value + "::port"),
			Database:     cfg.String(value + "::dbname"),
			Username:     cfg.String(value + "::user"),
			Password:     cfg.String(value + "::password"),
			TimeZone:     cfg.String(value + "::timezone"),
			MaxIdleConns: 30,
			MaxOpenConns: 30,
		}
	}
}

func (m *Manager) InitDatabase() {
	logs.Info("init " + m.TypeName + " database")
	err := orm.RegisterDriver(m.TypeName, orm.DRMySQL)
	if err != nil {
		logs.Error(err)
		panic(err)
	}
	m.Database = make(map[string]*sql.DB, 0)
	for key, value := range m.Config {
		dataSource := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=%s",
			value.Username,
			value.Password,
			value.Host,
			value.Port,
			value.Database,
			url.QueryEscape(value.TimeZone))
		logs.Info("datasource=[%s]", dataSource)
		// beego需要一个数据库alias = default
		if key == m.Default {
			key = "default"
		}
		err = orm.RegisterDataBase(key, m.TypeName, dataSource, value.MaxIdleConns, value.MaxOpenConns)
		if err != nil {
			logs.Error(err)
			panic(err)
		}
		if m.Env != "prod" {
			orm.Debug = true
		}
		db, err := orm.GetDB(key)
		if err != nil {
			logs.Error(err)
			panic(err)
		}
		m.Database[key] = db
	}
}

func InitMysql() {
	mgr := new(Manager)
	mgr.Init("mysql")
	db = mgr.Database
}

func GetConn(alias string) orm.Ormer {
	if db[alias] == nil {
		panic("请先初始化" + alias + "数据库")
	}
	ormer, err := orm.NewOrmWithDB("mysql", alias, db[alias])
	if err != nil {
		panic(err)
	}
	err = ormer.Using(alias)
	if err != nil {
		panic(err)
	}
	return ormer
}
