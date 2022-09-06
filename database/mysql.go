package database

import (
	"fmt"

	"github.com/fimreal/go-gin-jwt/database/models"
	verify "github.com/fimreal/goutils/parse"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewORM(dsn string) (*DBORM, error) {
	gormOpts := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Info),
	}
	if ShowSQL {
		gormOpts.Logger = logger.Default.LogMode(logger.Info)
	}

	gormDB, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN: dsn,
		}), gormOpts,
	)
	return &DBORM{
		DB: gormDB,
	}, err
}

func (db *DBORM) AddUser(user models.User) (models.User, error) {
	err := db.Create(&user).Error
	return user, err
}

func (db *DBORM) GetUser(account interface{}) (models.User, error) {
	var u models.User
	var err error

	switch t := account.(type) {
	case int:
		err = db.Model(&models.User{}).Select("user_id", "username", "password", "level").Where("user_id = ?", account).Scan(&u).Error
	case string:
		if verify.IsEmail(account.(string)) {
			err = db.Model(&models.User{}).Select("user_id", "username", "password", "level").Where("email = ?", account).Scan(&u).Error
		} else {
			err = db.Model(&models.User{}).Select("user_id", "username", "password", "level").Where("username = ?", account).Scan(&u).Error
		}
	default:
		return u, fmt.Errorf("err: not support type %v[%v]", t, account)
	}

	return u, err
}
