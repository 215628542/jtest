package testmodel

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func init() {

	var err error
	//dsn := "select:ZEsqk8J_8fZGx7Z@tcp(gz-tdsqlshard-1i3vlmv7.sql.tencentcdb.com:23750)/vmall?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := ""
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err.Error())
	}
}

type GormStrSlice []string

func (tp GormStrSlice) Value() (driver.Value, error) {
	if len(tp) == 0 {
		return []byte("[]"), nil
	}

	return json.Marshal(tp)
}

func (tp *GormStrSlice) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return json.Unmarshal(b, tp)
}

type GormJson []byte

func (tp GormJson) Value() (driver.Value, error) {
	if len(tp) == 0 {
		return []byte("{}"), nil
	}

	return tp, nil
}

func (tp *GormJson) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	*tp = b
	return nil
}
