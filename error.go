package main

import "fmt"

type ErrInvalidCSVHeader struct {
	headers []string
}

func (e *ErrInvalidCSVHeader) Error() string {
	return fmt.Sprint("not matched headers: ", e.headers)
}
