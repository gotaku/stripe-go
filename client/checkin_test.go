package client

import (
	"fmt"
	"testing"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/currency"
	"github.com/stripe/stripe-go/plan"
)

const testKey = "tGN0bIwXnHdwOa85VABjPdSn8nWY7G7I"

func TestCheckinConnectivity(t *testing.T) {
	c := &Api{}
	c.Init(testKey, nil)

	target, err := c.Account.Get()

	if err != nil {
		t.Error(err)
	}

	if target.Id != "cuD9Rwx8pgmRZRpVe02lsuR9cwp2Bzf7" {
		t.Errorf("Invalid account id: %q\n", target.Id)
	}

	if target.Name != "Stripe Test" {
		t.Errorf("Invalid account name: %q\n", target.Name)
	}

	if target.Email != "test+bindings@stripe.com" {
		t.Errorf("Invalid account email: %q\n", target.Email)
	}
}

func TestCheckinError(t *testing.T) {
	c := &Api{}
	c.Init("bad_key", nil)

	_, err := c.Account.Get()

	if err == nil {
		t.Errorf("Expected an error")
	}

	stripeErr := err.(*stripe.Error)

	if stripeErr.Type != stripe.InvalidRequest {
		t.Errorf("Type %v does not match expected type\n", stripeErr.Type)
	}
}

func TestCheckinPost(t *testing.T) {
	c := &Api{}
	c.Init(testKey, nil)

	charge := &stripe.ChargeParams{
		Amount:   100,
		Currency: currency.USD,
		Card: &stripe.CardParams{
			Name:   "Go Bindings Cardholder",
			Number: "4242424242424242",
			Month:  "12",
			Year:   "24",
		},
	}

	target, err := c.Charges.New(charge)

	if err != nil {
		t.Error(err)
	}

	if target.Amount != charge.Amount {
		t.Errorf("Amount %v does not match expected amount %v\n", target.Amount, charge.Amount)
	}

	if target.Currency != charge.Currency {
		t.Errorf("Currency %q does not match expected currency %q\n", target.Currency, charge.Currency)
	}

	if target.Card.Name != charge.Card.Name {
		t.Errorf("Card name %q does not match expected name %q\n", target.Card.Name, charge.Card.Name)
	}
}

func TestCheckinDel(t *testing.T) {
	c := &Api{}
	c.Init(testKey, nil)

	plan := &stripe.PlanParams{
		Id:       "go_binding",
		Name:     "Go Test Plan",
		Amount:   100,
		Currency: currency.USD,
		Interval: plan.Month,
	}

	_, err := c.Plans.New(plan)

	if err != nil {
		t.Error(err)
	}

	err = c.Plans.Del(plan.Id)

	if err != nil {
		t.Error(err)
	}
}

func TestCheckinList(t *testing.T) {
	const runs = 4
	c := &Api{}
	c.Init(testKey, nil)

	for i := 0; i < runs; i++ {
		plan := &stripe.PlanParams{
			Id:       fmt.Sprintf("go_binding_%v", i),
			Name:     fmt.Sprintf("Go Test Plan %v", i),
			Amount:   100,
			Currency: currency.USD,
			Interval: plan.Month,
		}

		c.Plans.New(plan)
	}

	params := &stripe.PlanListParams{}
	params.Filters.AddFilter("limit", "", "2")
	params.Single = true

	i := c.Plans.List(params)
	for !i.Stop() {
		target, err := i.Next()

		if err != nil {
			t.Error(err)
		}

		if i.Meta() == nil {
			t.Error("No metadata returned")
		}

		if target.Amount != 100 {
			t.Errorf("Amount %v does not match expected value\n", target.Amount)
		}
	}

	for i := 0; i < runs; i++ {
		c.Plans.Del(fmt.Sprintf("go_binding_%v", i))
	}
}
