package plpr

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Log struct {
	User     string
	DBName   string
	Host     string
	HostPort string
	Pid      int64
	Time     time.Time
	Duration float64
	Query    string
}

type parser struct {
	regex        *regexp.Regexp
	usingFormats []string
	logs         []*Log
}

func Parse(data, format string) []*Log {
	formats := map[string][]string{
		"%u": {"User", "([0-9a-zA-Z\\.\\-\\_\\[\\]]*)?"},
		"%d": {"DBName", "([0-9a-zA-Z\\.\\-\\_\\[\\]]*)?"},
		"%h": {"Host", "(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})?"},
		"%r": {"HostPort", "(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\(\\d+\\))?"},
		"%p": {"Pid", "(\\d+)*"},
		"%t": {"Time", "(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2} \\D{3})"},
		"%m": {"Mtime", "(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}\\.\\d* \\D{3})"},
		"%i": {"Command", "([0-9a-zA-Z\\.\\-\\_]*)?"},
		"%c": {"SessionID", "([0-9a-f\\.]*)"},
		"%l": {"SequenceNum", "(\\d+)*"},
	}

	escapes := map[string]string{
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

	regexBase := "LOG: +duration: (\\d+\\.\\d+) ms +(execute|statement).*: (.+)?"

	matches := regexp.MustCompile("[\\.\\-\\_\\*\\?\\+\\[\\]\\{\\}]").FindAllStringSubmatch(format, -1)
	for _, v := range matches {
		format = strings.Replace(format, v[0], escapes[v[0]], 1)
	}

	var usingFormats []string
	matches = regexp.MustCompile("%[udhrptmicl]").FindAllStringSubmatch(format, -1)
	for _, v := range matches {
		format = strings.Replace(format, v[0], formats[v[0]][1], 1)
		usingFormats = append(usingFormats, formats[v[0]][0])
	}

	regex := regexp.MustCompile(format + regexBase)
	lines := regexp.MustCompile("\r\n|\r|\n").Split(data, -1)

	logs := make([]*Log, 0)
	p := &parser{
		regex:        regex,
		usingFormats: usingFormats,
		logs:         logs,
	}

	content := []string{}
	for _, v := range lines {
		if strings.Contains(v, "LOG:") {
			if len(content) == 0 {
				content = append(content, v)
				continue
			} else {
				p.parse(strings.Join(content, ""))
				content = []string{}
			}
			content = append(content, v)
		}
	}
	p.parse(strings.Join(content, ""))
	return p.logs
}

func (p *parser) parse(line string) {
	m := p.regex.FindAllStringSubmatch(line, -1)
	if len(m) == 0 {
		return
	}
	matches := m[0]

	log := &Log{}
	duration, _ := strconv.ParseFloat(matches[len(p.usingFormats)+1], 64)
	query := matches[len(p.usingFormats)+3]
	elem := reflect.ValueOf(log).Elem()

	for i := 0; i < len(p.usingFormats); i++ {
		key := p.usingFormats[i]
		value := matches[i+1]

		if key == "Time" {
			endtime, _ := time.Parse("2006-01-02 15:04:05 MST", value)
			log.Time = endtime
		} else if key == "Mtime" {
			endtime, _ := time.Parse("2006-01-02 15:04:05.000 MST", value)
			log.Time = endtime
		} else if key == "Pid" {
			pid, _ := strconv.ParseInt(value, 10, 64)
			elem.FieldByName(key).SetInt(pid)
		} else {
			elem.FieldByName(key).SetString(value)
		}
	}
	log.Duration = duration
	log.Query = query

	p.logs = append(p.logs, log)
}
