package main

import (
	"math/rand"
	"sort"
	"strconv"
)

type Data struct {
	Metrics     []Metric
	Percentiles []Percentile
	Range       Range
}

type Percentile struct {
	Label string
	Value float64
}

type Metric struct {
	Label  string
	Values []float64
}

func RandomData() *Data {
	data := &Data{}

	data.Percentiles = []Percentile{
		{Label: "p0.1", Value: 0.001},
		{Label: "p1", Value: 0.01},
		{Label: "p5", Value: 0.05},
		{Label: "p10", Value: 0.10},
		{Label: "p25", Value: 0.25},
		{Label: "p50", Value: 0.50},
		{Label: "p75", Value: 0.75},
		{Label: "p90", Value: 0.90},
		{Label: "p95", Value: 0.95},
		{Label: "p99", Value: 0.99},
		{Label: "p99.9", Value: 0.999},
	}
	data.Metrics = make([]Metric, 32)

	for i := range data.Metrics {
		metric := &data.Metrics[i]
		metric.Label = strconv.Itoa(i)
		metric.Values = make([]float64, len(data.Percentiles))
		for i := range metric.Values {
			v := rand.ExpFloat64() * 100
			if v < 0 {
				v = -v
			}
			if v > 1000 {
				v = 1000
			}

			data.Range.Max = max(data.Range.Max, v)
			metric.Values[i] = v
		}
		sort.Float64s(metric.Values)
	}

	return data
}
