package utils

import (
	"Crocodile6/internal/conf"
	"Crocodile6/internal/consts"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

func InitDB(c conf.Mysql) error {
	var dbConfig gorm.Dialector
	var err error
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/?charset=%v&parseTime=True&multiStatements=true&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		//Conf.Get("mysql.database"),
		"utf8mb4",
	)
	dbConfig = mysql.New(mysql.Config{
		DSN: dsn,
	})
	db, err := gorm.Open(dbConfig, &gorm.Config{
		// 日志功能
		// Logger: logger.Default.LogMode(logger.Info),
		// 外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 禁用默认事务（提高运行速度）
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			// 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
			SingularTable: true,
		},
	})
	if err != nil {
		return err
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * 10)
	consts.DB = db
	return nil
}
