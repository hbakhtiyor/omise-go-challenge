package main

import (
	"bufio"
	"donate/cipher"
	"io"
)

func processCSV(reader io.ReadCloser, ch chan *Payment, options *Options) error {
	defer reader.Close()
	newReader := bufio.NewReader(reader)

	if options.Extension == ".rot128" {
		decoder, err := cipher.NewRot128Reader(newReader)
		if err != nil {
			return err
		}
		return parsePayments(decoder, ch, options)
	}

	return parsePayments(newReader, ch, options)
}
