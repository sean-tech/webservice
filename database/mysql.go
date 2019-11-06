package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sean-tech/webservice/config"
)

var (
	/** 数据中心id 关联 db Map **/
	dbMap map[int]*sqlx.DB
	/** 数据中心数量 **/
	dataCenterCount int = 0
)

type mysqlManagerImpl struct {}

/**
 * 数据库open
 * db: DB 对象
 */
func (this *mysqlManagerImpl) Open() {
	var (
		dbType, dbName, user, password string
	)
	dbType = config.Database.Type
	dbName = config.Database.Name
	user = config.Database.User
	password = config.Database.Password
	dbMap = make(map[int]*sqlx.DB)
	for id, host := range config.Database.Hosts {
		var dbLink = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName)
		db, err := sqlx.Open(dbType, dbLink)
		if err != nil {
			panic(err)
		}
		db.SetMaxIdleConns(config.Database.MaxIdle)
		db.SetMaxOpenConns(config.Database.MaxOpen)
		db.SetConnMaxLifetime(config.Database.MaxLifetime)
		dbMap[id] = db
		dataCenterCount += 1
	}
}

/**
 * 根据用户名基因确定数据库对象
 */
func (this *mysqlManagerImpl) GetDbByUserName(userName string) (db *sqlx.DB, err error) {
	dna, err := Dna(userName)
	if err != nil {
		return nil, err
	}
	dataCenterId := dna % dataCenterCount
	return dbMap[dataCenterId], nil
}

/**
 * 获取所有数据库对象
 */
func (this *mysqlManagerImpl) GetAllDbs() (dbs []*sqlx.DB) {
	for _, v := range dbMap {
		dbs = append(dbs, v)
	}
	return dbs
}

//var (
//	tokenBucketOnce sync.Once
//	tokenBucket		*ratelimit.Bucket
//)
//
//func getTokenBucket() *ratelimit.Bucket {
//	tokenBucketOnce.Do(func() {
//		tokenBucket	= ratelimit.NewBucket(config.Server.RateLimitFillInterval, config.Server.RateLimitCapacity)
//	})
//	return tokenBucket
//}

