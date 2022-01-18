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
)

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
func (c *Client) post(route string, d interface{}, r interface{}) error {
	return c.req("POST", route, d, r)
}

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

type FillBidParams struct {
	fillStr string
	skynet  bool
	b64     bool
}

func (c *Client) FillBid(fillStr string, skynet bool, b64 bool) (r string, err error) {
	params := FillBidParams{fillStr: fillStr, skynet: skynet, b64: b64}
	err = c.post("/bids/fill", params, r)
	return
}

type PlaceBidParams struct {
	inStr  string
	outStr string
	skynet bool
	b64    bool
}

func (c *Client) PlaceBid(inStr string, outStr string, skynet bool, b64 bool) (r string, err error) {
	params := PlaceBidParams{inStr: inStr, outStr: outStr, skynet: skynet, b64: b64}
	err = c.post("/bids/place", params, r)
	return
}

// NewClient returns a Client that connects to the specified embarcadero server.
func NewClient(addr string) *Client {
	return &Client{addr: addr}
}
