package main

import (
	"bufio"
	"fmt"
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

var path = "./"
var valuesName = "keystore_values.kv" 

var commands = map[string]string {
	"GET": "GET",
	"PUT": "PUT",
	"DELETE": "DELETE",
}
const Tombstone = "<TOMBSTONE>"
// keys point to position in values
type Store struct {
	kvFile *os.File
	index map[string]string
}
var store Store

func main() {
	
	var opts = Opts {
		FilePath: path,
		ValuesName: valuesName,
	}

	args := os.Args[1:]
	mainArg := args[0]

	boot(opts)
	switch mainArg {
	case commands["GET"]:
		key := args[1]
		Get(key)
	case commands["PUT"]:
		key := args[1]
		val := args[2]
		Put(key, val)
	case commands["DELETE"]:
		key := args[1]
		Delete(key)
	default: 
		fmt.Println("Invalid argument supplied")
	}
}


func boot(opts Opts) {
	
	keysPath := opts.getValuesPath()	
	keysFile, err := os.OpenFile(keysPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("ERR: ", err)
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


func Put(key string, value string) {
	store.kvFile.WriteString(key+`\t`+value+"\n")
	store.index[key] = value
}

func Get(key string) {
	out, ok := store.index[key]
	if !ok {
		fmt.Println("No Key for", out)	
	} else {
		fmt.Println(key, out)
	}
}

func Delete(key string) {
	Put(key, Tombstone)
	delete(store.index, key)
}
