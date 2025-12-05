package repository

import (
	"fmt"
	"strings"
)

const allFields = "*"

type Columns struct {
	readable []string
	writable []string
	alias    string
	idField  string
}

func NewColumns(readable []string, writable []string, alias, idField string) *Columns {
	return &Columns{
		readable: readable,
		writable: writable,
		alias:    alias,
		idField:  idField,
	}
}

func (c *Columns) ForInsert() []string {
	return c.writable
}

func (c *Columns) ForSelect(rawFields []string) []string {
	if len(rawFields) == 0 {
		return c.readable
	}

	fields := make([]string, 0, len(c.readable))
	for _, rawField := range rawFields {
		if rawField == allFields {
			fields = append(fields, c.readable...)
		} else if rawField != c.alias+"."+allFields {
			fields = append(fields, rawField)
		} else {
			for _, field := range c.readable {
				fields = append(fields, c.alias+"."+field)
			}
		}
	}
	return fields
}

func (c *Columns) GetIDField() string {
	return c.idField
}

func (c *Columns) GetAlias() string {
	return c.alias
}

func (c *Columns) OnConflict() string {
	if len(c.writable) == 0 || c.idField == "" {
		return ""
	}

	var statements []string
	for _, col := range c.writable {
		if col == c.idField {
			continue
		}
		statements = append(statements, fmt.Sprintf("%s = EXCLUDED.%s", col, col))
	}

	if len(statements) == 0 {
		return fmt.Sprintf("ON CONFLICT (%s) DO NOTHING", c.idField)
	}

	return fmt.Sprintf("ON CONFLICT (%s) DO UPDATE SET %s", c.idField, strings.Join(statements, ","))
}

func (c *Columns) GetOnConflictStatement() string {
	statements := make([]string, 0, len(c.writable))

	needChangeUpdatedAt := false
	for _, column := range c.writable {
		if !needChangeUpdatedAt && column == "updated_at" {
			needChangeUpdatedAt = true
			continue
		}

		statement := fmt.Sprintf("%s=EXCLUDED.%s", column, column)
		statements = append(statements, statement)
	}

	if needChangeUpdatedAt {
		statements = append(statements, "updated_at=NOW()")
	}
	return strings.Join(statements, ", ")
}
