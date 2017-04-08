package plpr

import (
	"time"
	"regexp"
)

type Log struct {
	User string
	DBName string
	Host	string
	HostPort string
	Pid	int
	StarTime time.Time
	EndTime	time.Time
	Duration float32
	Query string
}

var formats map[string][]string = map[string][]string{
	"%u": {"User", "([0-9a-zA-Z\\.\\-\\_\\[\\]]*)?"},
	"%d": {"DBName", "([0-9a-zA-Z\\.\\-\\_\\[\\]]*)?"},
	"%h": {"Host", "(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})?"},
	"%r": {"HostPort", "(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\(\\d+\\))?"},
	"%p": {"Pid", "(\\d+)*"},
	"%t": {"Endtime", "(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2} \\D{3})"},
	"%m": {"Endtime", "(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}\\.\\d* \\D{3})"},
}

var formatRegex = regexp.MustCompile("%[udhrptm]")
