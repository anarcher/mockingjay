package main

import (
	"fmt"
	"strconv"
	"time"
)

type Metric struct {
	ID        string `storm:"id"`
	Namespace string
	Name      string
	Value     float64
	CreatedAt time.Time

	AutoScalingGroupName string
}

func NewMetric(ns string, name string, value float64) *Metric {
	m := &Metric{
		Namespace: ns,
		Name:      name,
		Value:     value,
		CreatedAt: time.Now().UTC(),
	}
	m.ID = m.getID()
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
	m.ID = m.getID()

	return m, nil

}

func (m *Metric) getID() string {
	return getID(m.Namespace, m.Name, m.AutoScalingGroupName)
}

func getID(ns, name, asgname string) string {
	return fmt.Sprintf("%s:%s:%s", ns, name, asgname)
}
