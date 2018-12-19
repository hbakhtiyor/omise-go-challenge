package main

import "net/url"

func isValidURL(_url string) bool {
	if _, err := url.ParseRequestURI(_url); err != nil {
		return false
	}
	return true
}
