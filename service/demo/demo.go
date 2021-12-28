package demo

import (
	"context"

	"nautilus/dao/demo"
	"nautilus/pkg/log"
)

func TestTimeout(ctx context.Context) (err error) {
	item, err := demo.QueryByID(ctx, 11)
	if err != nil {
		return
	}

	log.Get(ctx).Infof("item id: %d", item.ID)
	return
}
