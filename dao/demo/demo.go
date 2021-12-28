package demo

import (
	"context"
	"time"

	"nautilus/pkg/sqlx"
)

type Item struct {
	ID         int32     `db:"id"`
	Name       string    `db:"name"`
	CreateTime time.Time `db:"create_time"`
	ModifyTime time.Time `db:"modify_time"`
}

func QueryByID(ctx context.Context, id int32) (item Item, err error) {
	conn := sqlx.Get(ctx, "default")
	sql := "select * from t_demo where id = ?"

	err = conn.GetContext(ctx, &item, sql, id)
	return
}
