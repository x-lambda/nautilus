package admin

import (
	"context"
	"time"

	"nautilus/pkg/sqlx"
)

// Profile 管理员信息
type Profile struct {
	ID       int64     `db:"id"`
	Username string    `db:"username"`
	Password string    `db:"password"`
	Phone    string    `db:"phone"`     // 手机号，普通管理员可以通过账号密码登录，超级管理员必须通过手机号登录
	RoleType int32     `db:"role_type"` // 0: 普通管理员  1: 超级管理员
	CTime    time.Time `db:"ctime"`     // 创建时间
	MTime    time.Time `db:"mtime"`     // 修改时间
}

// TableName 返回表名，必须实现
func (p Profile) TableName() string {
	return "t_admin"
}

// KeyName 返回主键，必须实现
func (p Profile) KeyName() string {
	return "id"
}

// CreateAdmin 创建管理员账号
func CreateAdmin(ctx context.Context, p Profile) (id int64, err error) {
	conn := sqlx.Get(ctx, "pension")

	now := time.Now()
	if p.CTime.IsZero() {
		p.CTime = now
	}

	if p.MTime.IsZero() {
		p.MTime = now
	}

	result, err := conn.InsertContext(ctx, &p)
	if err != nil {
		return
	}

	id, err = result.LastInsertId()
	return
}

// QueryByUsername 根据用户名查询
func QueryByUsername(ctx context.Context, username string) (p Profile, err error) {
	if username == "" {
		return
	}

	conn := sqlx.Get(ctx, "pension")
	err = conn.GetContext(ctx, &p, "select * from t_admin where username=?", username)

	// 如果没查询到，则id为0
	if sqlx.IsNoRowErr(err) {
		err = nil
	}

	return
}

// QueryByUID 根据uid查询
func QueryByUID(ctx context.Context, uid int64) (p Profile, err error) {
	if uid == 0 {
		return
	}

	conn := sqlx.Get(ctx, "pension")
	err = conn.GetContext(ctx, &p, "select * from t_admin where id=?", uid)

	// 如果没查询到，则id为0
	if sqlx.IsNoRowErr(err) {
		err = nil
	}
	return
}

// deleteByID 删除指定用户
func deleteByID(ctx context.Context, uid int64) (err error) {
	conn := sqlx.Get(ctx, "pension")
	p := Profile{ID: uid}
	_, err = conn.DeleteContext(ctx, p)
	return
}
