package sqlx

import (
	"database/sql/driver"
	"regexp"
	"strings"
)

func values(args []driver.NamedValue) []driver.Value {
	values := make([]driver.Value, 0, len(args))

	for _, v := range args {
		values = append(values, v.Value)
	}

	return values
}

var sqlreg = regexp.MustCompile(`(?i)` +
	`(?P<cmd>select)\s+.+?from\s+(?P<table>\w+)\s+|` +
	`(?P<cmd>update)\s+(?P<table>\w+)\s+|` +
	`(?P<cmd>delete)\s+from\s+(?P<table>\w+)\s+|` +
	`(?P<cmd>insert)\s+into\s+(?P<table>\w+)`)

// parseSQL 提取sql中的表名和指令
func parseSQL(sql string) (table string, cmd string) {
	matches := sqlreg.FindStringSubmatch(sql)

	results := map[string]string{}
	names := sqlreg.SubexpNames()

	for i, match := range matches {
		if match != "" {
			results[names[i]] = match
		}
	}

	table = strings.ToLower(results["table"])
	cmd = strings.ToLower(results["cmd"])

	return
}
