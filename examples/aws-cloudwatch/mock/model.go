package main

import (
	//	"github.com/asdine/storm"
	"github.com/asdine/storm/q"

	"strconv"
	"time"
)

type Metric struct {
	ID        int `storm:"id,increment"`
	Namespace string
	Name      string
	Value     float64
	CreatedAt time.Time `storm:index`

	AutoScalingGroupName string
}

func NewMetric(ns string, name string, value float64) *Metric {
	m := &Metric{
		Namespace: ns,
		Name:      name,
		Value:     value,
		CreatedAt: time.Now().UTC(),
	}
	return m
}

func NewMetricString(ns, name, value string) (*Metric, error) {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, err
	}

	m := NewMetric(ns, name, v)
	return m, nil

}

func NewASGInServiceInstancesMetric(ns, asgname, value string) (*Metric, error) {

	m, err := NewMetricString(ns, "GroupInServiceInstances", value)
	if err != nil {
		return m, err
	}
	m.AutoScalingGroupName = asgname

	return m, nil

}

func MetricStartEndTimeMatcher(from, to string) (q.Matcher, error) {
	f, err := time.Parse(time.RFC3339, from)
	if err != nil {
		return nil, err
	}
	t, err := time.Parse(time.RFC3339, to)
	if err != nil {
		return nil, err
	}

	m := q.And(
		q.Gte("CreatedAt", f),
		q.Lte("CreatedAt", t),
	)
	return m, nil
}

func MetricDimMatcher(dims map[string]string) q.Matcher {
	var ms []q.Matcher
	for k, v := range dims {
		m := q.Eq(k, v)
		ms = append(ms, m)
	}

	ret := q.And(ms...)
	return ret

}

func MetricNameMatcher(metricName, namespace string) q.Matcher {
	var ms []q.Matcher
	ms = append(ms, q.Eq("Name", metricName))
	if namespace != "" {
		ms = append(ms, q.Eq("Namespace", namespace))
	}
	return q.And(ms...)
}
