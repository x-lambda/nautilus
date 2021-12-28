package sqlx

import (
	"context"
	"database/sql/driver"
	"time"

	"nautilus/pkg/log"

	"github.com/ngrok/sqlmw"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// observer 拦截器：观察所有sql执行情况
// 执行SQL前会回调对应的函数
// 实现自 github.com/ngrok/sqlmw::Interceptor
type observer struct {
	sqlmw.NullInterceptor
	name string
}

// ConnExecContext 执行Exec SQL
func (o observer) ConnExecContext(ctx context.Context, conn driver.ExecerContext,
	query string, args []driver.NamedValue) (result driver.Result, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Exec")
	defer span.Finish()

	ext.Component.Set(span, "sqlx")
	ext.DBInstance.Set(span, o.name)
	ext.DBStatement.Set(span, query)

	s := time.Now()
	result, err = conn.ExecContext(ctx, query, args)
	d := time.Since(s)

	log.Get(ctx).Debugf("[sqlx] name: %s exec: %s args: %v, cost: %v",
		o.name, query, values(args), d)

	table, cmd := parseSQL(query)
	sqlDurations.WithLabelValues(o.name, table, cmd).Observe(d.Seconds())

	return
}

// ConnQueryContext 执行Query SQL
func (o observer) ConnQueryContext(ctx context.Context, conn driver.QueryerContext,
	query string, args []driver.NamedValue) (rows driver.Rows, err error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "Query")
	defer span.Finish()

	ext.Component.Set(span, "sqlx")
	ext.DBInstance.Set(span, o.name)
	ext.DBStatement.Set(span, query)

	s := time.Now()
	rows, err = conn.QueryContext(ctx, query, args)
	d := time.Since(s)

	log.Get(ctx).Debugf("[sqlx] name: %s query: %s args: %v cost: %v",
		o.name, query, values(args), d)

	table, cmd := parseSQL(query)
	sqlDurations.WithLabelValues(o.name, table, cmd).Observe(d.Seconds())

	return
}

// ConnPrepareContext prepare
// mysql-driver会向MySQL发起 prepared statement请求，获取到对应的stmt后将其返回
// 参考: https://manjusaka.itscoder.com/posts/2020/01/05/simple-introdution-about-sql-prepared/
func (o observer) ConnPrepareContext(ctx context.Context, conn driver.ConnPrepareContext,
	query string) (stmt driver.Stmt, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "prepare")
	defer span.Finish()

	ext.Component.Set(span, "sqlx")
	ext.DBInstance.Set(span, o.name)
	ext.DBStatement.Set(span, query)

	s := time.Now()
	stmt, err = conn.PrepareContext(ctx, query)
	d := time.Since(s)

	log.Get(ctx).Debugf("[sqlx] name: %s prepare: %s args: %v cost: %v",
		o.name, query, nil, d)

	table, _ := parseSQL(query)
	sqlDurations.WithLabelValues(o.name, table, "prepare").Observe(d.Seconds())

	return
}

// StmtExecContext exec stmt
func (o observer) StmtExecContext(ctx context.Context, stmt driver.StmtExecContext,
	query string, args []driver.NamedValue) (result driver.Result, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PreparedExec")
	defer span.Finish()

	ext.Component.Set(span, "sqlx")
	ext.DBInstance.Set(span, o.name)
	ext.DBStatement.Set(span, query)

	s := time.Now()
	result, err = stmt.ExecContext(ctx, args)
	d := time.Since(s)

	log.Get(ctx).Debugf("[sqlx] name: %s exec stmt: %s, args: %v, cost: %v",
		o.name, query, values(args), d)

	table, cmd := parseSQL(query)
	sqlDurations.WithLabelValues(o.name, table, cmd+"-stmt").Observe(d.Seconds())

	return
}

// StmtQueryContext query stmt
func (o observer) StmtQueryContext(ctx context.Context, stmt driver.StmtQueryContext,
	query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PreparedQuery")
	defer span.Finish()

	ext.Component.Set(span, "sqlx")
	ext.DBInstance.Set(span, o.name)
	ext.DBStatement.Set(span, query)

	s := time.Now()
	rows, err = stmt.QueryContext(ctx, args)
	d := time.Since(s)

	log.Get(ctx).Debugf("[sqlx] name: %s, query stmt: %s, args: %v, cost: %v",
		o.name, query, values(args), d)

	table, cmd := parseSQL(query)
	sqlDurations.WithLabelValues(o.name, table, cmd+"-stmt").Observe(d.Seconds())

	return
}

func (o observer) ConnBeginTx(ctx context.Context, conn driver.ConnBeginTx, txOpts driver.TxOptions) (tx driver.Tx, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "begin")
	defer span.Finish()

	ext.Component.Set(span, "sqlx")
	ext.DBInstance.Set(span, o.name)

	s := time.Now()
	tx, err = conn.BeginTx(ctx, txOpts)
	d := time.Since(s)

	log.Get(ctx).Debugf("[sqlx] name: %s, begin, cost: %v", o.name, d)
	sqlDurations.WithLabelValues(o.name, "", "begin").Observe(d.Seconds())

	return
}

func (o observer) TxCommit(ctx context.Context, tx driver.Tx) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commit")
	defer span.Finish()

	ext.Component.Set(span, "sqlx")
	ext.DBInstance.Set(span, o.name)

	s := time.Now()
	err = tx.Commit()
	d := time.Since(s)

	log.Get(ctx).Debugf("[sqlx] name: %s, commit, cost: %v", o.name, d)
	sqlDurations.WithLabelValues(o.name, "", "commit").Observe(d.Seconds())

	return
}

func (o observer) TxRollback(ctx context.Context, tx driver.Tx) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Rollback")
	defer span.Finish()

	ext.Component.Set(span, "sqldb")
	ext.DBInstance.Set(span, o.name)

	s := time.Now()
	err = tx.Rollback()
	d := time.Since(s)

	log.Get(ctx).Debugf("[sqldb] name:%s, rollback, cost: %v", o.name, d)

	sqlDurations.WithLabelValues(o.name, "", "rollback").Observe(d.Seconds())

	return
}
