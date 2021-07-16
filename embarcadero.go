package embarcadero

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"gitlab.com/NebulousLabs/Sia/modules"
	"gitlab.com/NebulousLabs/Sia/types"
)

// BidPrefix is the identifier for embarcadero bids in arbitrary data.
var BidPrefix = append(modules.PrefixNonSia[:], "embarcadero"...)

// A Trade is a SF<->SC transaction.
type Trade struct {
	Bid         Bid
	Transaction types.Transaction
	Height      types.BlockHeight
}

// A Bid is an unfilled trade.
type Bid struct {
	Transaction types.Transaction
	ID          types.OutputID
	Height      types.BlockHeight
	SF, SC      types.Currency
	OfferingSF  bool
	Invalid     bool `json:"omitempty"`
}

// A Client connects to an embarcadero server.
type Client struct {
	addr string
}

func (c *Client) req(method string, route string, data, resp interface{}) error {
	var body io.Reader
	if data != nil {
		js, _ := json.Marshal(data)
		body = bytes.NewReader(js)
	}
	req, err := http.NewRequest(method, fmt.Sprintf("%v%v", c.addr, route), body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer io.Copy(ioutil.Discard, r.Body)
	defer r.Body.Close()
	if r.StatusCode != 200 {
		err, _ := ioutil.ReadAll(r.Body)
		return errors.New(strings.TrimSpace(string(err)))
	}
	if resp == nil {
		return nil
	}
	return json.NewDecoder(r.Body).Decode(resp)
}

func (c *Client) get(route string, r interface{}) error { return c.req("GET", route, nil, r) }

// Trades returns all trades known to the embarcadero server.
func (c *Client) Trades() (ts []Trade, err error) {
	err = c.get("/trades", &ts)
	return
}

// Bids returns all bids known to the embarcadero server.
func (c *Client) Bids() (bs []Bid, err error) {
	err = c.get("/bids", &bs)
	return
}

// NewClient returns a Client that connects to the specified embarcadero server.
func NewClient(addr string) *Client {
	return &Client{addr: addr}
}
