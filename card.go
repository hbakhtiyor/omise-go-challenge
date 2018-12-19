package main

import "time"

type Card struct {
	Name            string
	Number          string
	SecurityCode    string
	ExpirationMonth time.Month
	ExpirationYear  int
}
