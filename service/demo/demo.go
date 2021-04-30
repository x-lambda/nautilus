package demo

import (
	"context"
	"fmt"

	"nautilus/dao/demo"
)

func TestTimeout(ctx context.Context) (err error) {
	// fmt.Println("你看不到我😛")
	item, err := demo.QueryByID(ctx, 11)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(item.ID)
	return
}
