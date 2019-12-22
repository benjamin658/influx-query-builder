package influxquerybuilder

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

var validUnquotedTimeRegexp = regexp.MustCompile(`^\d+(ns|u|ms|s|m|h|d|w)?$`)

// Duration Duration interface
type Duration interface {
	Nanoseconds(uint) Duration
	Microseconds(uint) Duration
	Milliseconds(uint) Duration
	Second(uint) Duration
	Minute(uint) Duration
	Hour(uint) Duration
	Day(uint) Duration
	Week(uint) Duration
	getDuration() string
}

// DurationType DurationType struct
type DurationType struct {
	unit  string
	value uint
}

// NewDuration New Duration
func NewDuration() Duration {
	return &DurationType{}
}

// Nanoseconds Nanoseconds
func (t *DurationType) Nanoseconds(d uint) Duration {
	t.unit = "ns"
	t.value = d

	return t
}

// Microseconds Microseconds
func (t *DurationType) Microseconds(d uint) Duration {
	t.unit = "u"
	t.value = d

	return t
}

// Milliseconds Milliseconds
func (t *DurationType) Milliseconds(d uint) Duration {
	t.unit = "ms"
	t.value = d

	return t
}

// Second Second
func (t *DurationType) Second(d uint) Duration {
	t.unit = "s"
	t.value = d

	return t
}

// Minute Minute
func (t *DurationType) Minute(d uint) Duration {
	t.unit = "m"
	t.value = d

	return t
}

// Hour Hour
func (t *DurationType) Hour(d uint) Duration {
	t.unit = "h"
	t.value = d

	return t
}

// Day Day
func (t *DurationType) Day(d uint) Duration {
	t.unit = "d"
	t.value = d

	return t
}

// Week Week
func (t *DurationType) Week(d uint) Duration {
	t.unit = "w"
	t.value = d

	return t
}

func (t *DurationType) getDuration() string {
	return fmt.Sprintf("time(%d%s)", t.value, t.unit)
}

// QueryBuilder QueryBuilder interface
type QueryBuilder interface {
	Select(fields ...string) QueryBuilder
	From(string) QueryBuilder
	Where(string, string, interface{}) QueryBuilder
	And(string, string, interface{}) QueryBuilder
	Or(string, string, interface{}) QueryBuilder
	WhereBrackets(QueryBuilder) QueryBuilder
	AndBrackets(QueryBuilder) QueryBuilder
	OrBrackets(QueryBuilder) QueryBuilder
	GroupBy(string) QueryBuilder
	GroupByTime(Duration) QueryBuilder
	GroupByTag(string) QueryBuilder
	Fill(interface{}) QueryBuilder
	Limit(uint) QueryBuilder
	Offset(uint) QueryBuilder
	Desc() QueryBuilder
	Asc() QueryBuilder
	Build() string
	Clean() QueryBuilder
	GetQueryStruct() CurrentQuery
}

// Tag Tag struct
type Tag struct {
	key   string
	op    string
	value interface{}
}

// Query Query struct
type Query struct {
	measurement   string
	fields        []string
	where         Tag
	and           []Tag
	or            []Tag
	whereBrackets QueryBuilder
	andBrackets   []QueryBuilder
	orBrackets    []QueryBuilder
	groupBy       string
	groupByTime   string
	groupByTag    string
	order         string
	limit         uint
	_limit        bool
	offset        uint
	_offset       bool
	fill          interface{}
}

// CurrentQuery Get current query
type CurrentQuery struct {
	Measurement   string
	Where         Tag
	And           []Tag
	Or            []Tag
	WhereBrackets QueryBuilder
	AndBrackets   []QueryBuilder
	OrBrackets    []QueryBuilder
	Fields        []string
	GroupBy       string
	GroupByTime   string
	GroupByTag    string
	Limit         uint
	Offset        uint
	Order         string
	IsLimitSet    bool
	IsOffsetSet   bool
}

// New New QueryBuilder
func New() QueryBuilder {
	return &Query{}
}

// Clean Clean current builder and get a new one
func (q *Query) Clean() QueryBuilder {
	return New()
}

// Select Select fields...
func (q *Query) Select(fields ...string) QueryBuilder {
	q.fields = append(q.fields, fields...)
	return q
}

// From From measurement
func (q *Query) From(measurement string) QueryBuilder {
	q.measurement = measurement
	return q
}

// Where Where criteria
func (q *Query) Where(key string, op string, value interface{}) QueryBuilder {
	q.where = Tag{key, op, value}
	return q
}

// And And criteria
func (q *Query) And(key string, op string, value interface{}) QueryBuilder {
	q.and = append(q.and, Tag{key, op, value})
	return q
}

// Or Or criteria
func (q *Query) Or(key string, op string, value interface{}) QueryBuilder {
	q.or = append(q.or, Tag{key, op, value})
	return q
}

// WhereBrackets WHERE (...)
func (q *Query) WhereBrackets(builder QueryBuilder) QueryBuilder {
	q.whereBrackets = builder
	return q
}

// AndBrackets AND (...)
func (q *Query) AndBrackets(builder QueryBuilder) QueryBuilder {
	q.andBrackets = append(q.andBrackets, builder)
	return q
}

// OrBrackets OR (...)
func (q *Query) OrBrackets(builder QueryBuilder) QueryBuilder {
	q.orBrackets = append(q.orBrackets, builder)
	return q
}

// GroupBy GROUP BY time
func (q *Query) GroupBy(time string) QueryBuilder {
	q.groupBy = time
	return q
}

// GroupByTime GROUP BY time
func (q *Query) GroupByTime(duration Duration) QueryBuilder {
	q.groupByTime = duration.getDuration()
	return q
}

// GroupByTag GROUP By tag
func (q *Query) GroupByTag(tag string) QueryBuilder {
	q.groupByTag = tag
	return q
}

// Fill FILL(...)
func (q *Query) Fill(fill interface{}) QueryBuilder {
	q.fill = fill
	return q
}

// Limit LIMIT x
func (q *Query) Limit(limit uint) QueryBuilder {
	q._limit = true
	q.limit = limit
	return q
}

// Offset OFFSET x
func (q *Query) Offset(offset uint) QueryBuilder {
	q._offset = true
	q.offset = offset
	return q
}

// Desc ORDER BY time DESC
func (q *Query) Desc() QueryBuilder {
	q.order = "DESC"
	return q
}

// Asc ORDER BY time ASC
func (q *Query) Asc() QueryBuilder {
	q.order = "ASC"
	return q
}

// GetQueryStruct Get query struct
func (q *Query) GetQueryStruct() CurrentQuery {
	return CurrentQuery{
		Measurement:   q.measurement,
		Where:         q.where,
		And:           q.and,
		Or:            q.or,
		WhereBrackets: q.whereBrackets,
		AndBrackets:   q.andBrackets,
		OrBrackets:    q.orBrackets,
		Fields:        q.fields,
		GroupBy:       q.groupBy,
		Limit:         q.limit,
		Offset:        q.offset,
		Order:         q.order,
		IsLimitSet:    q._limit,
		IsOffsetSet:   q._offset,
	}
}

// Build Build query string
func (q *Query) Build() string {
	var buffer bytes.Buffer

	buffer.WriteString(q.buildFields())
	buffer.WriteString(q.buildFrom())
	buffer.WriteString(q.buildWhere())
	buffer.WriteString(q.buildGroupBy())
	buffer.WriteString(q.buildFill())
	buffer.WriteString(q.buildOrder())
	buffer.WriteString(q.buildLimit())
	buffer.WriteString(q.buildOffset())

	return strings.TrimSpace(buffer.String())
}

var functionMatcher = regexp.MustCompile(`.+\(.+\)$`)
var mathMatcher = regexp.MustCompile(`^(.+?)((\s*)([\-\+\/\*])(\s*)(-?\d+)(\.\d+)?)+`)

func (q *Query) buildFields() string {
	if q.fields == nil {
		return ""
	}

	fields := make([]string, len(q.fields))

	for i := range fields {
		splitByAs := strings.Split(q.fields[i], "AS")
		selectField := strings.TrimSpace(splitByAs[0])
		selectAs := ""

		if selectField == "*" {
			return "SELECT * "
		}

		if len(splitByAs) == 2 {
			selectAs = strings.TrimSpace(splitByAs[1])
		}

		if functionMatcher.MatchString(selectField) {
			fields[i] = selectField
        } else if mathMatcher.MatchString(selectField) {
            fields[i] = selectField
		} else {
			fields[i] = fmt.Sprintf("\"%s\"", selectField)
		}

		if selectAs != "" {
			fields[i] = fields[i] + " AS " + fmt.Sprintf("\"%s\"", selectAs)
		}
	}

	return fmt.Sprintf("SELECT %s ", strings.Join(fields, ","))
}

func (q *Query) buildFrom() string {
	if q.measurement == "" {
		return ""
	}

	return fmt.Sprintf(`FROM "%s" `, q.measurement)
}

func (q *Query) buildWhere() string {
	var buffer bytes.Buffer
	var whereCriteria string
	andCriteria := make([]string, 0)
	orCriteria := make([]string, 0)

	if q.where != (Tag{}) || q.whereBrackets != nil {
		if q.where != (Tag{}) {
			buffer.WriteString("WHERE ")
			whereCriteria = getCriteriaTemplate(q.where)
			buffer.WriteString(whereCriteria)
			buffer.WriteString(" ")
		} else if q.whereBrackets != nil {
			buffer.WriteString("WHERE (")
			buffer.WriteString(strings.Replace(q.whereBrackets.Build(), "WHERE ", "", 1))
			buffer.WriteString(") ")
		}

		if q.and != nil {
			buffer.WriteString("AND ")
			for _, tag := range q.and {
				andCriteria = append(
					andCriteria,
					getCriteriaTemplate(tag),
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
					getCriteriaTemplate(tag),
				)
			}
			buffer.WriteString(strings.Join(orCriteria, " OR "))
			buffer.WriteString(" ")
		}

		if q.andBrackets != nil {
			for _, g := range q.andBrackets {
				buffer.WriteString("AND (")
				buffer.WriteString(strings.Replace(g.Build(), "WHERE ", "", 1))
				buffer.WriteString(") ")
			}
		}

		if q.orBrackets != nil {
			for _, g := range q.orBrackets {
				buffer.WriteString("OR (")
				buffer.WriteString(strings.Replace(g.Build(), "WHERE ", "", 1))
				buffer.WriteString(") ")
			}
		}
	}

	return buffer.String()
}

func (q *Query) buildGroupBy() string {
	var buffer bytes.Buffer

	if q.groupBy != "" {
		buffer.WriteString(
			fmt.Sprintf("GROUP BY time(%s)", q.groupBy),
		)

		buffer.WriteString(" ")
	} else if q.groupByTime != "" {
		buffer.WriteString(
			fmt.Sprintf("GROUP BY %s", q.groupByTime),
		)

		buffer.WriteString(" ")
	} else if q.groupByTag != "" {
		buffer.WriteString(
			fmt.Sprintf("GROUP BY %s", q.groupByTag),
		)

		buffer.WriteString(" ")
	}

	return buffer.String()
}

func (q *Query) buildFill() string {
	var buffer bytes.Buffer

	if q.fill != nil {
		buffer.WriteString(
			fmt.Sprintf(`FILL(%s)`, fmt.Sprint(q.fill)),
		)

		buffer.WriteString(" ")
	}

	return buffer.String()
}

func (q *Query) buildOrder() string {
	var buffer bytes.Buffer

	if q.order != "" {
		buffer.WriteString(
			fmt.Sprintf(`ORDER BY time %s`, q.order),
		)

		buffer.WriteString(" ")
	}

	return buffer.String()
}

func (q *Query) buildLimit() string {
	var buffer bytes.Buffer

	if q._limit {
		buffer.WriteString(
			fmt.Sprintf(`LIMIT %v`, q.limit),
		)

		buffer.WriteString(" ")
	}

	return buffer.String()
}

func (q *Query) buildOffset() string {
	var buffer bytes.Buffer

	if q._offset {
		buffer.WriteString(
			fmt.Sprintf(`OFFSET %v`, q.offset),
		)

		buffer.WriteString(" ")
	}

	return buffer.String()
}

func getCriteriaTemplate(tag Tag) string {
	switch tag.value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf(`"%s" %s %d`, tag.key, tag.op, tag.value)
	case float32, float64:
		return fmt.Sprintf(`"%s" %s %g`, tag.key, tag.op, tag.value)
	case bool:
		return fmt.Sprintf(`"%s" %s %t`, tag.key, tag.op, tag.value)
	default:
		if tag.key == "time" && validUnquotedTimeRegexp.MatchString(tag.value.(string)) {
			// 'time' key accepts non-quoted string value (eg: 1535313431000ns)
			return fmt.Sprintf(`%s %s %s`, tag.key, tag.op, tag.value)
		}
		return fmt.Sprintf(`"%s" %s '%s'`, tag.key, tag.op, tag.value)
	}
}
