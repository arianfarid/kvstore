package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Opts struct {
	FilePath string
	ValuesName string
}

func (r Opts) getValuesPath() string {
	return r.FilePath + r.ValuesName	
}
const SocketPath = "/tmp/kvstore.sock"

const (
	Err = "ERROR"
	NotFound = "NOTFOUND"
	Ok = "OK"
	Value = "VALUE"
)

var path = "./"
var valuesName = "keystore_values.kv" 

var commands = map[string]string {
	"GET": "GET",
	"PUT": "PUT",
	"DELETE": "DELETE",
	"LIST": "LIST",
}
const Tombstone = "<TOMBSTONE>"
// keys point to position in values
type Store struct {
	kvFile *os.File
	index map[string]string
}

func (s Store) Put(key string, value string) string {
	_,err := s.kvFile.WriteString(key+`\t`+value+"\n")
	if err != nil {
		return Err + ": " + err.Error()
	}
	s.index[key] = value

	return Ok
}

func (s Store) Get(key string) string {
	out, ok := s.index[key]
	if !ok {
		return NotFound
	}
	return Value + " " + out
}

func (s Store) Delete(key string) string {
	_, ok := s.index[key]
	if !ok {
		return NotFound
	}
	s.Put(key, Tombstone)
	delete(s.index, key)
	return Ok
}
func (s Store) List() []string {
	keys := make([]string, 0, len(s.index))
	for k:=range s.index {
		keys = append(keys, k)
	}
	return keys
}
var store Store


func main() {
	os.Remove(SocketPath)
	l, sock_err := net.Listen("unix", SocketPath)
	if sock_err != nil {
		panic(sock_err)
	}
	defer l.Close()
	
	var opts = Opts {
		FilePath: path,
		ValuesName: valuesName,
	}
	boot(opts)

	fmt.Println("KV daemon running on", SocketPath)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Accept err:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	
	defer c.Close()
	reader := bufio.NewReader(c)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		args := strings.SplitN(line, " ", 3)

		mainArg := args[0]
		
		var output string
		switch mainArg {
		case commands["GET"]:
			
			key := args[1]
			if len(key) < 1 {
				output = Err + ": Must provide valid key"
			}
			output = store.Get(key)
		case commands["PUT"]:
			key := args[1]
			if len(key) < 1 {
				output = Err + ": Must provide valid key"
			}
			val := args[2]
			output = store.Put(key, val)
		case commands["DELETE"]:
			key := args[1]
			if len(key) < 1 {
				output = Err + ": Must provide valid key"
			}
			output = store.Delete(key)
		case commands["LIST"]:
			keys := store.List()
			if len(keys) == 0 {
				c.Write([]byte("EMPTY\n"))
				continue
			}
			output = strings.Join(keys, " ") + "\n"
		default: 
			output = Err + ": provide valid arguments"
		}
		c.Write([]byte(output + "\n"))
	}
}

func boot(opts Opts) {
	
	keysPath := opts.getValuesPath()	
	keysFile, err := os.OpenFile(keysPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	store.index = make(map[string]string)
	store.kvFile = keysFile
	
	scanner := bufio.NewScanner(keysFile)
	for scanner.Scan() {
		row := scanner.Text()
		l := strings.Split(row, `\t`)
		k := l[0]
		v := l[1]
		if v == Tombstone {
			delete(store.index, k)
			continue
		}
		store.index[k] = v
	}
}


