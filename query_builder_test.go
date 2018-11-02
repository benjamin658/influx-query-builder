package influxquerybuilder

import (
	"fmt"
	"testing"
)

func assert(t *testing.T, q string, expected string) {
	if q != expected {
		t.Error(fmt.Sprintf("Expected %s but got %s", expected, q))
	}
}

func TestClean(t *testing.T) {
	builder := New()
	var expected = ""
	builder.
		Select([]string{"temperature", "humidity"}).
		From("measurement").
		Build()

	builder = builder.Clean()
	q := builder.Build()
	assert(t, q, expected)
}

func TestSelect(t *testing.T) {
	builder := New()
	var expected = `SELECT "temperature","humidity" FROM "measurement"`
	q := builder.
		Select([]string{"temperature", "humidity"}).
		From("measurement").
		Build()

	assert(t, q, expected)
	builder = builder.Clean()
}

func TestSelectFunction(t *testing.T) {
	var expected = `SELECT MEAN("temperature"),SUM("humidity") FROM "measurement"`
	builder := New()
	q := builder.
		Select([]string{`MEAN("temperature")`, `SUM("humidity")`}).
		From("measurement").
		Build()

	assert(t, q, expected)
}

func TestWhere(t *testing.T) {
	var expected = `SELECT "temperature","humidity" FROM "measurement" WHERE "time" < '2018-11-02T09:35:25Z'`
	builder := New()
	q := builder.
		Select([]string{"temperature", "humidity"}).
		From("measurement").Where("time", "<", "2018-11-02T09:35:25Z").
		Build()

	assert(t, q, expected)
}

func TestAnd(t *testing.T) {
	var expected = `SELECT "temperature","humidity" FROM "measurement" WHERE "time" > '2018-11-01T06:33:57.503Z' AND "time" < '2018-11-02T09:35:25Z'`
	builder := New()
	q := builder.
		Select([]string{"temperature", "humidity"}).
		From("measurement").
		Where("time", ">", "2018-11-01T06:33:57.503Z").
		And("time", "<", "2018-11-02T09:35:25Z").
		Build()

	assert(t, q, expected)
}

func TestOr(t *testing.T) {
	var expected = `SELECT "temperature","humidity" FROM "measurement" WHERE "time" > '2018-11-01T06:33:57.503Z' OR "time" < '2018-11-02T09:35:25Z'`
	builder := New()
	q := builder.
		Select([]string{"temperature", "humidity"}).
		From("measurement").
		Where("time", ">", "2018-11-01T06:33:57.503Z").
		Or("time", "<", "2018-11-02T09:35:25Z").
		Build()

	assert(t, q, expected)
}

func TestWhereAndOr(t *testing.T) {
	var expected = `SELECT "temperature","humidity" FROM "measurement" WHERE "time" > '2018-11-01T06:33:57.503Z' AND "time" < '2018-11-02T09:35:25Z' OR "tag" = 't'`
	builder := New()
	q := builder.
		Select([]string{"temperature", "humidity"}).
		From("measurement").
		Where("time", ">", "2018-11-01T06:33:57.503Z").
		And("time", "<", "2018-11-02T09:35:25Z").
		Or("tag", "=", "t").
		Build()

	assert(t, q, expected)
}

func TestGroupBy(t *testing.T) {
	var expected = `SELECT "temperature","humidity" FROM "measurement" GROUP BY time(10m)`
	builder := New()
	q := builder.
		Select([]string{"temperature", "humidity"}).
		From("measurement").
		GroupBy("10m").
		Build()

	assert(t, q, expected)
}

func TestFill(t *testing.T) {
	var expected = `SELECT "temperature","humidity" FROM "measurement" FILL(1)`
	builder := New()
	q := builder.
		Select([]string{"temperature", "humidity"}).
		From("measurement").
		Fill(1).
		Build()

	assert(t, q, expected)
}

func TestOrder(t *testing.T) {
	var expected = `SELECT "temperature","humidity" FROM "measurement" ORDER BY time DESC`
	builder := New()
	q := builder.
		Select([]string{"temperature", "humidity"}).
		From("measurement").
		Desc().
		Build()

	assert(t, q, expected)
}

func TestLimitOffset(t *testing.T) {
	var expected = `SELECT "temperature","humidity" FROM "measurement" LIMIT 10 OFFSET 5`
	builder := New()
	q := builder.
		Select([]string{"temperature", "humidity"}).
		From("measurement").
		Limit(10).
		Offset(5).
		Build()

	assert(t, q, expected)
}
