package plpr

import (
	"time"
	"regexp"
	"strings"
	"strconv"
	"reflect"
)

type Log struct {
	User string
	DBName string
	Host	string
	HostPort string
	Pid	int64
	StartTime time.Time
	EndTime	time.Time
	Duration float64
	Query string
}

var regex *regexp.Regexp
var usingFormats []string
var logs []*Log

func Parse(data, format string) []*Log {
	formats := map[string][]string{
		"%u": {"User", "([0-9a-zA-Z\\.\\-\\_\\[\\]]*)?"},
		"%d": {"DBName", "([0-9a-zA-Z\\.\\-\\_\\[\\]]*)?"},
		"%h": {"Host", "(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})?"},
		"%r": {"HostPort", "(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\(\\d+\\))?"},
		"%p": {"Pid", "(\\d+)*"},
		"%t": {"Endtime", "(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2} \\D{3})"},
		"%m": {"MEndtime", "(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}\\.\\d* \\D{3})"},
		"%i": {"Command", "([0-9a-zA-Z\\.\\-\\_]*)?"},
		"%c": {"SessionID", "([0-9a-f\\.]*)"},
		"%l": {"SequenceNum", "(\\d+)*"},
	}

	escapes := map[string]string {
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

	regexBase := "LOG:  duration: (\\d+\\.\\d+) ms  (execute|statement).*: (.+)?"

	matches := regexp.MustCompile("[\\.\\-\\_\\*\\?\\+\\[\\]\\{\\}]").FindAllStringSubmatch(format, -1)
	for _, v := range matches {
		format = strings.Replace(format, v[0], escapes[v[0]], 1)
	}

	matches = regexp.MustCompile("%[udhrptmicl]").FindAllStringSubmatch(format, -1)
	for _, v := range matches {
		format = strings.Replace(format, v[0], formats[v[0]][1], 1)
		usingFormats = append(usingFormats, formats[v[0]][0])
	}

	regex = regexp.MustCompile(format+regexBase)
	lines := regexp.MustCompile("\r\n|\r|\n").Split(data, -1)

	content := []string{}
	for _, v := range lines {
		if strings.Contains(v, ":  ") {
			if len(content) == 0 {
				content = append(content, v)
				continue
			} else {
				parse(strings.Join(content, ""))
				content = []string{}
			}
			content = append(content, v)
		}
	}
	parse(strings.Join(content, ""))
	return logs
}

func parse(line string) {
	m := regex.FindAllStringSubmatch(line, -1)
	if len(m) == 0 {
		return
	}
	matches := m[0]

	log := &Log{}
	duration, _ := strconv.ParseFloat(matches[len(usingFormats)+1], 64)
	query := matches[len(usingFormats)+3]
	elem := reflect.ValueOf(log).Elem()

	for i := 0; i < len(usingFormats); i++ {
		key := usingFormats[i]
		value := matches[i+1]

		if key == "Endtime" {
			endtime, _ := time.Parse("2006-01-02 15:04:05 MST", value)
			starttime := time.Unix(endtime.Unix() - int64(duration / 1000), 0).UTC()
			log.StartTime = starttime
			log.EndTime = starttime
		} else if key == "MEndtime" {
			endtime, _ := time.Parse("2006-01-02 15:04:05.000 MST", value)
			starttime := time.Unix(0, endtime.UnixNano() - int64(duration * 1000000)).UTC()
			log.StartTime = starttime
			log.EndTime = endtime
		} else if key == "Pid" {
			pid, _ := strconv.ParseInt(value, 10, 64)
			elem.FieldByName(key).SetInt(pid)
		} else {
			elem.FieldByName(key).SetString(value)
		}
	}
	log.Duration = duration
	log.Query = query

	logs = append(logs, log)
}
