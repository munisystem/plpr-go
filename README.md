# plpr-go
plpr-go is PostgreSQL log parser library.

# Install
```
$ go get -d munisystem/plpr-go
```

## Usage
```go
package main

import (
	"github.com/munisystem/plpr-go"
	"fmt"
)

func main() {
	format := "[%m] %h:%d "
	body := `[2007-09-01 16:44:49.244 ADT] 192.168.2.10:ossecdb LOG:  duration: 4.550 ms  statement: SELECT id FROM location WHERE name = 'enigma->/var/log/messages' AND server_id = '1'`
	logs := plpr.Parse(body, format)
	fmt.Println(logs[0].Query)
	//=> "SELECT id FROM location WHERE name = 'enigma->/var/log/messages' AND server_id = '1'"
}
```

# License
MIT © munisystem
