package main

import (
	"github.com/munisystem/plpr-go"
	"io/ioutil"
	"fmt"
)

func main() {
	format := "[%m] %h:%d "
	b, err := ioutil.ReadFile("postgresql.log")
	if err != nil {
		fmt.Println(err)
	}

	body := string(b)
	logs := plpr.Parse(body, format)
	for _, v := range logs {
		fmt.Println("-----")
		fmt.Println("Time: ", v.Time)
		fmt.Println("Host: ", v.Host)
		fmt.Println("DBName: ", v.DBName)
		fmt.Println("Duration: ", v.Duration)
		fmt.Println("Query : ", v.Query)
	}
}
