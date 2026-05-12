package csp

import "time"

type GoJqFilter struct {
	Name string
	GoJq string
}

type Package struct {
	StartTime time.Time
	EndTime   time.Time
	Input     any
	Output    any
	Filename  string
	Filters   []GoJqFilter
	Index     int
}
