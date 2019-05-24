# Influx Query Builder

[![Build Status](https://travis-ci.org/benjamin658/influx-query-builder.svg?branch=master)](https://travis-ci.org/benjamin658/influx-query-builder.svg?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/benjamin658/influx-query-builder/badge.svg?branch=master)](https://coveralls.io/github/benjamin658/influx-query-builder?branch=master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/dbdc74b92709494b839a0e0e72d1f6a6)](https://app.codacy.com/app/benjamin658/influx-query-builder?utm_source=github.com&utm_medium=referral&utm_content=benjamin658/influx-query-builder&utm_campaign=Badge_Grade_Dashboard)
[![Go Report Card](https://goreportcard.com/badge/github.com/benjamin658/influx-query-builder)](https://goreportcard.com/report/github.com/benjamin658/influx-query-builder)

> The super lightweight InfluxDB query builder implemented in Go.

## Installation

`go get -u github.com/benjamin658/influx-query-builder`

## Query Builder Usage

### Simple query

```go
builder := New()
query := builder.
  Select("temperature", "humidity").
  From("measurement).
  Build()
```

Output:

```sql
SELECT temperature,humidity FROM "measurement"
```

### Function query

```go
builder := New()
query := builder.
  Select(`MEAN("temperature")`, `SUM("humidity")`).
  From("measurement").
  Build()
```

Output:

```sql
SELECT MEAN("temperature"),SUM("humidity") FROM "measurement"
```

### Query with criteria

```go
builder := New()
query := builder.
  Select("temperature", "humidity").
  From("measurement").
  Where("time", ">", "2018-11-01T06:33:57.503Z").
  And("time", "<", "2018-11-02T09:35:25Z").
  Or("tag", "=", "t").
  Build()
```

Output:

```sql
SELECT temperature,humidity FROM "measurement" WHERE "time" > '2018-11-01T06:33:57.503Z' AND "time" < '2018-11-02T09:35:25Z' OR "tag" = 't'
```

### Brackets criteria

Noted: If you use `Where` with `WhereBrackets`, `Where` will override the `WhereBrackets`.

#### Where Brackets

```go
builder := New()
query := builder.
  Select("temperature", "humidity").
  From("measurement").
  WhereBrackets(
    // Passing a new builder as the param
    New().
      Where("time", ">", "2018-11-01T06:33:57.503Z").
      And("time", "<", "2018-11-02T09:35:25Z").
  ).
  Or("tag", "=", "t").
  Build()
```

Output:

```sql
SELECT temperature,humidity FROM "measurement" WHERE ("time" > '2018-11-01T06:33:57.503Z' AND "time" < '2018-11-02T09:35:25Z') OR "tag" = 't'
```

#### And Brackets

```go
builder := New()
query := builder.
  Select("temperature", "humidity").
  From("measurement").
  Where("time", ">", "2018-11-01T06:33:57.503Z").
  AndBrackets(
    // Passing a new builder as the param
    New().
      Where("time", "<", "2018-11-02T09:35:25Z").
      Or("tag", "=", "t"),
  ).
  Build()
```

Output:

```sql
SELECT temperature,humidity FROM "measurement" WHERE "time" > '2018-11-01T06:33:57.503Z' AND ("time" < '2018-11-02T09:35:25Z' OR "tag" = 't')
```

#### Or Brackets

```go
builder := New()
query := builder.
  Select("temperature", "humidity").
  From("measurement").
  Where("time", ">", "2018-11-01T06:33:57.503Z").
  OrBrackets(
    // Passing a new builder as the param
    New().
      Where("time", "<", "2018-11-02T09:35:25Z").
      Or("tag", "=", "t"),
  ).
  Build()
```

Output:

```sql
SELECT temperature,humidity FROM "measurement" WHERE "time" > '2018-11-01T06:33:57.503Z' OR ("time" < '2018-11-02T09:35:25Z' OR "tag" = 't')
```

### Group By time

```go
builder := New()
duration := NewDuration()
query := builder.
  Select("temperature", "humidity").
  From("measurement").
  GroupByTime(duration.Minute(10)).
  Build()
```

Output:

```sql
SELECT temperature,humidity FROM "measurement" GROUP BY time(10m)
```

### Group By Tag

```go
builder := New()
duration := NewDuration()
query := builder.
  Select("temperature", "humidity").
  From("measurement").
  GroupByTag("sensorId").
  Build()
```

Output:

```sql
SELECT temperature,humidity FROM "measurement" GROUP BY sensorId
```

### Order By time

```go
builder := New()
query := builder.
  Select("temperature", "humidity").
  From("measurement").
  Desc().
  Build()
```

Output:

```sql
SELECT temperature,humidity FROM "measurement" ORDER BY time DESC
```

### Limit and Offset

```go
builder := New()
query := builder.
  Select("temperature", "humidity").
  From("measurement").
  Limit(10).
  Offset(5).
  Build()
```

Output:

```sql
SELECT temperature,humidity FROM "measurement" LIMIT 10 OFFSET 5
```

### Reset builder and get a new one

```go
builder := New()
// some code...
builder = builder.Clean()
```

### Get current query struct

```go
builder := New()
query := builder.
  Select("temperature", "humidity").
  From("measurement").
  Limit(100).
  Offset(100).
  Asc().
  GetQueryStruct()

/*
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
  Limit         uint
  Offset        uint
  Order         string
  IsLimitSet    bool
  IsOffsetSet   bool
}
*/
```

## Deprecated

### [Deprecated] Group By time

```go
builder := New()
query := builder.
  Select("temperature", "humidity").
  From("measurement").
  GroupBy("10m").
  Build()
```

Output:

```sql
SELECT temperature,humidity FROM "measurement" GROUP BY time(10m)
```

## License

-------

© Ben Hu (benjamin658), 2018-NOW

Released under the [MIT License](https://github.com/benjamin658/influx-query-builder/blob/master/LICENSE)