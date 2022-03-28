package connector

import (
	"bytes"
	"fmt"
	"gorm.io/gorm/schema"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/url"
	"strconv"
)

type mysqlConfig struct {
	UserName          string `json:"useName"`
	Password          string `json:"password"`
	Host              string `json:"host"`
	Port              string `json:"port"`
	DataBase          string `json:"dataBase"`
	Charset           string `json:"charset"`
	ParseTime         bool   `json:"parseTime"`
	MultiStatements   bool   `json:"multiStatements"`
	Loc               string `json:"local"`
	MaxOpenConns      int    `json:"maxOpenConns"`
	MaxIdleConns      int    `json:"maxIdleConns"`
	ConnMaxLifeSecond int    `json:"connMaxLifeSecond"`
}

var defaultMysqlConfig = &mysqlConfig{
	UserName:        "root",
	Password:        "123756",
	Host:            "localhost",
	Port:            "3306",
	DataBase:        "common",
	Charset:         "utf8mb4",
	ParseTime:       true, //自动将时间戳转换时间
	MultiStatements: true,
	Loc:             "Local",
	MaxOpenConns:    8,
	MaxIdleConns:    8,
}

func NewDefaultMysqlConfig() *mysqlConfig {
	return defaultMysqlConfig
}

func NewEmptyMysqlConfig() *mysqlConfig {
	return &mysqlConfig{}
}

func (m *mysqlConfig) GetMySQLEngine() (db *gorm.DB, err error) {
	// 连接中间件实例
	db, err = gorm.Open(mysql.Open(m.GetDsnByString()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		}})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if m.ConnMaxLifeSecond > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(m.ConnMaxLifeSecond) * time.Second)
	}
	sqlDB.SetMaxIdleConns(m.MaxIdleConns)
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)

	return db, nil

}

func (m mysqlConfig) GetDsnByString() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=%v&multiStatements=%v&loc=%v",
		m.UserName,
		m.Password,
		m.Host,
		m.Port,
		m.DataBase,
		m.Charset,
		m.ParseTime,
		m.MultiStatements,
		m.Loc,
	)
}

func (m mysqlConfig) GetDsnByBuffer() string {
	var buf bytes.Buffer
	buf.WriteString(m.UserName)
	buf.WriteString(":")
	buf.WriteString(m.Password)
	buf.WriteString("@tcp(")
	buf.WriteString(m.Host)
	buf.WriteString(")/")
	if m.DataBase != "" {
		buf.WriteString(m.DataBase)
	}
	buf.WriteString("?parseTime=" + strconv.FormatBool(m.ParseTime))
	buf.WriteString("&multiStatements=" + strconv.FormatBool(m.MultiStatements))
	if m.Charset != "" {
		buf.WriteString("&charset=" + m.Charset)
	}
	if m.Loc == "" {
		buf.WriteString("&loc=Local")
	} else {
		buf.WriteString("&loc=" + url.QueryEscape(m.Loc))
	}
	return buf.String()
}
