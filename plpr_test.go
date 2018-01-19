package plpr

import (
	"fmt"
	"reflect"
	"testing"
)

const format = "[%m] %h:%d "

var TestParseCase map[string]interface{} = map[string]interface{}{
	"data": `[2007-09-01 16:44:49.244 ADT] 192.168.2.10:testdb LOG: duration: 4.550 ms statement: SELECT id FROM users WHERE name = 'alice';
[2007-09-01 16:44:49.251 ADT] 192.168.2.10:testdb LOG: checkpoint starting: time
`,
	"expect": []*Log{
		&Log{
			Host:   "192.168.2.10",
			DBName: "testdb",
			Query:  "SELECT id FROM users WHERE name = 'alice';",
		},
	},
}

func TestParse(t *testing.T) {
	expect := TestParseCase["expect"].([]*Log)
	actual := Parse(TestParseCase["data"].(string), format)
	fmt.Println(actual)
	if len(expect) != len(actual) {
		t.Fatalf("didn't match result length, expect: %v, actual: %v", len(expect), len(actual))
	}
	for i, log := range actual {
		expect[i].Time = log.Time
		if reflect.DeepEqual(expect[i], log) {
			t.Errorf("didn't match result, expect: %v, actual: %v", expect[i], log)
		}
	}
}

var TestParseWithMultiLineQueryCase map[string]interface{} = map[string]interface{}{
	"data": `[2007-09-01 16:44:49.244 ADT] 192.168.2.10:testdb LOG: duration: 4.550 ms statement:
SELECT id
FROM users
WHERE name = 'alice';
[2007-09-01 16:44:49.251 ADT] 192.168.2.10:testdb LOG: duration: 5.252 ms statement: INSERT INTO users(name, age) VALUES ('bob', 24);
`,
	"expect": []*Log{
		&Log{
			Host:   "192.168.2.10",
			DBName: "testdb",
			Query:  "SELECT id FROM users WHERE name = 'alice';",
		},
		&Log{
			Host:   "192.168.2.10",
			DBName: "testdb",
			Query:  "INSERT INTO users(name, age) VALUES ('bob', 24);",
		},
	},
}

func TestParseWithMultiLineQuery(t *testing.T) {
	expect := TestParseWithMultiLineQueryCase["expect"].([]*Log)
	actual := Parse(TestParseWithMultiLineQueryCase["data"].(string), format)
	if len(expect) != len(actual) {
		t.Fatalf("didn't match result length, expect: %v, actual: %v", len(expect), len(actual))
	}
	for i, log := range actual {
		expect[i].Time = log.Time
		if reflect.DeepEqual(expect[i], log) {
			t.Errorf("didn't match result, expect: %v, actual: %v", expect[i], log)
		}
	}
}

var TestParseWithOldFormatCase map[string]interface{} = map[string]interface{}{
	"data": `[2007-09-01 16:44:49.244 ADT] 192.168.2.10:testdb LOG:  duration: 4.550 ms  statement: SELECT id FROM users WHERE name = 'alice';
`,
	"expect": []*Log{
		&Log{
			Host:   "192.168.2.10",
			DBName: "testdb",
			Query:  "SELECT id FROM users WHERE name = 'alice';",
		},
	},
}

func TestParseWithOldFormat(t *testing.T) {
	expect := TestParseWithOldFormatCase["expect"].([]*Log)
	actual := Parse(TestParseWithOldFormatCase["data"].(string), format)
	if len(expect) != len(actual) {
		t.Fatalf("didn't match result length, expect: %v, actual: %v", len(expect), len(actual))
	}
	for i, log := range actual {
		expect[i].Time = log.Time
		if reflect.DeepEqual(expect[i], log) {
			t.Errorf("didn't match result, expect: %v, actual: %v", expect[i], log)
		}
	}
}
