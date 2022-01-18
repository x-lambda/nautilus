package admin

import (
	"time"
)

type Token struct {
	ID    int64     `db:"id"`
	UID   int64     `db:"uid"`
	Key   string    `db:"key"`
	CTime time.Time `db:"ctime"`
	MTime time.Time `db:"mtime"`
}

// TableName 返回表名，必须实现
func (t Token) TableName() string {
	return "t_token"
}

// KeyName 返回主键，必须实现
func (t Token) KeyName() string {
	return "id"
}
