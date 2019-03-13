package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"sort"
)

func KeyOfStringMap(m map[string]interface{}) []string {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func ValueOfStringMap(m map[string]interface{}) []interface{} {
	values := make([]interface{}, 0)
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func ForeachStringKeysInOrder(keys []string, handler func(string)) {
	sort.Strings(keys)
	for _, k := range keys {
		handler(k)
	}
}

func ReadFile(path string) (string, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func WriteNewFile(path string, content string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer f.Close()
	if err != nil {
		return err
	}
	n, err := f.WriteString(content)
	if err != nil {
		return err
	}
	if n < len(content) {
		return errors.New("write file " + path + " uncompleted!")
	}
	return nil
}
