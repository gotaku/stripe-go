// package balance provides the /balance APIs
package balance

import (
	"net/url"
	"strconv"

	stripe "github.com/stripe/stripe-go"
)

const (
	TxAvailable stripe.TransactionStatus = "available"
	TxPending   stripe.TransactionStatus = "pending"

	TxCharge         stripe.TransactionType = "charge"
	TxRefund         stripe.TransactionType = "refund"
	TxAdjust         stripe.TransactionType = "adjustment"
	TxAppFee         stripe.TransactionType = "application_fee"
	TxFeeRefund      stripe.TransactionType = "application_fee_refund"
	TxTransfer       stripe.TransactionType = "transfer"
	TxTransferCancel stripe.TransactionType = "transfer_cancel"
	TxTransferFail   stripe.TransactionType = "transfer_failure"
)

// Client is used to invoke /balance and transaction-related APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// Get returns the details of your balance.
// For more details see https://stripe.com/docs/api#retrieve_balance.
func Get(params *stripe.BalanceParams) (*stripe.Balance, error) {
	return getC().Get(params)
}

func (c Client) Get(params *stripe.BalanceParams) (*stripe.Balance, error) {
	var body *url.Values

	if params != nil {
		body = &url.Values{}
		params.AppendTo(body)
	}

	balance := &stripe.Balance{}
	err := c.B.Call("GET", "/balance", c.Key, body, balance)

	return balance, err
}

// GetTx returns the details of a balance transaction.
// For more details see	https://stripe.com/docs/api#retrieve_balance_transaction.
func GetTx(id string, params *stripe.TxParams) (*stripe.Transaction, error) {
	return getC().GetTx(id, params)
}

func (c Client) GetTx(id string, params *stripe.TxParams) (*stripe.Transaction, error) {
	var body *url.Values

	if params != nil {
		body = &url.Values{}
		params.AppendTo(body)
	}

	balance := &stripe.Transaction{}
	err := c.B.Call("GET", "/balance/history/"+id, c.Key, body, balance)

	return balance, err
}

// List returns a list of balance transactions.
// For more details see https://stripe.com/docs/api#balance_history.
func List(params *stripe.TxListParams) *Iter {
	return getC().List(params)
}

func (c Client) List(params *stripe.TxListParams) *Iter {
	var body *url.Values
	var lp *stripe.ListParams

	if params != nil {
		body = &url.Values{}

		if params.Created > 0 {
			body.Add("created", strconv.FormatInt(params.Created, 10))
		}

		if params.Available > 0 {
			body.Add("available_on", strconv.FormatInt(params.Available, 10))
		}

		if len(params.Currency) > 0 {
			body.Add("currency", params.Currency)
		}

		if len(params.Src) > 0 {
			body.Add("source", params.Src)
		}

		if len(params.Transfer) > 0 {
			body.Add("transfer", params.Transfer)
		}

		if len(params.Type) > 0 {
			body.Add("type", string(params.Type))
		}

		params.AppendTo(body)
		lp = &params.ListParams
	}

	return &Iter{stripe.GetIter(lp, body, func(b url.Values) ([]interface{}, stripe.ListMeta, error) {
		type transactionList struct {
			stripe.ListMeta
			Values []*stripe.Transaction `json:"data"`
		}

		list := &transactionList{}
		err := c.B.Call("GET", "/balance/history", c.Key, &b, list)

		ret := make([]interface{}, len(list.Values))
		for i, v := range list.Values {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is a iterator for list responses.
type Iter struct {
	Iter *stripe.Iter
}

// Next returns the next value in the list.
func (i *Iter) Next() (*stripe.Transaction, error) {
	t, err := i.Iter.Next()
	if err != nil {
		return nil, err
	}

	return t.(*stripe.Transaction), err
}

// Stop returns true if there are no more iterations to be performed.
func (i *Iter) Stop() bool {
	return i.Iter.Stop()
}

// Meta returns the list metadata.
func (i *Iter) Meta() *stripe.ListMeta {
	return i.Iter.Meta()
}

func getC() Client {
	return Client{stripe.GetBackend(), stripe.Key}
}
