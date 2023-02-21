/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : mysql.go
*   coder: zemanzeng
*   date : 2021-09-29 16:09:57
*   desc : mysql连接
*
================================================================*/

package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
)

var dbs = new(sync.Map)
var _ = &mysql.MySQLDriver{} // keep

// MysqlSource mysql连接相关配置
type MysqlSource struct {
	Source       string `json:"source" yaml:"source"`
	MaxOpenConns int    `json:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns int    `json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxLifeTime  int    `json:"max_life_time" yaml:"max_life_time"`
}

func GetDB(source MysqlSource) (*sql.DB, error) {
	if len(source.Source) == 0 {
		return nil, errors.New("invalid source")
	}
	if value, exist := dbs.Load(source.Source); exist {
		if db, ok := value.(*sql.DB); ok {
			return db, nil
		}
		return nil, fmt.Errorf("invalid db:%v type:%T", source.Source, value)
	}

	db, err := sql.Open("mysql", source.Source)
	if err != nil {
		return db, err
	}
	if err := db.Ping(); err != nil {
		return db, fmt.Errorf("ping error:%w", err)
	}
	dbs.Store(source.Source, db)

	db.SetMaxOpenConns(GetNonZero(source.MaxOpenConns, 10))
	db.SetMaxIdleConns(GetNonZero(source.MaxIdleConns, 2))

	lifeTime := GetNonZero(source.MaxLifeTime, 1*60*60)
	db.SetConnMaxLifetime(time.Duration(lifeTime) * time.Second)
	return db, nil
}

var DBErr_AffectedIsZero = errors.New("affected is zero")
var DBErr_NotFindMatchData = errors.New("not find match data")

func CheckExecResult(result sql.Result, expect int) error {

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return DBErr_AffectedIsZero
	}
	if expect > 0 && affected != int64(expect) {
		return fmt.Errorf("affected:%d not equal expected:%d", affected, expect)
	}

	return nil
}

func I64Arr2String(arr []int64) string {
	if len(arr) == 0 {
		return "()"
	}

	var s string
	for i, item := range arr {
		if i == 0 {
			s = strconv.FormatInt(item, 10)
		} else {
			s = s + "," + strconv.FormatInt(item, 10)
		}
	}

	return "(" + s + ")"
}

func StringArr2String(arr []string) string {
	if len(arr) == 0 {
		return "()"
	}

	var s string
	for i, item := range arr {
		if i == 0 {
			s = "'" + item + "'"
		} else {
			s = s + ", '" + item + "'"
		}
	}

	return "(" + s + ")"
}

func GetNonZero(nums ...int) int {
	for _, number := range nums {
		if number != 0 {
			return number
		}
	}
	return 0
}

func PlaceHolder(count int) string {
	var holder string
	for index := 0; index < count; index++ {
		if index == 0 {
			holder = holder + "?"
			continue
		}
		holder = holder + ", ? "
	}
	return holder
}

func ParseTime(t string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", t)
}
