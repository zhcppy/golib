package db

import (
	"errors"
	"fmt"
	"strings"
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
	Port         int64  `json:"port"`
	MaxIdleConns int64  `json:"maxIdleConns"`
	MaxOpenConns int64  `json:"maxOpenConns"`
	TablePrefix  string `json:"tablePrefix"`
	IsLog        bool   `json:"isLog"`
	Init         bool   `json:"init"`
	IsDrop       bool   `json:"isDrop"`
}

func DefMysqlCfg() *Config {
	return &Config{
		Username:     "root",
		Password:     "12345678",
		Host:         "127.0.0.1",
		Port:         3306,
		Database:     "test",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		TablePrefix:  "",
		IsLog:        false,
		Init:         true,
		IsDrop:       false,
	}
}

func NewMysql(cfg *Config, models ...interface{}) (*gorm.DB, error) {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		if strings.HasPrefix(defaultTableName, cfg.TablePrefix) {
			return defaultTableName
		}
		return cfg.TablePrefix + "_" + defaultTableName
	}
	if err := CreateDB(cfg); err != nil {
		return nil, err
	}

	cStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	openedDb, err := gorm.Open("mysql", cStr)
	if err != nil {
		logger.L.Error("mysql open connection error: ", err.Error())
		return nil, err
	}
	logger.L.Info("init mysql db connection, db_host: ", cfg.Host, " user: ", cfg.Username, " db_name: ", cfg.Database)

	openedDb.DB().SetMaxIdleConns(int(cfg.MaxIdleConns)) //最大空闲连接数
	openedDb.DB().SetMaxOpenConns(int(cfg.MaxOpenConns)) //最大连接池大小
	openedDb.DB().SetConnMaxLifetime(time.Hour * 2)      //连接最大存活时间

	openedDb.SingularTable(true)
	openedDb.LogMode(cfg.IsLog) // 启用详细日志（包含sql语句）

	logger.L.Info("database connection successful by mysql, db_name: ", cfg.Database)

	if cfg.IsDrop {
		logger.L.Info("===========> create table")
		if err := openedDb.CreateTable(models...).Error; err != nil {
			return nil, err
		}
	}

	if cfg.Init {
		logger.L.Info("===========> auto migrate table")
		if err := openedDb.AutoMigrate(models...).Error; err != nil {
			return nil, err
		}
	}
	return openedDb, nil
}

func CreateDB(cfg *Config) (err error) {
	if cfg == nil {
		return errors.New("mysql connection config is nil")
	}
	if !cfg.Init {
		return nil
	}
	logger.L.Info("===========> init database")

	cStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/information_schema?charset=utf8&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	openedDb, err := gorm.Open("mysql", cStr)
	if err != nil {
		return err
	}
	defer openedDb.Close()
	openedDb.LogMode(cfg.IsLog)

	if cfg.IsDrop {
		dropDbSQL := "DROP DATABASE IF EXISTS " + cfg.Database + ";"
		if err := openedDb.Exec(dropDbSQL).Error; err != nil {
			logger.L.Error("failed to mysql drop database, sql: ", dropDbSQL, "err: ", err.Error())
			return err
		}
		logger.L.Info("success to drop mysql db: ", cfg.Database)
	}

	createDbSQL := "CREATE DATABASE IF NOT EXISTS " + cfg.Database + " DEFAULT CHARSET utf8 COLLATE utf8_general_ci;"
	if err := openedDb.Exec(createDbSQL).Error; err != nil {
		logger.L.Error("failed to mysql create database, sql: ", createDbSQL, "err: ", err.Error())
		return err
	}
	logger.L.Info("success to create mysql db: ", cfg.Database)
	return nil
}
