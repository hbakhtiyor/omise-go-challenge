package main

import (
	"fmt"
	"strings"

	"golang.org/x/text/message"
)

type DonationResult struct {
	TotalPayments int
	SuccessSum    uint64
	FailSum       uint64
	TopPayments   []*Payment
	Currency      string
}

var p *message.Printer

func (r *DonationResult) sortTopPayments(payment *Payment) {
	len := len(r.TopPayments)
	for i, p := range r.TopPayments {
		if payment.Amount > p.Amount {
			r.TopPayments = append(r.TopPayments[:i], append([]*Payment{payment}, r.TopPayments[i:len-1]...)...)
			break
		}
	}
}

func (r DonationResult) printSummary(clean bool) {
	if clean && p != nil {
		r.cleanSummary()
	}
	if p == nil {
		p = message.NewPrinter(message.MatchLanguage("en"))
	}
	r.Currency = strings.ToUpper(r.Currency)

	p.Printf("\n%23s %s %14d.00\n", "total received:", r.Currency, r.SuccessSum+r.FailSum)
	p.Printf("%23s %s %14d.00\n", "successfully donated:", r.Currency, r.SuccessSum)
	p.Printf("%23s %s %14d.00\n\n", "faulty donation:", r.Currency, r.FailSum)
	p.Printf("%23s %s %17.2f\n", "average per person:", r.Currency, float64(r.SuccessSum+r.FailSum)/float64(r.TotalPayments))
	p.Printf("%23s", "top donors:")
	for i, payment := range r.TopPayments {
		if payment == nil || payment.Card == nil {
			// reserving line for future top donors, clean-up linces will be mess
			p.Printf("\n")
			continue
		}
		if i == 0 {
			p.Printf(" %s\n", payment.Card.Name)
		} else {
			p.Printf("%23s %s\n", "", payment.Card.Name)
		}
	}
}

func (r DonationResult) cleanSummary() {
	if p != nil {
		len := len(r.TopPayments)
		for i := 0; i < len+6; i++ {
			fmt.Printf("\033[1A\033[K")
		}
	}
}
