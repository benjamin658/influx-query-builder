package influxqcursor

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

// QueryBuilder QueryBuilder interface
type QueryBuilder interface {
	Select(fields []string) QueryBuilder
	From(string) QueryBuilder
	Where(string, string, interface{}) QueryBuilder
	And(string, string, interface{}) QueryBuilder
	Or(string, string, interface{}) QueryBuilder
	GroupBy(string) QueryBuilder
	Fill(interface{}) QueryBuilder
	Limit(int) QueryBuilder
	Offset(int) QueryBuilder
	Desc() QueryBuilder
	Asc() QueryBuilder
	Build() string
	Clean() QueryBuilder
}

type nullValue struct{}

type tag struct {
	key   string
	op    string
	value string
}

type query struct {
	measurement string
	fields      []string
	where       tag
	and         []tag
	or          []tag
	groupBy     string
	order       string
	limit       int
	offset      int
	fill        interface{}
}

// New New QueryBuilder
func New() QueryBuilder {
	return &query{
		limit:  -1,
		offset: -1,
	}
}

// Clean Clean current builder and get a new one
func (q *query) Clean() QueryBuilder {
	return New()
}

func (q *query) Select(fields []string) QueryBuilder {
	q.fields = append(q.fields, fields...)
	return q
}

func (q *query) From(measurement string) QueryBuilder {
	q.measurement = measurement
	return q
}

func (q *query) Where(key string, op string, value interface{}) QueryBuilder {
	q.where = tag{key, op, fmt.Sprint(value)}
	return q
}

func (q *query) And(key string, op string, value interface{}) QueryBuilder {
	q.and = append(q.and, tag{key, op, fmt.Sprint(value)})
	return q
}

func (q *query) Or(key string, op string, value interface{}) QueryBuilder {
	q.or = append(q.or, tag{key, op, fmt.Sprint(value)})
	return q
}

func (q *query) GroupBy(time string) QueryBuilder {
	q.groupBy = time
	return q
}

func (q *query) Fill(fill interface{}) QueryBuilder {
	q.fill = fill
	return q
}

func (q *query) Limit(limit int) QueryBuilder {
	q.limit = limit
	return q
}

func (q *query) Offset(offset int) QueryBuilder {
	q.offset = offset
	return q
}

func (q *query) Desc() QueryBuilder {
	q.order = "DESC"
	return q
}

func (q *query) Asc() QueryBuilder {
	q.order = "ASC"
	return q
}

func (q *query) Build() string {
	var buffer bytes.Buffer
	selectStmt := q.buildFields()
	fromStmt := q.buildFrom()

	if selectStmt != "" && fromStmt != "" {
		buffer.WriteString(selectStmt)
		buffer.WriteString(fromStmt)
		buffer.WriteString(q.buildWhere())
		buffer.WriteString(q.buildGroupBy())
		buffer.WriteString(q.buildFill())
		buffer.WriteString(q.buildOrder())
		buffer.WriteString(q.buildLimit())
		buffer.WriteString(q.buildOffset())
	}

	return strings.TrimSpace(buffer.String())
}

var functionMatcher = regexp.MustCompile(`.+\(.+\)$`)

func (q *query) buildFields() string {
	if q.fields == nil {
		return ""
	}

	tmpl := `"%s"`
	fields := make([]string, len(q.fields))

	for i := range fields {
		if functionMatcher.MatchString(q.fields[i]) {
			fields[i] = q.fields[i]
		} else {
			fields[i] = fmt.Sprintf(tmpl, q.fields[i])
		}
	}

	return fmt.Sprintf("SELECT %s ", strings.Join(fields, ","))
}

func (q *query) buildFrom() string {
	if q.measurement == "" {
		return ""
	}

	return fmt.Sprintf(`FROM "%s" `, q.measurement)
}

func (q *query) buildWhere() string {
	var buffer bytes.Buffer

	if q.where != (tag{}) {
		var whereCriteria string
		andCriteria := make([]string, 0)
		orCriteria := make([]string, 0)

		buffer.WriteString("WHERE ")
		whereCriteria = fmt.Sprintf(`"%s" %s %s`, q.where.key, q.where.op, q.where.value)
		buffer.WriteString(whereCriteria)
		buffer.WriteString(" ")

		if q.and != nil {
			buffer.WriteString("AND ")
			for _, tag := range q.and {
				andCriteria = append(
					andCriteria,
					fmt.Sprintf(`"%s" %s %s`, tag.key, tag.op, tag.value),
				)
			}
			buffer.WriteString(strings.Join(andCriteria, " AND "))
			buffer.WriteString(" ")
		}

		if q.or != nil {
			buffer.WriteString("OR ")
			for _, tag := range q.or {
				orCriteria = append(
					orCriteria,
					fmt.Sprintf(`"%s" %s %s`, tag.key, tag.op, tag.value),
				)
			}
			buffer.WriteString(strings.Join(orCriteria, " OR "))
			buffer.WriteString(" ")
		}

		buffer.WriteString(" ")
	}

	return buffer.String()
}

func (q *query) buildGroupBy() string {
	var buffer bytes.Buffer

	if q.groupBy != "" {
		buffer.WriteString(
			fmt.Sprintf("GROUP BY time(%s)", q.groupBy),
		)

		buffer.WriteString(" ")
	}

	return buffer.String()
}

func (q *query) buildFill() string {
	var buffer bytes.Buffer

	if q.fill != nil {
		buffer.WriteString(
			fmt.Sprintf(`FILL(%s)`, fmt.Sprint(q.fill)),
		)

		buffer.WriteString(" ")
	}

	return buffer.String()
}

func (q *query) buildOrder() string {
	var buffer bytes.Buffer

	if q.order != "" {
		buffer.WriteString(
			fmt.Sprintf(`ORDER BY time %s`, q.order),
		)

		buffer.WriteString(" ")
	}

	return buffer.String()
}

func (q *query) buildLimit() string {
	var buffer bytes.Buffer

	if q.limit != -1 {
		buffer.WriteString(
			fmt.Sprintf(`LIMIT %v`, q.limit),
		)

		buffer.WriteString(" ")
	}

	return buffer.String()
}

func (q *query) buildOffset() string {
	var buffer bytes.Buffer

	if q.offset != -1 {
		buffer.WriteString(
			fmt.Sprintf(`OFFSET %v`, q.offset),
		)

		buffer.WriteString(" ")
	}

	return buffer.String()
}
