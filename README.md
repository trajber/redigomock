redigomock
==========

[![Build Status](https://travis-ci.org/rafaeljusto/redigomock.png?branch=master)](https://travis-ci.org/rafaeljusto/redigomock)
[![GoDoc](https://godoc.org/github.com/rafaeljusto/redigomock?status.png)](https://godoc.org/github.com/rafaeljusto/redigomock)

Easy way to unit test projects using [redigo library](https://github.com/garyburd/redigo) (Redis client in go). You can find the latest release [here](https://github.com/rafaeljusto/redigomock/releases).

install
-------

```
go get -u github.com/rafaeljusto/redigomock
```

usage
-----

Here is an example of using redigomock, for more information please check the [API documentation](https://godoc.org/github.com/rafaeljusto/redigomock).

```go
package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"
)

type Person struct {
	Name string `redis:"name"`
	Age  int    `redis:"age"`
}

func RetrievePerson(conn redis.Conn, id string) (Person, error) {
	var person Person

	values, err := redis.Values(conn.Do("HGETALL", fmt.Sprintf("person:%s", id)))
	if err != nil {
		return person, err
	}

	err = redis.ScanStruct(values, &person)
	return person, err
}

func main() {
	// Simulate command result

	redigomock.Command("HGETALL", "person:1").ExpectMap(map[string]string{
		"name": "Mr. Johson",
		"age":  "42",
	})

	person, err := RetrievePerson(redigomock.NewConn(), "1")
	if err != nil {
		fmt.Println(err)
		return
	}

	if person.Name != "Mr. Johson" {
		fmt.Printf("Invalid name. Expected 'Mr. Johson' and got '%s'\n", person.Name)
		return
	}

	if person.Age != 42 {
		fmt.Printf("Invalid age. Expected '42' and got '%d'\n", person.Age)
		return
	}

	// Simulate command error

	redigomock.Command("HGETALL", "person:1").ExpectError(fmt.Errorf("Simulate error!"))

	person, err = RetrievePerson(redigomock.NewConn(), "1")
	if err == nil {
		fmt.Println("Should return an error!")
		return
	}

	fmt.Println("Success!")
}
```
