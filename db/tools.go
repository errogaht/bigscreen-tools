package db

import (
	"fmt"
	"strconv"
	"strings"
)

type TableMetadata struct {
	Name string
	Cols []string
	PK   string
}

func (m *TableMetadata) Comma() string {
	return strings.Join(m.Cols, ", ")
}

func (m *TableMetadata) GetUpsertSql() string {
	return fmt.Sprintf("insert into %s (%s) values(%s) on conflict (%s) do update set %s", m.Name, m.Comma(), m.Params(), m.PK, m.DoUpdate())
}

func (m *TableMetadata) GetFindBySql(cond string) (sql string) {
	sql = fmt.Sprintf("SELECT %s FROM %s", m.Comma(), m.Name)
	if cond != "" {
		sql += " WHERE " + cond
	}
	return
}
func (m *TableMetadata) GetINPreparedString(p []string) (str string) {
	if p == nil {
		return ""
	}
	for _, i := range p {
		str += `'` + i + `', `
	}
	str = str[:len(str)-2]
	return
}

func (m *TableMetadata) Params() string {
	params := make([]string, len(m.Cols))
	for i, _ := range m.Cols {
		params[i] = "$" + strconv.Itoa(i+1)
	}
	return strings.Join(params, ", ")
}

func (m *TableMetadata) Params2(par []interface{}) string {
	params := make([]string, len(par))
	for i, _ := range par {
		params[i] = "$" + strconv.Itoa(i+1)
	}
	return strings.Join(params, ", ")
}
func (m *TableMetadata) DoUpdate() string {
	var params []string
	for _, col := range m.Cols {
		if col == m.PK {
			continue
		}
		//created_at = excluded.created_at
		params = append(params, col+" = excluded."+col)
	}
	return strings.Join(params, ", ")
}
