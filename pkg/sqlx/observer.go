package sqlx

import (
	"context"
	"database/sql/driver"
	"time"

	"nautilus/pkg/log"
	"nautilus/pkg/trace"

	"github.com/ngrok/sqlmw"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
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
	tr := otel.Tracer("MySQL-Operation")
	ctx, span := tr.Start(ctx, "Exec")

	span.SetAttributes(trace.DBSystemValue)
	span.SetAttributes(trace.DBNameKey.String(o.name))
	span.SetAttributes(trace.DBStatementKey.String(query))

	s := time.Now()
	result, err = conn.ExecContext(ctx, query, args)
	d := time.Since(s)

	// log.Get(ctx).Debugf("[sqlx] name: %s exec: %s args: %v, cost: %v",
	//	o.name, query, values(args), d)

	table, cmd := parseSQL(query)
	sqlDurations.WithLabelValues(o.name, table, cmd).Observe(d.Seconds())

	span.SetAttributes(trace.DBOperationKey.String(cmd))
	span.SetAttributes(trace.DBTableKey.String(table))
	onSpanErr(span, err)
	return
}

// ConnQueryContext 执行Query SQL
func (o observer) ConnQueryContext(ctx context.Context, conn driver.QueryerContext,
	query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	// tracing
	tr := otel.Tracer("MySQL-Operation")
	ctx, span := tr.Start(ctx, "Query")

	span.SetAttributes(trace.DBSystemValue)
	span.SetAttributes(trace.DBNameKey.String(o.name))
	span.SetAttributes(trace.DBStatementKey.String(query))

	s := time.Now()
	rows, err = conn.QueryContext(ctx, query, args)
	d := time.Since(s)

	// log.Get(ctx).Debugf("[sqlx] name: %s query: %s args: %v cost: %v",
	//	o.name, query, values(args), d)

	table, cmd := parseSQL(query)
	sqlDurations.WithLabelValues(o.name, table, cmd).Observe(d.Seconds())

	span.SetAttributes(trace.DBOperationKey.String(cmd))
	span.SetAttributes(trace.DBTableKey.String(table))
	onSpanErr(span, err)
	return
}

// ConnPrepareContext prepare
// mysql-driver会向MySQL发起 prepared statement请求，获取到对应的stmt后将其返回
// 参考: https://manjusaka.itscoder.com/posts/2020/01/05/simple-introdution-about-sql-prepared/
func (o observer) ConnPrepareContext(ctx context.Context, conn driver.ConnPrepareContext,
	query string) (stmt driver.Stmt, err error) {
	tr := otel.Tracer("MySQL-Operation")
	ctx, span := tr.Start(ctx, "Prepare")

	span.SetAttributes(trace.DBSystemValue)
	span.SetAttributes(trace.DBNameKey.String(o.name))
	span.SetAttributes(trace.DBStatementKey.String(query))

	s := time.Now()
	stmt, err = conn.PrepareContext(ctx, query)
	d := time.Since(s)

	// log.Get(ctx).Debugf("[sqlx] name: %s prepare: %s args: %v cost: %v",
	//	o.name, query, nil, d)

	table, cmd := parseSQL(query)
	sqlDurations.WithLabelValues(o.name, table, "prepare").Observe(d.Seconds())

	span.SetAttributes(trace.DBOperationKey.String(cmd))
	span.SetAttributes(trace.DBTableKey.String(table))
	onSpanErr(span, err)
	return
}

// StmtExecContext exec stmt
func (o observer) StmtExecContext(ctx context.Context, stmt driver.StmtExecContext,
	query string, args []driver.NamedValue) (result driver.Result, err error) {
	tr := otel.Tracer("MySQL-Operation")
	ctx, span := tr.Start(ctx, "StmtExec")

	span.SetAttributes(trace.DBSystemValue)
	span.SetAttributes(trace.DBNameKey.String(o.name))
	span.SetAttributes(trace.DBStatementKey.String(query))

	s := time.Now()
	result, err = stmt.ExecContext(ctx, args)
	d := time.Since(s)

	log.Get(ctx).Debugf("[sqlx] name: %s exec stmt: %s, args: %v, cost: %v",
		o.name, query, values(args), d)

	table, cmd := parseSQL(query)
	sqlDurations.WithLabelValues(o.name, table, cmd+"-stmt").Observe(d.Seconds())

	span.SetAttributes(trace.DBOperationKey.String(cmd))
	span.SetAttributes(trace.DBTableKey.String(table))
	onSpanErr(span, err)
	return
}

// StmtQueryContext query stmt
func (o observer) StmtQueryContext(ctx context.Context, stmt driver.StmtQueryContext,
	query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	tr := otel.Tracer("MySQL-Operation")
	ctx, span := tr.Start(ctx, "StmtQuery")

	span.SetAttributes(trace.DBSystemValue)
	span.SetAttributes(trace.DBNameKey.String(o.name))
	span.SetAttributes(trace.DBStatementKey.String(query))

	s := time.Now()
	rows, err = stmt.QueryContext(ctx, args)
	d := time.Since(s)

	log.Get(ctx).Debugf("[sqlx] name: %s, query stmt: %s, args: %v, cost: %v",
		o.name, query, values(args), d)

	table, cmd := parseSQL(query)
	sqlDurations.WithLabelValues(o.name, table, cmd+"-stmt").Observe(d.Seconds())

	span.SetAttributes(trace.DBOperationKey.String(cmd))
	span.SetAttributes(trace.DBTableKey.String(table))
	onSpanErr(span, err)
	return
}

func (o observer) ConnBeginTx(ctx context.Context, conn driver.ConnBeginTx, txOpts driver.TxOptions) (tx driver.Tx, err error) {
	tr := otel.Tracer("MySQL-Operation")
	ctx, span := tr.Start(ctx, "trans")

	span.SetAttributes(trace.DBSystemValue)
	span.SetAttributes(trace.DBNameKey.String(o.name))
	span.SetAttributes(trace.DBStatementKey.String("begin"))

	s := time.Now()
	tx, err = conn.BeginTx(ctx, txOpts)
	d := time.Since(s)

	log.Get(ctx).Debugf("[sqlx] name: %s, begin, cost: %v", o.name, d)
	sqlDurations.WithLabelValues(o.name, "", "begin").Observe(d.Seconds())
	onSpanErr(span, err)
	return
}

func (o observer) TxCommit(ctx context.Context, tx driver.Tx) (err error) {
	tr := otel.Tracer("MySQL-Operation")
	ctx, span := tr.Start(ctx, "trans")

	span.SetAttributes(trace.DBSystemValue)
	span.SetAttributes(trace.DBNameKey.String(o.name))
	span.SetAttributes(trace.DBStatementKey.String("commit"))

	s := time.Now()
	err = tx.Commit()
	d := time.Since(s)

	log.Get(ctx).Debugf("[sqlx] name: %s, commit, cost: %v", o.name, d)
	sqlDurations.WithLabelValues(o.name, "", "commit").Observe(d.Seconds())
	onSpanErr(span, err)
	return
}

func (o observer) TxRollback(ctx context.Context, tx driver.Tx) (err error) {
	tr := otel.Tracer("MySQL-Operation")
	ctx, span := tr.Start(ctx, "trans")

	span.SetAttributes(trace.DBSystemValue)
	span.SetAttributes(trace.DBNameKey.String(o.name))
	span.SetAttributes(trace.DBStatementKey.String("rollback"))

	s := time.Now()
	err = tx.Rollback()
	d := time.Since(s)

	log.Get(ctx).Debugf("[sqldb] name:%s, rollback, cost: %v", o.name, d)

	sqlDurations.WithLabelValues(o.name, "", "rollback").Observe(d.Seconds())
	onSpanErr(span, err)
	return
}

// onSpanErr 记录span err
func onSpanErr(span oteltrace.Span, err error) {
	defer span.End()

	if err == nil {
		return
	}

	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}
