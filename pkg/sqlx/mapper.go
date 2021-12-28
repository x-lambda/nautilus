package sqlx

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/dlmiddlecote/sqlstats"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/ngrok/sqlmw"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/singleflight"
)

var (
	sfg singleflight.Group
	rwl sync.RWMutex

	dbs = map[string]*DB{}
)

// DB sqlx DB 封装
type DB struct {
	*sqlx.DB
}

// Tx sqlx Tx 封装
type Tx struct {
	*sqlx.Tx
}

// Get 根据配置名字创建并返回 DB 连接池对象
// Get 是并发安全的，可以在多协程下使用
//
// DB 配置名字格式为 DB_{$name}_DSN
// DB 配置内容格式请参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name
//
// Warning: 如果model中的字段设置time.Time格式，数据库中存储了timestamp/datetime类型，scan的时候自动转换，则需要在dsn中指定参数parseTime=true
func Get(ctx context.Context, name string) *DB {
	rwl.RLock()
	if db, ok := dbs[name]; ok {
		rwl.RUnlock()
		return db
	}
	rwl.RUnlock()

	v, _, _ := sfg.Do(name, func() (interface{}, error) {
		dsn := "root:Y#*dUjC9%U%j@tcp(10.19.22.20:3306)/AAI_DEV?parseTime=true" // TODO 从conf中取
		driverName := "db-mysql:" + name
		driver := sqlmw.Driver(mysql.MySQLDriver{}, observer{name: name})

		sql.Register(driverName, driver)
		sdb := sqlx.MustOpen(driverName, dsn)

		// TODO 通过配置拿，需要注意扩容时，总的连接数不要超过2k
		sdb.SetMaxOpenConns(20)
		sdb.SetMaxIdleConns(10)
		sdb.SetConnMaxLifetime(1 * time.Hour)
		sdb.SetConnMaxLifetime(5 * time.Minute)

		db := &DB{sdb}

		rwl.Lock()
		defer rwl.Unlock()
		dbs[name] = db

		collector := sqlstats.NewStatsCollector(name, db)
		prometheus.MustRegister(collector)

		return db, nil
	})

	return v.(*DB)
}

// MustBegin 封装 sqlx.DB.MustBegin
func (db *DB) MustBegin() *Tx {
	tx := db.DB.MustBegin()
	return &Tx{tx}
}

// Beginx 封装 sqlx.DB.Beginx
func (db *DB) Beginx() (*Tx, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}

	return &Tx{tx}, nil
}

func (db *DB) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.DB.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &Tx{tx}, nil
}

// InsertContext 生成并执行 insert 语句
func (db *DB) InsertContext(ctx context.Context, m Modeler) (sql.Result, error) {
	return insert(ctx, db, m)
}

func (db *DB) Insert(m Modeler) (sql.Result, error) {
	return db.InsertContext(context.TODO(), m)
}

// UpdateContext 生成并执行 update 语句，注意必须指定主键
func (db *DB) UpdateContext(ctx context.Context, m Modeler) (sql.Result, error) {
	return update(ctx, db, m)
}

func (db *DB) Update(m Modeler) (sql.Result, error) {
	return db.UpdateContext(context.Background(), m)
}

// DeleteContext 生成并执行 delete 语句，注意必须指定主键
func (db *DB) DeleteContext(ctx context.Context, m Modeler) (sql.Result, error) {
	return deletex(ctx, db, m)
}

// GetMapper 添加 GetMapper 方法，方便与 Tx 统一
func (db *DB) GetMapper() *reflectx.Mapper {
	return db.Mapper
}

// InsertContext 生成并执行 insert 语句
func (tx *Tx) InsertContext(ctx context.Context, m Modeler) (sql.Result, error) {
	return insert(ctx, tx, m)
}

func (tx *Tx) Insert(m Modeler) (sql.Result, error) {
	return tx.InsertContext(context.Background(), m)
}

// UpdateContext 生成并执行 update 语句
func (tx *Tx) UpdateContext(ctx context.Context, m Modeler) (sql.Result, error) {
	return update(ctx, tx, m)
}

func (tx *Tx) Update(m Modeler) (sql.Result, error) {
	return tx.UpdateContext(context.Background(), m)
}

// GetMapper 添加 GetMapper 方法，方便与 DB 统一
func (tx *Tx) GetMapper() *reflectx.Mapper {
	return tx.Mapper
}
