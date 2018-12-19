package main

import (
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

func createTokenAndCharge(client *omise.Client, payment *Payment) (bool, error) {
	// Creates a token from a test card.
	token, createToken := &omise.Token{}, &operations.CreateToken{
		Name:            payment.Card.Name,
		Number:          payment.Card.Number,
		ExpirationMonth: payment.Card.ExpirationMonth,
		ExpirationYear:  payment.Card.ExpirationYear,
		SecurityCode:    payment.Card.SecurityCode,
	}
	if err := client.Do(token, createToken); err != nil {
		return false, err
	}

	// Creates a charge from the token
	charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
		Amount:   int64(payment.Amount),
		Currency: DefaultCurrency,
		Card:     token.ID,
	}
	if err := client.Do(charge, createCharge); err != nil {
		return false, err
	}

	return charge.Paid, nil
}
