package main

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"strconv"
	"time"
)

type Payment struct {
	Amount uint64
	Card   *Card
}

func parsePayments(reader io.Reader, ch chan *Payment, options *Options) error {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = len(options.Headers)
	csvReader.Comma = options.Comma

	// Read CSV headers first, and skip from parsing to object
	if firstRecord, err := csvReader.Read(); err == io.EOF {
		return errors.New("Empty file")
	} else if err != nil {
		return err
	} else if options.isMatchHeaders(firstRecord) {
		// Ensure CSV headers matches
		return &ErrInvalidCSVHeader{headers: firstRecord}
	}

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
			continue
		}

		ch <- convertToPayment(record)
	}

	close(ch)
	return nil
}

func convertToPayment(record []string) *Payment {
	amount, err := strconv.ParseUint(record[1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	month, err := strconv.ParseInt(record[4], 10, 0)
	if err != nil {
		log.Fatal(err)
	}
	year, err := strconv.ParseInt(record[5], 10, 0)
	if err != nil {
		log.Fatal(err)
	}

	return &Payment{
		Amount: amount,
		Card: &Card{
			Name:            record[0],
			Number:          record[2],
			SecurityCode:    record[3],
			ExpirationMonth: time.Month(month),
			ExpirationYear:  int(year),
		},
	}
}
