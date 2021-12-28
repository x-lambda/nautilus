package sqlx

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx/reflectx"
)

// Modeler 接口提供查询模型的表结构信息
// 所有模型都需要实现接口
type Modeler interface {
	// TableName 返回表名
	TableName() string

	// KeyName 返回主键字段
	KeyName() string
}

// mapExecer 统一DB和Tx对象
type mapExecer interface {
	DriverName() string
	GetMapper() *reflectx.Mapper
	Rebind(string) string
	ExecContext(ctx context.Context, query string, args ...interface{}) (result sql.Result, err error)
}

func bindArgs(names []string, arg interface{}, m *reflectx.Mapper) ([]interface{}, error) {
	args := make([]interface{}, 0, len(names))

	v := reflect.ValueOf(arg)
	for v = reflect.ValueOf(arg); v.Kind() == reflect.Ptr; {
		v = v.Elem()
	}

	err := m.TraversalsByNameFunc(v.Type(), names, func(i int, t []int) error {
		if len(t) == 0 {
			return fmt.Errorf("could not find name %s in %#v", names[i], arg)
		}

		val := reflectx.FieldByIndexesReadOnly(v, t)
		args = append(args, val.Interface())

		return nil
	})

	return args, err
}

// bindModeler 解析出model的字段映射到数据库的字段，和对应的参数值
func bindModeler(arg interface{}, m *reflectx.Mapper) ([]string, []interface{}, error) {
	t := reflect.TypeOf(arg)
	names := []string{}

	for k := range m.TypeMap(t).Names {
		names = append(names, k)
	}

	args, err := bindArgs(names, arg, m)
	if err != nil {
		return nil, nil, err
	}

	return names, args, nil
}

// insert sql insert封装接口
func insert(ctx context.Context, db mapExecer, m Modeler) (result sql.Result, err error) {
	names, args, err := bindModeler(m, db.GetMapper())
	if err != nil {
		return nil, err
	}

	marks := ""
	var k int
	for i := 0; i < len(names); i++ {
		if names[i] == m.KeyName() {
			args = append(args[:i], args[i+1:]...)
			k = i
			continue
		}

		marks += "?,"
	}

	names = append(names[:k], names[k+1:]...)
	marks = marks[:len(marks)-1]
	query := "INSERT INTO " + m.TableName() + "(" + strings.Join(names, ",") + ") VALUES (" + marks + ")"

	// 将查询占位符(bindvars)转成每个db驱动识别的占位符
	// ? ----> mysql: ?
	// ? ----> sqlite: $1/?
	// ? ----> oracle: :name
	// 参考: https://www.liwenzhou.com/posts/Go/sqlx/#autoid-0-4-0
	// https://github.com/jmoiron/sqlx/blob/master/sqlx_test.go#L1319
	query = db.Rebind(query)
	return db.ExecContext(ctx, query, args...)

	return
}

// update sql update封装接口
// 全量更新，需要指定主键
func update(ctx context.Context, db mapExecer, m Modeler) (sql.Result, error) {
	names, args, err := bindModeler(m, db.GetMapper())
	if err != nil {
		return nil, err
	}

	query := "UPDATE " + m.TableName() + " SET "
	var id interface{}
	for i := 0; i < len(names); i++ {
		name := names[i]
		if name == m.KeyName() {
			id = args[i]
			args = append(args[:i], args[i+1:]...)
			continue
		}
		query += name + "=?,"
	}

	// 去除最后一个逗号
	query = query[:len(query)-1] + " WHERE " + m.KeyName() + " = ?"
	query = db.Rebind(query)
	args = append(args, id)

	return db.ExecContext(ctx, query, args...)
}

// deletex sql delete 封装接口
// 根据主键id删除，必须指定model的主键值
func deletex(ctx context.Context, db mapExecer, m Modeler) (result sql.Result, err error) {
	names, args, err := bindModeler(m, db.GetMapper())
	if err != nil {
		return nil, err
	}

	query := "DELETE FROM " + m.TableName() + " WHERE " + m.KeyName() + " = ?"
	var id interface{}
	for i := 0; i < len(names); i++ {
		// 取主键值
		if names[i] == m.KeyName() {
			id = args[i]
			break
		}
	}

	query = db.Rebind(query)
	return db.ExecContext(ctx, query, id)
}
