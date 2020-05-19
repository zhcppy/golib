package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/zhcppy/golib/logger"
)

type Config struct {
	Host         string `json:"host"`
	Database     string `json:"database"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Port         uint64 `json:"port"`
	MaxIdleConns int    `json:"maxIdleConns"`
	MaxOpenConns int    `json:"maxOpenConns"`
	TablePrefix  string `json:"tablePrefix"`
	IsLog        bool   `json:"isLog"`
	Init         bool   `json:"init"`
}

func DefaultCfg() *Config {
	return &Config{
		Username:     "root",
		Password:     "12345678",
		Host:         "127.0.0.1",
		Port:         3306,
		Database:     "fx_test",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		TablePrefix:  "fx_",
		IsLog:        true,
	}
}

// 创建Mysql连接
func NewMysql(cfg *Config, models ...func() []interface{}) (*gorm.DB, error) {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return cfg.TablePrefix + defaultTableName
	}
	if cfg.Init {
		logger.L.Debugf("===========> init database")
		if err := CreateDB(cfg); err != nil {
			return nil, err
		}
	}

	cStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	openedDb, err := gorm.Open("mysql", cStr)
	if err != nil {
		logger.L.Error("mysql open connection error: ", err.Error())
		return nil, err
	}
	logger.L.Debug("init mysql db connection,", "db_host: ", cfg.Host, "user: ", cfg.Username, "db_name: ", cfg.Database)

	openedDb.DB().SetMaxIdleConns(cfg.MaxIdleConns) //最大空闲连接数
	openedDb.DB().SetMaxOpenConns(cfg.MaxOpenConns) //最大连接池大小
	openedDb.DB().SetConnMaxLifetime(time.Hour * 2) //连接最大存活时间

	openedDb.SingularTable(true)
	openedDb.LogMode(cfg.IsLog) // 启用详细日志（包含sql语句）

	logger.L.Debug("database connection successful by mysql,", "db_name: ", cfg.Database)

	if cfg.Init {
		for _, model := range models {
			if err := openedDb.AutoMigrate(model()...).Error; err != nil {
				return nil, err
			}
		}
	}
	return openedDb, nil
}

// 创建数据库
func CreateDB(cfg *Config) (err error) {
	if cfg == nil {
		return errors.New("mysql connection config is nil")
	}
	cStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/information_schema?charset=utf8&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	openedDb, err := gorm.Open("mysql", cStr)
	if err != nil {
		return err
	}
	defer openedDb.Close()

	createDbSQL := "CREATE DATABASE IF NOT EXISTS " + cfg.Database + " DEFAULT CHARSET utf8 COLLATE utf8_general_ci;"
	if err := openedDb.Exec(createDbSQL).Error; err != nil {
		logger.L.Error("failed to mysql create database, sql: ", createDbSQL, "err: ", err.Error())
		return err
	}
	logger.L.Debug("success to create mysql db: ", cfg.Database)
	return nil
}

// 删除数据库
func DropDB(openedDb *gorm.DB, database string) (err error) {
	dropDbSQL := "DROP DATABASE IF EXISTS " + database + ";"
	if err := openedDb.Exec(dropDbSQL).Error; err != nil {
		logger.L.Error("failed to mysql drop database, sql: ", dropDbSQL, "err: ", err.Error())
		return err
	}
	logger.L.Debug("success to drop mysql db: ", database)
	return nil
}

// 删除所有表中的数据，并不删除表结构
func ClearAllData(openedDb *gorm.DB) (err error) {
	rs, err := openedDb.Raw("SHOW TABLES;").Rows()
	if err != nil {
		logger.L.Error("failed to get table name list: ", err.Error())
		return err
	}
	var tName string
	for rs.Next() {
		if err := rs.Scan(&tName); err != nil {
			logger.L.Error("failed to scam table name: ", tName, "error: ", err.Error())
			return err
		}
		if tName == "" {
			continue
		}
		if err := openedDb.Exec(fmt.Sprintf("DELETE FROM %s", tName)).Error; err != nil {
			logger.L.Error("failed to delete table all data, error: ", err.Error())
			return err
		}
	}
	logger.L.Debug("success to clear all data")
	return nil
}
