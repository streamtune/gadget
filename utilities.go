package main

import (
	"crypto/rand"
	"encoding/json"
	"strings"
)

const shortLen = 12

func asJson(o interface{}) string {
	b, err := json.Marshal(o)
	if err != nil {
		parrot.Error("Warning", err)
		return "{}"
	}
	return string(b)
}

func TruncateID(id string) string {
	if i := strings.IndexRune(id, ':'); i >= 0 {
		id = id[i+1:]
	}
	trimTo := shortLen
	if len(id) < shortLen {
		trimTo = len(id)
	}
	return id[:trimTo]
}

func Truncate(str string) string {
	suffix := "..."
	trimTo := shortLen
	if len(str) < shortLen {
		trimTo = len(str)
		suffix = ""
	}

	return str[:trimTo] + suffix
}

func random() string {

	var dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	var bytes = make([]byte, 12)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

func tail(a []string) []string {
	if len(a) >= 2 {
		return []string(a)[1:]
	}
	return []string{}
}

func check(e error) {
	if e != nil {
		parrot.Error("Error...", e)
		return
	}
}

func fatal(e error) {
	if e != nil {
		parrot.Error("Fatal...", e)
		panic(e)
	}
}
