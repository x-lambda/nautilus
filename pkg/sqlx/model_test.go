package sqlx

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type user struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

func (u user) TableName() string {
	return "t_test_orm"
}

func (u user) KeyName() string {
	return "id"
}

func TestDelete(t *testing.T) {
	ctx := context.TODO()
	conn := Get(ctx, "")

	u := user{
		Name: "evdsc",
		Age:  10,
	}

	result, err := conn.DeleteContext(ctx, u)
	fmt.Printf("err: %+v\n", err)
	fmt.Println(result.RowsAffected())
}

func TestQuery(t *testing.T) {
	ctx := context.TODO()
	conn := Get(ctx, "")

	u := user{Name: "test1" + time.Now().Format(time.RFC3339), Age: 10}
	result, err := conn.InsertContext(ctx, u)
	assert.Nil(t, err)
	id1, _ := result.LastInsertId()
	fmt.Println("---------------")

	u = user{Name: "test2" + time.Now().Format(time.RFC3339), Age: 10}
	result, err = conn.InsertContext(ctx, u)
	assert.Nil(t, err)
	id2, _ := result.LastInsertId()
	fmt.Println("---------------")

	var dst user
	err = conn.QueryRowxContext(ctx, "SELECT id, name, age FROM t_test_orm where id = ?", id1).Scan(&dst.ID, &dst.Name, &dst.Age)
	assert.Nil(t, err)
	fmt.Println("---------------")

	err = conn.QueryRowxContext(ctx, "SELECT id, name, age FROM t_test_orm where id = ?", id1).Scan(&dst.ID, &dst.Name, &dst.Age)
	assert.Nil(t, err)
	fmt.Println("---------------")

	err = conn.QueryRowxContext(ctx, "SELECT id, name, age FROM t_test_orm where id = ?", id2).Scan(&dst.ID, &dst.Name, &dst.Age)
	assert.Nil(t, err)
	fmt.Println("---------------")
}

func TestModel(t *testing.T) {
	ctx := context.TODO()
	conn := Get(ctx, "")

	u := user{
		Name: "evdsc",
		Age:  10,
	}
	// insert
	result, err := conn.InsertContext(ctx, u)
	if err != nil {
		panic(err)
	}

	// update
	u.ID, _ = result.LastInsertId()
	u.Name = "dscc"
	result, err = conn.UpdateContext(ctx, u)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())

	// get
	var dst user
	err = conn.GetContext(ctx, &dst, "select * from t_test_orm where id = ?", u.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(dst.ID, dst.Age, dst.Name)

	// 查询不到
	err = conn.GetContext(ctx, &dst, "select * from t_test_orm where id = 100")
	fmt.Printf("err: %+v\n", err)

	var all []user
	err = conn.SelectContext(ctx, &all, "select * from t_test_orm order by id desc")
	if err != nil {
		panic(err)
	}
	fmt.Println("---------------")

	result, err = conn.DeleteContext(ctx, u)
	fmt.Printf("err: %+v\n", err)
	fmt.Println(result.RowsAffected())
}

func TestRawExec(t *testing.T) {
	ctx := context.TODO()
	conn := Get(ctx, "")

	// exec insert
	result, err := conn.ExecContext(ctx, "insert into t_test_orm(name, age) values (?, ?)", "c", 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.LastInsertId())
	fmt.Println("------------------------------------------------------")

	result, err = conn.ExecContext(ctx, "update t_test_orm set name=? where id=?", "test", 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())
	fmt.Println("------------------------------------------------------")

	rows, err := conn.QueryContext(ctx, "select id, name, age from t_test_orm")
	if err != nil {
		panic(err)
	}

	var id int
	var name string
	var age int

	for rows.Next() {
		if rows.Scan(&id, &name, &age); err != nil {
			panic(err)
		}
		fmt.Printf("id: %d name: %s age: %d\n", id, name, age)
	}
	fmt.Println("------------------------------------------------------")

	err = conn.QueryRowxContext(ctx, "select id, name, age from t_test_orm where id = ?", 1).Scan(&id, &name, &age)
	if err != nil {
		panic(err)
	}
	fmt.Printf("id: %d name: %s age: %d\n", id, name, age)
	fmt.Println("------------------------------------------------------")

	result, err = conn.ExecContext(ctx, "delete from t_test_orm where id = ?", 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())
	fmt.Println("------------------------------------------------------")
}

func TestTransaction(t *testing.T) {
	ctx := context.TODO()
	conn := Get(ctx, "")

	tx, err := conn.Beginx()
	if err != nil {
		panic(err)
	}

	defer func() {
		if p := recover(); p != nil {
			// 回滚，继续向上panic
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// 回滚，向上抛 err
			tx.Rollback()
		} else {
			// 提交事务
			err = tx.Commit()
		}
	}()

	// 事务1
	u := user{ID: 11, Name: "lalala", Age: 100}
	result, err := tx.UpdateContext(ctx, u)
	if err != nil {
		return
	}

	affect, err := result.RowsAffected()
	if err != nil {
		return
	}

	if affect < 1 {
		err = fmt.Errorf("no affect")
		return
	}

	// 事务2
	err = trans(ctx, tx)
	if err != nil {
		return
	}

	// 事务n...

	return
}

func trans(ctx context.Context, conn *Tx) (err error) {
	u := user{ID: 1000, Name: "None", Age: 999}
	result, err := conn.UpdateContext(ctx, u)
	if err != nil {
		return
	}

	affect, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affect < 1 {
		err = fmt.Errorf("no affect")
		return
	}
	return
}
