# Influx Query Cursor - InfluxDB Query Builder & Cursor-Based Pagination

[![Build Status](https://travis-ci.org/benjamin658/express-inject-middleware.svg?branch=master)](https://travis-ci.org/benjamin658/influx-query-cursor.svg?branch=master)

> The lightweight InfluxDB query builder and cursor-based pagination implement in Go.

This project is still under active development.

## Query Builder Usage

### Simple query

```go
  builder := New()
  query := builder.Select([]string{"temperature", "humidity"}).From("measurement).Build()
```

### Function query

```go
  builder := New()
  query := builder.Select([]string{`MEAN("temperature")`, `SUM("humidity")`}).From("measurement").Build()
```

Output:  `SELECT MEAN("temperature"),SUM("humidity") FROM "measurement"`

### Query with criteria

```go
  builder := New()
  query := builder.Select([]string{"temperature", "humidity"}).From("measurement").Where("temperature", ">", 30).And("humidity", "<", 10).Or("humidity", ">", 20).Build()
```

Output:  `SELECT "temperature","humidity" FROM "measurement" WHERE "temperature" > 30 AND "humidity" < 10 OR "humidity" > 20`

### Group By time

```go
  builder := New()
  query := builder.Select([]string{"temperature", "humidity"}).From("measurement").GroupBy("10m").Build()
```

Output:  `SELECT "temperature","humidity" FROM "measurement" GROUP BY time(10m)`

### Order By time

```go
  builder := New()
  query := builder.Select([]string{"temperature", "humidity"}).From("measurement").Desc().Build()
```

Output: `SELECT "temperature","humidity" FROM "measurement" ORDER BY time DESC`

### Limit and Offset

```go
  builder := New()
  query := builder.Select([]string{"temperature", "humidity"}).From("measurement").Limit(10).Offset(5).Build()
```

Output: `SELECT "temperature","humidity" FROM "measurement" LIMIT 10 OFFSET 5`

### Reset builder and get a new one

```go
  builder := New()
  // some code...
  builder = builder.Clean()
```