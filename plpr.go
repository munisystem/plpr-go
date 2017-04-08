package plpr

import (
	"time"
	"regexp"
	"strings"
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

func Parse(format string) string {
	var formats map[string][]string = map[string][]string{
		"%u": {"User", "([0-9a-zA-Z\\.\\-\\_\\[\\]]*)?"},
		"%d": {"DBName", "([0-9a-zA-Z\\.\\-\\_\\[\\]]*)?"},
		"%h": {"Host", "(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})?"},
		"%r": {"HostPort", "(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\(\\d+\\))?"},
		"%p": {"Pid", "(\\d+)*"},
		"%t": {"Endtime", "(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2} \\D{3})"},
		"%m": {"Endtime", "(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}\\.\\d* \\D{3})"},
		"%i": {"Command", "([0-9a-zA-Z\\.\\-\\_]*)?"},
		"%c": {"SessionID", "([0-9a-f\\.]*)"},
		"%l": {"SequenceNum", "(\\d+)*"},
	}

	var escapes map[string] string = map[string]string {
		".": "\\.",
		"-": "\\-",
		"_": "\\_",
		"*": "\\*",
		"?": "\\?",
		"+": "\\+",
		"[": "\\[",
		"]": "\\]",
		"{": "\\{",
		"}": "\\}",
	}

	var escapeRegex = regexp.MustCompile("[\\.\\-\\_\\*\\?\\+\\[\\]\\{\\}]")

	matches := escapeRegex.FindAllStringSubmatch(format, -1)
	for _, v := range matches {
		format = strings.Replace(format, v[0], escapes[v[0]], 1)
	}

	var formatRegex = regexp.MustCompile("%[udhrptmicl]")
	matches = formatRegex.FindAllStringSubmatch(format, -1)
	for _, v := range matches {
		format = strings.Replace(format, v[0], formats[v[0]][1], 1)
	}

	//regex := regexp.MustCompile(format)
	return format
}

