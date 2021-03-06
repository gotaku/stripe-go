// package refund provides the /refunds APIs
package refund

import (
	"fmt"
	"net/url"
	"strconv"

	stripe "github.com/stripe/stripe-go"
)

// Client is used to invoke /refunds APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New refunds a charge previously created.
// For more details see https://stripe.com/docs/api#refund_charge.
func New(params *stripe.RefundParams) (*stripe.Refund, error) {
	return getC().New(params)
}

func (c Client) New(params *stripe.RefundParams) (*stripe.Refund, error) {
	body := &url.Values{}

	if params.Amount > 0 {
		body.Add("amount", strconv.FormatUint(params.Amount, 10))
	}

	if params.Fee {
		body.Add("refund_application_fee", strconv.FormatBool(params.Fee))
	}

	params.AppendTo(body)

	refund := &stripe.Refund{}
	err := c.B.Call("POST", fmt.Sprintf("/charges/%v/refunds", params.Charge), c.Key, body, refund)

	return refund, err
}

// Get returns the details of a refund.
// For more details see https://stripe.com/docs/api#retrieve_refund.
func Get(id string, params *stripe.RefundParams) (*stripe.Refund, error) {
	return getC().Get(id, params)
}

func (c Client) Get(id string, params *stripe.RefundParams) (*stripe.Refund, error) {
	body := &url.Values{}

	params.AppendTo(body)

	refund := &stripe.Refund{}
	err := c.B.Call("GET", fmt.Sprintf("/charges/%v/refunds/%v", params.Charge, id), c.Key, body, refund)

	return refund, err
}

// Update updates a refund's properties.
// For more details see https://stripe.com/docs/api#update_refund.
func Update(id string, params *stripe.RefundParams) (*stripe.Refund, error) {
	return getC().Update(id, params)
}

func (c Client) Update(id string, params *stripe.RefundParams) (*stripe.Refund, error) {
	body := &url.Values{}

	params.AppendTo(body)

	refund := &stripe.Refund{}
	err := c.B.Call("POST", fmt.Sprintf("/charges/%v/refunds/%v", params.Charge, id), c.Key, body, refund)

	return refund, err
}

// List returns a list of refunds.
// For more details see https://stripe.com/docs/api#list_refunds.
func List(params *stripe.RefundListParams) *Iter {
	return getC().List(params)
}

func (c Client) List(params *stripe.RefundListParams) *Iter {
	body := &url.Values{}
	var lp *stripe.ListParams

	params.AppendTo(body)
	lp = &params.ListParams

	return &Iter{stripe.GetIter(lp, body, func(b url.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.RefundList{}
		err := c.B.Call("GET", fmt.Sprintf("/charges/%v/refunds", params.Charge), c.Key, &b, list)

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
func (i *Iter) Next() (*stripe.Refund, error) {
	r, err := i.Iter.Next()
	if err != nil {
		return nil, err
	}

	return r.(*stripe.Refund), err
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
