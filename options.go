package main

import "strings"

type Options struct {
	Comma     rune
	Extension string
	Headers   []string
	Force     bool
}

func (o *Options) isMatchHeaders(headers []string) bool {
	return strings.Join(headers, string(o.Comma)) != strings.Join(o.Headers, string(o.Comma))
}
