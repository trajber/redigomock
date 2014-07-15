// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package redigomock

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

var (
	commands map[string]*Cmd
)

func init() {
	commands = make(map[string]*Cmd)
}

// Cmd stores the registered information about a command to return it later when request by a
// command execution
type Cmd struct {
	Response interface{} // Response to send back when this command/arguments are called
	Error    error       // Error to send back when this command/arguments are called
}

// Command register a command in the mock system using the same arguments of a Do or Send commands.
// It will return a registered command object where you can set the response or error
func Command(commandName string, args ...interface{}) *Cmd {
	var cmd Cmd
	commands[generateKey(commandName, args)] = &cmd
	return &cmd
}

// GenericCommand register a command withot arguments. If a command with arguments doesn't match
// with any registered command, it will look for generic commands before throwing an error
func GenericCommand(commandName string) *Cmd {
	var cmd Cmd
	commands[generateKey(commandName, nil)] = &cmd
	return &cmd
}

// Expect sets a response for this command. Everytime a Do or Receive methods are executed for a
// registered command this response or error will be returned. You cannot set a response and a error
// for the same command/arguments
func (c *Cmd) Expect(response interface{}) {
	c.Response = response
	c.Error = nil
}

// ExpectMap works in the same way of the Expect command, but has a key/value input to make it
// easier to build test environments
func (c *Cmd) ExpectMap(response map[string]string) {
	var values []interface{}
	for key, value := range response {
		values = append(values, []byte(key))
		values = append(values, []byte(value))
	}

	c.Response = values
	c.Error = nil
}

// ExpectError allows you to force an error when executing a command/arguments
func (c *Cmd) ExpectError(err error) {
	c.Response = nil
	c.Error = err
}

// generateKey build an id for the command/arguments to make it easier to find in the registered
// commands
func generateKey(commandName string, args []interface{}) string {
	key := strings.TrimSpace(commandName)
	key = strings.ToUpper(key)

	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			key += " " + strings.TrimSpace(arg)
		case []byte:
			key += " " + strings.TrimSpace(string(arg))
		case int:
			key += " " + strconv.Itoa(arg)
		case int64:
			key += " " + strconv.FormatInt(arg, 10)
		case float64:
			key += " " + strconv.FormatFloat(arg, 'g', -1, 64)
		case bool:
			if arg {
				key += " 1"
			} else {
				key += " 0"
			}
		case nil:
			key += " "
		default:
			var buf bytes.Buffer
			fmt.Fprint(&buf, arg)
			key += " " + strings.TrimSpace(buf.String())
		}
	}

	return key
}
