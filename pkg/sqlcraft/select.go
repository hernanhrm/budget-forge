package sqlcraft

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hernanhrm/budget-forge/pkg/dafi"
)

type JoinType string

const (
	InnerJoinType JoinType = "INNER JOIN"
	LeftJoinType  JoinType = "LEFT JOIN"
	RightJoinType JoinType = "RIGHT JOIN"
)

type Join struct {
	Type      JoinType
	Table     string
	Condition string
}

type SelectQuery struct {
	table                  string
	columns                []string
	requiredColumns        map[string]struct{}
	sqlColumnByDomainField map[string]string

	filters    dafi.Filters
	sorts      dafi.Sorts
	pagination dafi.Pagination

	groups []string
	joins  []Join
}

func Select(columns ...string) SelectQuery {
	return SelectQuery{
		table:           "",
		columns:         columns,
		requiredColumns: make(map[string]struct{}),
	}
}

func (s SelectQuery) From(table string) SelectQuery {
	s.table = table

	return s
}

func (s SelectQuery) Where(filters ...dafi.Filter) SelectQuery {
	s.filters = filters

	return s
}

func (s SelectQuery) OrderBy(sorts ...dafi.Sort) SelectQuery {
	s.sorts = sorts

	return s
}

func (s SelectQuery) Limit(limit uint) SelectQuery {
	s.pagination.PageSize = limit

	return s
}

func (s SelectQuery) Page(page uint) SelectQuery {
	s.pagination.PageNumber = page

	return s
}

func (s SelectQuery) RequiredColumns(columns ...string) SelectQuery {
	for _, col := range columns {
		s.requiredColumns[col] = struct{}{}
	}

	return s
}

func (s SelectQuery) SQLColumnByDomainField(sqlColumnByDomainField map[string]string) SelectQuery {
	s.sqlColumnByDomainField = sqlColumnByDomainField

	return s
}

func (s SelectQuery) InnerJoin(table, condition string) SelectQuery {
	return s.addJoin(InnerJoinType, table, condition)
}

func (s SelectQuery) LeftJoin(table, condition string) SelectQuery {
	return s.addJoin(LeftJoinType, table, condition)
}

func (s SelectQuery) RightJoin(table, condition string) SelectQuery {
	return s.addJoin(RightJoinType, table, condition)
}

func (s SelectQuery) addJoin(joinType JoinType, table, condition string) SelectQuery {
	s.joins = append(s.joins, Join{
		Type:      joinType,
		Table:     table,
		Condition: condition,
	})

	return s
}

func (s SelectQuery) ToSQL() (Result, error) {
	if len(s.columns) == 0 {
		return Result{}, ErrEmptyColumns
	}

	if len(s.sqlColumnByDomainField) > 0 {
		requiredCols := make(map[string]struct{})
		for k := range s.requiredColumns {
			requiredSQLColumn, ok := s.sqlColumnByDomainField[k]
			if !ok {
				return Result{}, ErrInvalidFieldName
			}

			requiredCols[requiredSQLColumn] = struct{}{}
		}

		s.requiredColumns = requiredCols
	}

	builder := strings.Builder{}

	builder.WriteString("SELECT ")

	if len(s.requiredColumns) == 0 {
		builder.WriteString(strings.Join(s.columns, ", "))
	} else {
		selectedCols := make([]string, 0, len(s.requiredColumns))
		for _, col := range s.columns {
			if _, ok := s.requiredColumns[col]; ok {
				selectedCols = append(selectedCols, col)
			}
		}

		if len(selectedCols) == 0 {
			builder.WriteString(strings.Join(s.columns, ", "))
		} else {
			builder.WriteString(strings.Join(selectedCols, ", "))
		}
	}

	builder.WriteString(" FROM ")
	builder.WriteString(s.table)

	for _, join := range s.joins {
		builder.WriteString(" ")
		builder.WriteString(string(join.Type))
		builder.WriteString(" ")
		builder.WriteString(join.Table)
		builder.WriteString(" ON ")
		builder.WriteString(join.Condition)
	}

	args := []any{}
	if len(s.filters) > 0 {
		whereResult, err := WhereSafe(0, s.sqlColumnByDomainField, s.filters...)
		if err != nil {
			return Result{}, err
		}
		args = append(args, whereResult.Args...)

		builder.WriteString(whereResult.SQL)
	}

	if len(s.groups) > 0 {
		groupSQL, err := BuildGroupBy(s.groups, s.sqlColumnByDomainField)
		if err != nil {
			return Result{}, err
		}

		builder.WriteString(groupSQL)
	}

	if len(s.sorts) > 0 {
		sortSQL := BuildOrderBy(s.sorts, s.sqlColumnByDomainField)

		builder.WriteString(sortSQL)
	}

	paginationSQL := BuildPagination(s.pagination)
	builder.WriteString(paginationSQL)

	return Result{
		SQL:  builder.String(),
		Args: args,
	}, nil
}

func BuildOrderBy(sorts dafi.Sorts, sqlColumnByDomainField map[string]string) string {
	if sorts.IsZero() {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString(" ORDER BY ")
	for i, sort := range sorts {
		fieldName := string(sort.Field)

		if len(sqlColumnByDomainField) > 0 {
			if sqlColumn, ok := sqlColumnByDomainField[fieldName]; ok {
				fieldName = sqlColumn
			}
		}

		builder.WriteString(fieldName)

		if sort.Type != dafi.None {
			builder.WriteString(" ")
			builder.WriteString(strings.ToUpper(string(sort.Type)))
		}

		if i < len(sorts)-1 {
			builder.WriteString(", ")
		}
	}

	return builder.String()
}

func BuildPagination(pagination dafi.Pagination) string {
	if pagination.HasPageSize() && !pagination.HasPageNumber() {
		pagination.PageNumber = 1
	}

	if pagination.IsZero() {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString(" LIMIT ")
	const maxInt = int(^uint(0) >> 1)
	if pagination.PageSize > uint(maxInt) {
		builder.WriteString(strconv.Itoa(maxInt))
	} else {
		builder.WriteString(strconv.Itoa(int(pagination.PageSize)))
	}

	if pagination.HasPageNumber() {
		builder.WriteString(" OFFSET ")
		var offset uint
		if pagination.PageNumber > 0 {
			offset = pagination.PageSize * (pagination.PageNumber - 1)
		}
		if offset > uint(maxInt) {
			builder.WriteString(strconv.Itoa(maxInt))
		} else {
			builder.WriteString(strconv.Itoa(int(offset)))
		}
	}

	return builder.String()
}

func BuildGroupBy(groups []string, sqlColumnByDomainField map[string]string) (string, error) {
	if len(sqlColumnByDomainField) > 0 {
		for i, group := range groups {
			sqlColumnName, ok := sqlColumnByDomainField[group]
			if !ok {
				return "", fmt.Errorf("invalid field name for grouping: %s", group)
			}

			groups[i] = sqlColumnName
		}
	}

	return " GROUP BY " + strings.Join(groups, ", "), nil
}
