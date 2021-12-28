# sqlx

`sqlx`是标准库`sql`的升级版， 在`sqlx`的基础上线封装一套新的`api`以适应业务，实现以下目标

- 支持多`db`实例
- 记录`sql`执行日志和`cost`耗时
- 上报`opentracing`追踪数据
- 汇总`prometheus`监控指标
- 尽可能试`api`简洁，越少越好，标准库的`api`实际只有5个（连事务）

`sqlx`存在缺点
- 不够`orm`，只有部分操作可以真正的实现`orm`，例如`insert`，`update`(以`id`为查询条件时)
- 写法跟`orm`不同，需要适应一下
- api略多，5个基础操作+1个复杂`exec`+1个复杂查询+1个事务

### 使用示例
`model`定义
```go
type User struct {
    ID         int
    Name       string
    Age        int
    CreateTime time.Time
}

// TableName 返回表面，必须实现
func (u User) TableName() string {
    return "t_user"
}

// KeyName 返回主键，必须实现
func (u User) KeyName() string {
    return "id"
}
```

1. 单行查询
```go
// 选择某个数据库，可以支持多实例
conn := sqlx.Get(ctx, "db1")
ctx := context.TODO()

var u User
err = conn.GetContext(ctx, &u, "select * from users where id = ?", id)
if err != nil {
    return
}
```

2. 多行查询
```go
// 选择某个数据库，可以支持多实例
conn := sqlx.Get(ctx, "db2")
ctx := context.TODO()

var users []User
err = conn.SelectContext(ctx, &users, "select * from users order by id desc")
if err != nil {
    return
}
```

3. `insert`
```go
// 选择某个数据库，可以支持多实例
conn := sqlx.Get(ctx, "db3")
ctx := context.TODO()

u := User{
    Name: "demo",
    Age:  18,
}
result, err := conn.InsertContext(ctx, u)
if err != nil {
    return
}

id, _ := result.LastInsertId()
```

4. `update`
```go
// 选择某个数据库，可以支持多实例
conn := sqlx.Get(ctx, "db1")
ctx := context.TODO()

u.Name = "bar"
u.ID = int(id)

_, err = conn.UpdateContext(ctx, u)
if err != nil {
    return
}
```

5. `delete`
```go
// 选择某个数据库，可以支持多实例
conn := sqlx.Get(ctx, "db2")
ctx := context.TODO()

u.ID = int(id)
_, err = conn.DeleteContext(ctx, u)
if err != nil {
    return
}
```

6. 复杂查询，对于复杂的查询，例如连表操作，需要使用类似原生的api
```go
// 选择某个数据库，可以支持多实例
conn := sqlx.Get(ctx, "db2")
ctx := context.TODO()

rows, err := conn.QueryContext(ctx, "select a.t, a.B, b.T, b.N from t_test as a left join t_demo as b on a.id = b.xx_id")
if err != nil {
   return
}

// ........
```

7. 事务
```go
func dao_func(ctx context.Context) {
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
   if _, err = tx.UpdateContext(ctx, u); err != nil {
      return
   }

   // 事务2
   if err = trans(ctx, tx); err != nil {
      return
   }

   // 事务n...

   return
}

// trans 事务中的其他操作
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
```