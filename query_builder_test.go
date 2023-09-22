package influxquerybuilder

import (
	"fmt"
	"testing"
)

func assert(t *testing.T, q interface{}, expected interface{}) {
	if q != expected {
		t.Error(fmt.Sprintf("Expected %s but got %s", expected, q))
	}
}

func TestClean(t *testing.T) {
	expected := ""
	builder := New()
	builder.
		Select("temperature", "humidity").
		From("measurement").
		Build()

	builder = builder.Clean()
	q := builder.Build()
	assert(t, q, expected)
}

func TestSelect(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement"`
	builder := New()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		Build()

	assert(t, q, expected)
}

func TestSelectAs(t *testing.T) {
	expected := `SELECT "temperature" AS "temp","humidity" AS "hum" FROM "measurement"`
	builder := New()
	q := builder.
		Select("temperature AS temp", "humidity AS hum").
		From("measurement").
		Build()

	assert(t, q, expected)
}

func TestSelectAll(t *testing.T) {
	expected := `SELECT * FROM "measurement"`
	builder := New()
	q := builder.
		Select("*").
		From("measurement").
		Build()

	assert(t, q, expected)
}

func TestSelectFieldWithSpecialCharacter(t *testing.T) {
	expected := `SELECT "temperature-with-hyphen" FROM "measurement"`
	builder := New()
	q := builder.
		Select("temperature-with-hyphen").
		From("measurement").
		Build()
	assert(t, q, expected)
}

func TestSelectFunction(t *testing.T) {
	expected := `SELECT MEAN("temperature"),SUM("humidity") FROM "measurement"`
	builder := New()
	q := builder.
		Select(`MEAN("temperature")`, `SUM("humidity")`).
		From("measurement").
		Build()

	assert(t, q, expected)
}

func TestSelectFunctionAs(t *testing.T) {
	expected := `SELECT MEAN("temperature") AS "mt",SUM("humidity") AS "sh" FROM "measurement"`
	builder := New()
	q := builder.
		Select(`MEAN("temperature") AS mt`, `SUM("humidity") AS sh`).
		From("measurement").
		Build()

	assert(t, q, expected)
}

func TestFrom(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM rp_1h."measurement"`
	builder := New()
	q := builder.
		Select("temperature", "humidity").
		FromRP("rp_1h", "measurement").
		Build()

	assert(t, q, expected)
}

func TestWhere(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" WHERE "time" < '2018-11-02T09:35:25Z'`
	builder := New()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").Where("time", "<", "2018-11-02T09:35:25Z").
		Build()

	assert(t, q, expected)
}

func TestAnd(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" WHERE "time" > '2018-11-01T06:33:57.503Z' AND "time" < '2018-11-02T09:35:25Z'`
	builder := New()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		Where("time", ">", "2018-11-01T06:33:57.503Z").
		And("time", "<", "2018-11-02T09:35:25Z").
		Build()

	assert(t, q, expected)
}

func TestOr(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" WHERE "time" > '2018-11-01T06:33:57.503Z' OR "time" < '2018-11-02T09:35:25Z'`
	builder := New()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		Where("time", ">", "2018-11-01T06:33:57.503Z").
		Or("time", "<", "2018-11-02T09:35:25Z").
		Build()

	assert(t, q, expected)
}

func TestWhereAndOr(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" WHERE "time" > '2018-11-01T06:33:57.503Z' AND "time" < '2018-11-02T09:35:25Z' OR "tag" = 't'`
	builder := New()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		Where("time", ">", "2018-11-01T06:33:57.503Z").
		And("time", "<", "2018-11-02T09:35:25Z").
		Or("tag", "=", "t").
		Build()

	assert(t, q, expected)
}

func TestGroupBy(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" GROUP BY time(10m)`
	builder := New()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		GroupBy("10m").
		Build()

	assert(t, q, expected)
}

func TestGroupByTimeNanoSec(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" GROUP BY time(10ns)`
	builder := New()
	duration := NewDuration()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		GroupByTime(duration.Nanoseconds(10)).
		Build()

	assert(t, q, expected)
}

func TestGroupByTimeMicroSec(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" GROUP BY time(10u)`
	builder := New()
	duration := NewDuration()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		GroupByTime(duration.Microseconds(10)).
		Build()

	assert(t, q, expected)
}

func TestGroupByTimeMillSec(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" GROUP BY time(10ms)`
	builder := New()
	duration := NewDuration()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		GroupByTime(duration.Milliseconds(10)).
		Build()

	assert(t, q, expected)
}

func TestGroupByTimeSec(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" GROUP BY time(10s)`
	builder := New()
	duration := NewDuration()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		GroupByTime(duration.Second(10)).
		Build()

	assert(t, q, expected)
}

func TestGroupByTimeMinute(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" GROUP BY time(10m)`
	builder := New()
	duration := NewDuration()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		GroupByTime(duration.Minute(10)).
		Build()

	assert(t, q, expected)
}

func TestGroupByTimeHour(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" GROUP BY time(10h)`
	builder := New()
	duration := NewDuration()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		GroupByTime(duration.Hour(10)).
		Build()

	assert(t, q, expected)
}

func TestGroupByTimeDay(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" GROUP BY time(10d)`
	builder := New()
	duration := NewDuration()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		GroupByTime(duration.Day(10)).
		Build()

	assert(t, q, expected)
}

func TestGroupByTimeWeek(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" GROUP BY time(10w)`
	builder := New()
	duration := NewDuration()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		GroupByTime(duration.Week(10)).
		Build()

	assert(t, q, expected)
}

func TestGroupByTag(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" GROUP BY sensorId`
	builder := New()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		GroupByTag("sensorId").
		Build()

	assert(t, q, expected)
}

func TestGroupByTags(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" GROUP BY sensorId,location`
	builder := New()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		GroupByTag("sensorId", "location").
		Build()

	assert(t, q, expected)

	q = New().
		Select("temperature", "humidity").
		From("measurement").
		GroupByTag("sensorId").
		GroupByTag("location").
		Build()

	assert(t, q, expected)

	q = New().
		Select("temperature", "humidity").
		From("measurement").
		GroupByTime(NewDuration().Minute(5)).
		GroupByTag("sensorId", "location").
		Build()
	assert(t, q,
		`SELECT "temperature","humidity" FROM "measurement" GROUP BY time(5m),sensorId,location`)
}

func TestFill(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" FILL(1)`
	builder := New()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		Fill(1).
		Build()

	assert(t, q, expected)
}

func TestAscOrder(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" ORDER BY time ASC`
	builder := New()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		Asc().
		Build()

	assert(t, q, expected)
}

func TestDescOrder(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" ORDER BY time DESC`
	builder := New()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		Desc().
		Build()

	assert(t, q, expected)
}

func TestLimitOffset(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" LIMIT 10 OFFSET 5`
	builder := New()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		Limit(10).
		Offset(5).
		Build()

	assert(t, q, expected)
}

func TestBracketsWhere(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" WHERE ("time" > '2018-11-01T06:33:57.503Z' AND "time" < '2018-11-02T09:35:25Z') OR "tag" = 't'`
	builder := New()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		WhereBrackets(
			New().
				Where("time", ">", "2018-11-01T06:33:57.503Z").
				And("time", "<", "2018-11-02T09:35:25Z"),
		).
		Or("tag", "=", "t").
		Build()

	assert(t, q, expected)
}

func TestBracketsAndCriteria(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" WHERE "time" > '2018-11-01T06:33:57.503Z' AND ("time" < '2018-11-02T09:35:25Z' OR "tag" = 't')`
	builder := New()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		Where("time", ">", "2018-11-01T06:33:57.503Z").
		AndBrackets(
			New().
				Where("time", "<", "2018-11-02T09:35:25Z").
				Or("tag", "=", "t"),
		).
		Build()

	assert(t, q, expected)
}

func TestBracketsOrCriteria(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" WHERE "time" > '2018-11-01T06:33:57.503Z' OR ("time" < '2018-11-02T09:35:25Z' OR "tag" = 't')`
	builder := New()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		Where("time", ">", "2018-11-01T06:33:57.503Z").
		OrBrackets(
			New().
				Where("time", "<", "2018-11-02T09:35:25Z").
				Or("tag", "=", "t"),
		).
		Build()

	assert(t, q, expected)
}

func TestWhereTypeSqlQuote(t *testing.T) {
	expected := `SELECT "temperature","humidity" FROM "measurement" WHERE "temperature" > 20 OR "humidity" < 10.101`
	builder := New()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		Where("temperature", ">", 20).
		Or("humidity", "<", 10.101).
		Build()
	assert(t, q, expected)

	expected = `SELECT "temperature" FROM "measurement" WHERE "hot" = true`
	q = builder.
		Clean().
		Select("temperature").
		From("measurement").
		Where("hot", "=", true).
		Build()
	assert(t, q, expected)
}

func TestGetQueryStruct(t *testing.T) {
	var expected uint = 100
	builder := New()
	q := builder.
		Select("temperature", "humidity").
		From("measurement").
		Limit(100).
		Offset(100).
		Asc().
		GetQueryStruct()

	assert(t, q.Fields[0], "temperature")
	assert(t, q.Fields[1], "humidity")
	assert(t, q.Measurement, "measurement")
	assert(t, q.Limit, expected)
	assert(t, q.IsLimitSet, true)
	assert(t, q.Offset, expected)
	assert(t, q.IsOffsetSet, true)
	assert(t, q.Order, "ASC")
}
