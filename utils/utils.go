package utils

import (
	"io/ioutil"
	"os"
	"sort"
)

func KeyOfStringMap(m map[string]interface{}) []string {
	keys := make([]string, 0)
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}

func ForeachStringKeysInOrder(keys []string, handler func(string)) {
	sort.Strings(keys)
	for _, k := range keys {
		handler(k)
	}
}

func ReadFile(path string) string {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}
	return string(bytes)
}

func WriteNewFile(path string, content string) {
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		panic(err.Error())
	}
	n, err := f.WriteString(content)
	if err != nil {
		panic(err.Error())
	}
	if n < len(content) {
		panic("write file " + path + " uncompleted!")
	}
}


