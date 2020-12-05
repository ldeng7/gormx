package gormx

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Conf struct {
	Dsn     string
	Mysql   *mysql.Config
	Gorm    *gorm.Config
	MaxIdle int
	MaxOpen int
}

func InitGormMysql(conf *Conf) (*gorm.DB, error) {
	var dial gorm.Dialector
	if nil == conf.Mysql {
		dial = mysql.Open(conf.Dsn)
	} else {
		dial = mysql.New(*conf.Mysql)
	}
	db, err := gorm.Open(dial, conf.Gorm)
	if nil != err {
		return nil, err
	}

	maxIdle, maxOpen := conf.MaxIdle, conf.MaxOpen
	if maxIdle <= 0 {
		maxIdle = 8
	}
	if maxOpen <= 0 {
		maxOpen = 128
	}
	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(maxIdle)
	sqlDb.SetMaxOpenConns(maxOpen)
	return db, nil
}
