package dbutils

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sean-tech/webservice/config"
)

/** 数据中心id 关联 db Map **/
var dbMap map[int]*sqlx.DB

/**
 * 数据库open
 * db: DB 对象
 */
func DatabaseOpen() {
	var (
		dbType, dbName, user, password string
	)
	dbType = config.DatabaseSetting.Type
	dbName = config.DatabaseSetting.Name
	user = config.DatabaseSetting.User
	password = config.DatabaseSetting.Password
	dbMap = make(map[int]*sqlx.DB)
	for id, host := range config.DatabaseSetting.Hosts {
		var dbLink = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName)
		db, err := sqlx.Open(dbType, dbLink)
		if err != nil {
			panic(err)
		}
		db.SetMaxIdleConns(config.DatabaseSetting.MaxIdle)
		db.SetMaxOpenConns(config.DatabaseSetting.MaxOpen)
		db.SetConnMaxLifetime(config.DatabaseSetting.MaxLifetime)
		dbMap[id] = db
	}
}

const dataCenterCount int = 1

func DbByUserName(userName string) (db *sqlx.DB, err error) {
	dna, err := Dna(userName)
	if err != nil {
		return nil, err
	}
	dataCenterId := dna % dataCenterCount
	return dbMap[dataCenterId], nil
}

func GetAllDbs() (dbs []*sqlx.DB) {
	for _, v := range dbMap {
		dbs = append(dbs, v)
	}
	return dbs
}

