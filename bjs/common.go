package bjs

import "github.com/goccy/go-json"

func Unmarshal[T any](buf []byte) (*T, error) {
	t := new(T)
	return t, json.Unmarshal(buf, t)
}
